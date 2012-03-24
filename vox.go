package glvox

import (
	"fmt"
)

type Tracer interface {
	Trace(pos, dir Vec3) (dest Vec3, hit bool)
}

type Setter interface {
	Set(x, y, z int32, v int32)
}

type Dimensioner interface {
	Dim(w, h, d int32)
}

type DimSetter interface {
	Dimensioner
	Setter
}

type Octree struct {
	Index []int32
	WHD int32
}

type Grid  struct {
	data []int32
	W, H, D int32
}

type vox struct {
	d Vec3
	size float32
	alpha float32
}

const (
	MaxSteps = 10
)

func NewOctree() *Octree {
	oct := new(Octree)
	oct.Index = append(oct.Index, 0, 0, 0, 0,  0, 0, 0, 0)
	return oct
}

func (oct *Octree) Dim(w, h, d int32) {

	dim := w
	if h > dim { dim = h }
	if d > dim { dim = d }

	pow2 := int32(1)
	for dim > pow2 { pow2 *= 2 }

	oct.WHD = pow2
}

func (oct *Octree) Trace(ro, rd Vec3) (pos Vec3, hit bool) {

	var s Vec3
	s.X = float32(1.0); if rd.X < 0 { s.X = -1.0 }
	s.Y = float32(1.0); if rd.Y < 0 { s.Y = -1.0 }
	s.Z = float32(1.0); if rd.Z < 0 { s.Z = -1.0 }

	pos, hit = oct.boundingBox(ro, rd)
	if !hit { return }

	for i := 0; i < MaxSteps; i++ {
		v := oct.Voxel(pos, s)
		if v.alpha > 0.0 { hit = true; return }

		dist := v.d
		f := s.Mul(v.size).Minus(dist)
		f.X /= rd.X; f.Y /= rd.Y; f.Z /= rd.Z

		fmin := float32(100.0)
		if f.X > 0.0 && f.X < fmin { fmin = f.X }
		if f.Y > 0.0 && f.Y < fmin { fmin = f.Y }
		if f.Z > 0.0 && f.Z < fmin { fmin = f.Z }

		pos = pos.Plus(rd.Mul(fmin))
	}

	return
}

func (oct *Octree) boundingBox(p, d Vec3) (pos Vec3, hit bool) {

	size := float32(oct.WHD)

	if	p.X >= 0.0 && p.X <= size &&
		p.Y >= 0.0 && p.Y <= size &&
		p.Z >= 0.0 && p.Z <= size {

		return p, true
	}

	if d.X > 0.0 {
		pos, hit = face(p, d, 0.0, size, 0)
		if hit { return }
	} else {
		pos, hit = face(p, d, size, size, 0)
		if hit { return }
	}

	if d.Y > 0.0 {
		pos, hit = face(p, d, 0.0, size, 1)
		if hit { return }
	} else {
		pos, hit = face(p, d, size, size, 1)
		if hit { return }
	}

	if d.Z > 0.0 {
		pos, hit = face(p, d, 0.0, size, 3)
		if hit { return }
	} else {
		pos, hit = face(p, d, size, size, 3)
		if hit { return }
	}

	return
}

func face(p, d Vec3, off, size float32, axis int) (pos Vec3, hit bool) {

	var f float32
	switch axis {
	case 0:
		if d.X == 0.0 { return p, false }
		f = (off - p.X) / d.X
	case 1:
		if d.Y == 0.0 { return p, false }
		f = (off - p.Y) / d.Y
	case 2:
		if d.Z == 0.0 { return p, false }
		f = (off - p.Z) / d.Z
	}

	pos = p.Plus(d.Mul(f))
	hit = pos.X >= 0 && pos.X <= size &&
		pos.Y >= 0 && pos.Y <= size &&
		pos.Z >= 0 && pos.Z <= size

	return
}

func (oct *Octree) Voxel(pos, dir Vec3) vox {
	x, y, z := int32(pos.X), int32(pos.Y), int32(pos.Z)
	if dir.X < 0.0 { x-- }
	if dir.Y < 0.0 { y-- }
	if dir.Z < 0.0 { z-- }

	val, size := oct.Get(x, y, z)

	s := float32(size) / 2.0
	center :=
		Vec3{float32(x), float32(y), float32(z)}.Plus(
		Vec3{s, s, s})
	dist := pos.Minus(center)
	v := vox{dist, s, float32(val)}
	return v
}

func (oct *Octree) Get(x, y, z int32) (val int32, whd int32) {

	whd = oct.WHD
	if x < 0 || x >= whd || y < 0 || y >= whd || z < 0 || z >= whd {
		return
	}

	var i, off int32 = 0, 0
	for whd > 1 {

		whd >>= 1
		off = 0

		if z >= whd { off += 4; z -= whd }
		if y >= whd { off += 2; y -= whd }
		if x >= whd { off += 1; x -= whd }

		i = oct.Index[i*8 + off]
		if i <= 0 { val = -i; return }
	}

	return
}

func (oct *Octree) Set(x, y, z int32, v int32) {

	whd := oct.WHD
	if x < 0 || x >= whd || y < 0 || y >= whd || z < 0 || z >= whd {
		return
	}

	var i, off int32 = 0, 0
	for whd > 1 {

		whd >>= 1
		off = 0

		if z >= whd { off += 4; z -= whd }
		if y >= whd { off += 2; y -= whd }
		if x >= whd { off += 1; x -= whd }

		idx := oct.Index[i*8 + off]
		if -idx == v { return }
		if idx == e {
			oct.Index[i*8 + off] = -v
			return
		}

		if idx <= 0 {
			if whd > 1 {
				idx = oct.newIndex(idx)
				oct.Index[i*8 + off] = idx
			} else {
				oct.Index[i*8 + off] = -v
				return
			}
		}

		i = idx
	}
}

func (oct *Octree) newIndex(v int32) int32 {
	idx := len(oct.Index) / 8
	oct.Index = append(oct.Index, v, v, v, v,  v, v, v, v)
	//oct.Index = append(oct.Index, e, e, e, e,  e, e, e, e)
	return int32(idx)
}

func (oct *Octree) String() string {

	printer := func(n int32) string {
		if n > 0 { return fmt.Sprintf("%d", n) }
		return fmt.Sprintf("%c[1;34;40m%d%c[0;37;40m",
			0x1b, -n, 0x1b)
	}

	s := ""
	for i := 0; i < len(oct.Index) / 8; i++ {
		idx := oct.Index[i*8]
		s += fmt.Sprint(i, ": [", printer(idx))
		for o := 1; o < 8; o++ {
			idx := oct.Index[i*8 + o]
			s += fmt.Sprint(", ", printer(idx))
		}
		s += "]\n"
	}

	return s
}

func NewGrid() *Grid {
	g := new(Grid)
	return g
}

func (g *Grid) Dim(w, h, d int32) {
	g.data = make([]int32, w*h*d)
	g.W = w; g.H = h; g.D = d
}

func (g *Grid) Set(x, y, z int32, val int32) {
	g.data[z*g.H*g.W + y*g.W + x] = val
}

func (g *Grid) Get(x, y, z int32) int32 {
	return g.data[z*g.H*g.W + y*g.W + x]
}

func (g *Grid) Trace(ro, rd Vec3) (pos Vec3, hit bool) {

	sx := float32(1.0); if rd.X < 0 { sx = -1.0 }
	sy := float32(1.0); if rd.Y < 0 { sy = -1.0 }
	sz := float32(1.0); if rd.Z < 0 { sz = -1.0 }

	hit = false
	pos = ro
	for i := 0; i < MaxSteps; i++ {
		v := g.Voxel(pos, rd)
		if v.alpha > 0.0 { hit = true; return }

		dist := v.d
		fx := (sx * v.size - dist.X) / rd.X
		fy := (sy * v.size - dist.Y) / rd.Y
		fz := (sz * v.size - dist.Z) / rd.Z

		f := float32(100.0)
		if fx > 0.0 && fx < f { f = fx }
		if fy > 0.0 && fy < f { f = fy }
		if fz > 0.0 && fz < f { f = fz }

		pos = pos.Plus(rd.Mul(f))
	}

	return
}

func (g *Grid) Voxel(pos, dir Vec3) vox {
	x, y, z := int32(pos.X), int32(pos.Y), int32(pos.Z)
	if dir.X < 0 { x-- }
	if dir.Y < 0 { y-- }
	if dir.Z < 0 { z-- }

	val := int32(0)
	if x >= 0 && x < g.W && y >= 0 && y < g.H && z >= 0 && z < g.D {
		val = g.Get(x, y, z)
	}

	s := float32(.5)
	center := Vec3{float32(x), float32(y), float32(z)}.Plus(Vec3{s, s, s})
	dist := pos.Minus(center)
	v := vox{dist, s, float32(val)}
	return v
}

func (g *Grid) String() string {
	s := ""

	for z := int32(0); z < g.D; z++ {
		for y := int32(0); y < g.H; y++ {
			for x := int32(0); x < g.W; x++ {
				s += fmt.Sprint(" ", g.Get(x, y, z))
			}
			s += "\n"
		}
		s += "\n"
	}

	return s
}
