package glvox

import (
	"fmt"
)

type Grid  struct {
	data []int32
	W, H, D int32
}

func NewGrid(w, h, d int32) *Grid {
	g := new(Grid)
	g.data = make([]int32, w*h*d)
	g.W = w; g.H = h; g.D = d
	return g
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
		if v.Value > 0.0 { hit = true; return }

		dist := v.Dist
		fx := (sx * v.Size - dist.X) / rd.X
		fy := (sy * v.Size - dist.Y) / rd.Y
		fz := (sz * v.Size - dist.Z) / rd.Z

		f := float32(100.0)
		if fx > 0.0 && fx < f { f = fx }
		if fy > 0.0 && fy < f { f = fy }
		if fz > 0.0 && fz < f { f = fz }

		pos = pos.Plus(rd.Mul(f))
	}

	return
}

func (g *Grid) Voxel(pos, dir Vec3) Vox {
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
	v := Vox { dist, s, float32(val) }
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
