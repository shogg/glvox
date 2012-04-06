package glvox

import (
	"fmt"
)

const (
	MaxSteps = 10
)

type Octree struct {
	Index []int
	Size int
}

func NewOctree(size int) *Octree {

	oct := new(Octree)
	oct.Index = append(oct.Index, 0, 0, 0, 0,  0, 0, 0, 0)

	pow2 := 1
	for size > pow2 { pow2 *= 2 }

	oct.Size = pow2
	return oct
}

func (oct *Octree) Trace(ro, rd Vec3) (pos Vec3, hit bool) {

	var s Vec3
	s.X = float32(1.0); if rd.X < 0 { s.X = -1.0 }
	s.Y = float32(1.0); if rd.Y < 0 { s.Y = -1.0 }
	s.Z = float32(1.0); if rd.Z < 0 { s.Z = -1.0 }

	for i := 0; i < MaxSteps; i++ {
		v := oct.Voxel(pos, s)
		if v.Value > 0.0 { hit = true; return }

		dist := v.Dist
		f := s.Mul(v.Size).Minus(dist)
		f.X /= rd.X; f.Y /= rd.Y; f.Z /= rd.Z

		fmin := float32(100.0)
		if f.X > 0.0 && f.X < fmin { fmin = f.X }
		if f.Y > 0.0 && f.Y < fmin { fmin = f.Y }
		if f.Z > 0.0 && f.Z < fmin { fmin = f.Z }

		pos = pos.Plus(rd.Mul(fmin))
	}

	return
}

func (oct *Octree) Voxel(pos, dir Vec3) Vox {
	x, y, z := int(pos.X), int(pos.Y), int(pos.Z)
	if dir.X < 0.0 { x-- }
	if dir.Y < 0.0 { y-- }
	if dir.Z < 0.0 { z-- }

	val, size := oct.Get(x, y, z)

	s := float32(size) / 2.0
	center :=
		Vec3{float32(x), float32(y), float32(z)}.Plus(
		Vec3{s, s, s})
	dist := pos.Minus(center)
	v := Vox { dist, s, float32(val) }
	return v
}

func (oct *Octree) Get(x, y, z int) (val int, size int) {

	size = oct.Size
	if x < 0 || x >= size || y < 0 || y >= size || z < 0 || z >= size {
		return
	}

	var i, off int = 0, 0
	for size > 1 {

		size >>= 1
		off = 0

		if z >= size { off += 4; z -= size }
		if y >= size { off += 2; y -= size }
		if x >= size { off += 1; x -= size }

		i = oct.Index[i*8 + off]
		if i <= 0 { val = -i; return }
	}

	return
}

func (oct *Octree) Set(x, y, z int, v int) {

	size := oct.Size
	if x < 0 || x >= size || y < 0 || y >= size || z < 0 || z >= size {
		return
	}

	var i, off int = 0, 0
	for size > 1 {

		size >>= 1
		off = 0

		if z >= size { off += 4; z -= size }
		if y >= size { off += 2; y -= size }
		if x >= size { off += 1; x -= size }

		idx := oct.Index[i<<3 + off]
		if -idx == v { return }

		if idx <= 0 {
			if size > 1 {
				idx = oct.newIndex(idx)
				oct.Index[i<<3 + off] = idx
			} else {
				oct.Index[i<<3 + off] = -v
				return
			}
		}

		i = idx
	}
}

func (oct *Octree) newIndex(v int) int {
	idx := len(oct.Index) / 8
	oct.Index = append(oct.Index, v, v, v, v,  v, v, v, v)
	return int(idx)
}

func (oct *Octree) String() string {

	printer := func(n int) string {
		if n > 0 { return fmt.Sprintf("%d", n) }
		return fmt.Sprintf("%c[0;37;40m%d%c[1;37;40m",
			0x1b, -n, 0x1b)
	}

	length := len(oct.Index) / 8
	if length > 100 { length = 100 }

	s := ""
	for i := 0; i < length; i++ {
		idx := oct.Index[i<<3]
		s += fmt.Sprint(i, ": [", printer(idx))
		for o := 1; o < 8; o++ {
			idx := oct.Index[i<<3 + o]
			s += fmt.Sprint(", ", printer(idx))
		}
		s += "]\n"
	}

	return s
}

