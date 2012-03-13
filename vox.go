package glvox

import (
	"fmt"
)

type Tracer interface {
	Trace(pos, dir vec3) (dest vec3, hit bool)
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
	data []float32
	WHD int32
}

type Grid  struct {
	data []float32
	W, H, D int32
}

type vox struct {
	d vec3
	size float32
	alpha float32
}

const (
	maxSteps = 10
	IdxSize = 9
)

func NewOctree() *Octree {
	oct := new(Octree)
	oct.Index = make([]int32, IdxSize)
	oct.data = make([]float32, 2)
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

func (oct *Octree) Trace(ro, rd vec3) (pos vec3, hit bool) {

	var s vec3
	s.x = float32(1.0); if rd.x < 0 { s.x = -1.0 }
	s.y = float32(1.0); if rd.y < 0 { s.y = -1.0 }
	s.z = float32(1.0); if rd.z < 0 { s.z = -1.0 }

	pos, hit = oct.boundingBox(ro, rd)
	if !hit { return }

	for i := 0; i < maxSteps; i++ {
		v := oct.Voxel(pos, s)
		fmt.Println(pos, "vox", v)
		if v.alpha > 0.0 { hit = true; return }

		dist := v.d
		f := s.Mul(v.size).Minus(dist)
		f.x /= rd.x; f.y /= rd.y; f.z /= rd.z

		fmin := float32(100.0)
		if f.x > 0.0 && f.x < fmin { fmin = f.x }
		if f.y > 0.0 && f.y < fmin { fmin = f.y }
		if f.z > 0.0 && f.z < fmin { fmin = f.z }

		pos = pos.Plus(rd.Mul(fmin))
	}

	return
}

func (oct *Octree) boundingBox(p, d vec3) (pos vec3, hit bool) {

	size := float32(oct.WHD)

	if	p.x >= 0.0 && p.x <= size &&
		p.y >= 0.0 && p.y <= size &&
		p.z >= 0.0 && p.z <= size {

		return p, true
	}

	if d.x > 0.0 {
		pos, hit = face(p, d, 0.0, size, 0)
		if hit { return }
	} else {
		pos, hit = face(p, d, size, size, 0)
		if hit { return }
	}

	if d.y > 0.0 {
		pos, hit = face(p, d, 0.0, size, 1)
		if hit { return }
	} else {
		pos, hit = face(p, d, size, size, 1)
		if hit { return }
	}

	if d.z > 0.0 {
		pos, hit = face(p, d, 0.0, size, 3)
		if hit { return }
	} else {
		pos, hit = face(p, d, size, size, 3)
		if hit { return }
	}

	return
}

func face(p, d vec3, off, size float32, axis int) (pos vec3, hit bool) {

	var f float32
	switch axis {
	case 0:
		if d.x == 0.0 { return p, false }
		f = (off - p.x) / d.x
	case 1:
		if d.y == 0.0 { return p, false }
		f = (off - p.y) / d.y
	case 2:
		if d.z == 0.0 { return p, false }
		f = (off - p.z) / d.z
	}

	pos = p.Plus(d.Mul(f))
	hit = pos.x >= 0 && pos.x <= size &&
		pos.y >= 0 && pos.y <= size &&
		pos.z >= 0 && pos.z <= size

	return
}

func (oct *Octree) Voxel(pos, dir vec3) vox {
	x, y, z := int32(pos.x), int32(pos.y), int32(pos.z)
	if dir.x < 0.0 { x-- }
	if dir.y < 0.0 { y-- }
	if dir.z < 0.0 { z-- }

	val, size := oct.Get(x, y, z)
	alpha := oct.data[val]

	s := float32(size) / 2.0
	center :=
		vec3{float32(x), float32(y), float32(z)}.Plus(
		vec3{s, s, s})
	dist := pos.Minus(center)
	v := vox{dist, s, alpha}
	return v
}

func (oct *Octree) Get(x, y, z int32) (val int32, whd int32) {

	whd = oct.WHD
	if x < 0 || x >= whd || y < 0 || y >= whd || z < 0 || z >= whd {
		return
	}

	var i, off int32 = 0, 0
	for whd > 1 {

		val = oct.Index[i*IdxSize + 8]
		if val != -1 {
			return val, whd
		}

		whd >>= 1
		off = 0

		if z >= whd { off += 4; z -= whd }
		if y >= whd { off += 2; y -= whd }
		if x >= whd { off += 1; x -= whd }

		i = oct.Index[i*IdxSize + off]
		if i == 0 { panic("data corrupted") }
	}

	val = i
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

		val := oct.Index[i*IdxSize + 8]
		if val == v { return }

		oct.Index[i*IdxSize + 8] = -1
		if val != -1 && whd == 1 {
			for o := int32(0); o < 8; o++ {
				oct.Index[i*IdxSize + o] = val
			}
		}

		if z >= whd { off += 4; z -= whd }
		if y >= whd { off += 2; y -= whd }
		if x >= whd { off += 1; x -= whd }

		idx := oct.Index[i*IdxSize + off]
		if idx == 0 || idx == val {
			if whd > 1 {
				idx = oct.newIndex()
				oct.Index[i*IdxSize + off] = idx
				oct.Index[idx*IdxSize + 8] = v
			} else {
				oct.Index[i*IdxSize + off] = v
			}
		}

		i = idx
	}
}

func (oct *Octree) newIndex() int32 {
	idx := len(oct.Index) / IdxSize
	oct.Index = append(oct.Index, 0, 0, 0, 0,  0, 0, 0, 0,  0)
	return int32(idx)
}

func (g *Grid) Trace(ro, rd vec3) (pos vec3, hit bool) {

	sx := float32(1.0); if rd.x < 0 { sx = -1.0 }
	sy := float32(1.0); if rd.y < 0 { sy = -1.0 }
	sz := float32(1.0); if rd.z < 0 { sz = -1.0 }

	hit = false
	pos = ro
	for i := 0; i < maxSteps; i++ {
		v := g.Voxel(pos, rd)
		if v.alpha == 1.0 { hit = true; return }

		dist := v.d
		fx := (sx * v.size - dist.x) / rd.x
		fy := (sy * v.size - dist.y) / rd.y
		fz := (sz * v.size - dist.z) / rd.z

		f := float32(100.0)
		if fx > 0.0 && fx < f { f = fx }
		if fy > 0.0 && fy < f { f = fy }
		if fz > 0.0 && fz < f { f = fz }

//		fmt.Printf("f=%f,\tpos=%v\n", f, pos, )
		pos = pos.Plus(rd.Mul(f))
	}

	return
}

func (g *Grid) Voxel(pos, dir vec3) vox {
	x, y, z := int32(pos.x), int32(pos.y), int32(pos.z)
	if dir.x < 0 { x-- }
	if dir.y < 0 { y-- }
	if dir.z < 0 { z-- }

	alpha := float32(0.0)
	if x >= 0 && x < g.W && y >= 0 && y < g.H && z >= 0 && z < g.D {
		alpha = g.data[x + g.D*y + g.W*g.H*z]
	}

	s := float32(.5)
	center := vec3{float32(x), float32(y), float32(z)}.Plus(vec3{s, s, s})
	dist := pos.Minus(center)
	v := vox{dist, s, alpha}
	return v
}
