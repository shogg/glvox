package glvox

import (
//	"github.com/shogg/glvox"
	"fmt"
	"testing"
)

func buildOctree() *Octree {

	oct := NewOctree()
	oct.Dim(4, 4, 4)

	for z := int32(0); z < 4; z++ {
		for y := int32(0); y < 4; y++ {
			for x := int32(0); x < 4; x++ {
				oct.Set(x, y, z, 0)

				if x == 0 && y == 0 && z == 0 {
					oct.Set(0, 0, 0, 1)
				}

				if x == 3 && y == 3 && z == 3 {
					oct.Set(3, 3, 3, 2)
				}
			}
		}
	}

	return oct
}

func buildGrid() *Grid {

	g := NewGrid()
	g.Dim(4, 4, 4)

	g.Set(0, 0, 0, 1)
	g.Set(3, 3, 3, 2)

	return g
}

func TestOctreeTrace(t *testing.T) {

	oct := buildOctree()

	ro := vec3{ 1.0, 1.0, 1.0}
	rd := vec3{ 1.0, 1.0, 1.0}.Normalize()
	pos, hit := oct.Trace(ro, rd)
	exp := vec3{3.0, 3.0, 3.0}
	if pos != exp || !hit {
		t.Errorf("test1 hit expected at %v, was %v", exp, pos)
	}

	ro = vec3{ 3.0, 3.0, 3.0}
	rd = vec3{-1.0,-1.0,-1.0}.Normalize()
	pos, hit = oct.Trace(ro, rd)
	exp = vec3{1.0, 1.0, 1.0}
	if pos != exp || !hit {
		t.Errorf("test2 hit expected at %v, was %v", exp, pos)
	}

	ro = vec3{-10.0,-10.0,-10.0}
	rd = vec3{ 1.0, 1.0, 1.0}.Normalize()
	pos, hit = oct.Trace(ro, rd)
	exp = vec3{0.0, 0.0, 0.0}
	if pos != exp || !hit {
		t.Errorf("test3 hit expected at %v, was %v", exp, pos)
	}

	ro = vec3{ 10.0, 10.0, 10.0}
	rd = vec3{-1.0,-1.0,-1.0}.Normalize()
	pos, hit = oct.Trace(ro, rd)
	exp = vec3{4.0, 4.0, 4.0}
	if pos != exp || !hit {
		t.Errorf("test4 hit expected at %v, was %v", exp, pos)
	}
}

func TestGridTrace(t *testing.T) {

	g := buildGrid()

	ro := vec3{ 1.0, 1.0, 1.0}
	rd := vec3{ 1.0, 1.0, 1.0}.Normalize()

	pos, hit := g.Trace(ro, rd)
	exp := vec3{3.0, 3.0, 3.0}
	if pos != exp || !hit {
		t.Errorf("test1 hit expected at %v, was %v", exp, pos)
	}

	ro = vec3{ 3.0, 3.0, 3.0}
	rd = vec3{-1.0,-1.0,-1.0}.Normalize()

	pos, hit = g.Trace(ro, rd)
	exp = vec3{1.0, 1.0, 1.0}
	if pos != exp || !hit {
		t.Errorf("test2 hit expected at %v, was %v", exp, pos)
	}
}

func TestOctreeGet(t *testing.T) {

	oct := buildOctree()

	if v, s := oct.Get(0, 0, 0); v != 1 || s != 1 {
		t.Errorf("test1 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(1, 1, 1); v != 0 || s != 1 {
		t.Errorf("test2 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(2, 2, 2); v != 0 || s != 1 {
		t.Errorf("test3 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(3, 3, 3); v != 2 || s != 1 {
		t.Errorf("test4 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(0, 0, 3); v != 0 || s != 2 {
		t.Errorf("test5 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(0, 3, 0); v != 0 || s != 2 {
		t.Errorf("test6 v=%d d=%d", v, s)
	}
}

func TestOctreeGetAll(t *testing.T) {

	oct := buildOctree()
	grid := buildGrid()

	fmt.Println(grid)
	fmt.Println(oct)

	for z := int32(0); z < grid.D; z++ {
		for y := int32(0); y < grid.D; y++ {
			for x := int32(0); x < grid.D; x++ {
				expected := grid.Get(x, y, z)
				actual, _ := oct.Get(x, y, z)
				if expected != actual {
					t.Errorf("(%d, %d, %d): expected %d, was %d",
						x, y, z, expected, actual)
				}
			}
		}
	}
}
