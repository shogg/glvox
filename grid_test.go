package glvox

import (
	"testing"
)

func buildGrid() *Grid {

	s, s1 := int32(16), int32(15)

	g := NewGrid(s, s, s)
	for i := 0; i < int(s*s*s); i++ {
		g.data[i] = 0
	}

	g.Set(0, 0, 0, 5)
	g.Set(s1, s1, s1, 6)

	return g
}

func TestGridTrace(t *testing.T) {

	g := buildGrid()

	ro := Vec3{ 10.0, 10.0, 10.0}
	rd := Vec3{ 1.0, 1.0, 1.0}.Normalize()

	pos, hit := g.Trace(ro, rd)
	exp := Vec3{15.0, 15.0, 15.0}
	if pos != exp || !hit {
		t.Errorf("test1 hit expected at %v, was %v", exp, pos)
	}

	ro = Vec3{ 3.0, 3.0, 3.0}
	rd = Vec3{-1.0,-1.0,-1.0}.Normalize()

	pos, hit = g.Trace(ro, rd)
	exp = Vec3{1.0, 1.0, 1.0}
	if pos != exp || !hit {
		t.Errorf("test2 hit expected at %v, was %v", exp, pos)
	}
}
