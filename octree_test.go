package glvox

import (
	"testing"
	"fmt"
)

func buildOctree() *Octree {

	s, s1 := int32(16), int32(15)

	oct := NewOctree(s)
	for z := int32(0); z < s; z++ {
		for y := int32(0); y < s; y++ {
			for x := int32(0); x < s; x++ {
				oct.Set(x, y, z, 0)

				if x == 0 && y == 0 && z == 0 {
					oct.Set(0, 0, 0, 5)
				}

				if x == s1 && y == s1 && z == s1 {
					oct.Set(s1, s1, s1, 6)
				}
			}
		}
	}

	return oct
}

func TestOctreeTrace(t *testing.T) {

	oct := buildOctree()

	ro := Vec3{ 1.0, 1.0, 1.0}
	rd := Vec3{ 1.0, 1.0, 1.0}.Normalize()
	pos, hit := oct.Trace(ro, rd)
	exp := Vec3{15.0, 15.0, 15.0}
	if pos != exp || !hit {
		t.Errorf("test1 hit expected at %v, was %v", exp, pos)
	}

	ro = Vec3{ 3.0, 3.0, 3.0}
	rd = Vec3{-1.0,-1.0,-1.0}.Normalize()
	pos, hit = oct.Trace(ro, rd)
	exp = Vec3{1.0, 1.0, 1.0}
	if pos != exp || !hit {
		t.Errorf("test2 hit expected at %v, was %v", exp, pos)
	}

	ro = Vec3{-10.0,-10.0,-10.0}
	rd = Vec3{ 1.0, 1.0, 1.0}.Normalize()
	pos, hit = oct.Trace(ro, rd)
	exp = Vec3{0.0, 0.0, 0.0}
	if pos != exp || !hit {
		t.Errorf("test3 hit expected at %v, was %v", exp, pos)
	}

	ro = Vec3{ 20.0, 20.0, 20.0}
	rd = Vec3{-1.0,-1.0,-1.0}.Normalize()
	pos, hit = oct.Trace(ro, rd)
	exp = Vec3{16.0, 16.0, 16.0}
	if pos != exp || !hit {
		t.Errorf("test4 hit expected at %v, was %v", exp, pos)
	}
}

func TestOctreeGet(t *testing.T) {

	oct := buildOctree()

	if v, s := oct.Get(0, 0, 0); v != 5 || s != 1 {
		t.Errorf("test1 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(1, 1, 1); v != 0 || s != 1 {
		t.Errorf("test2 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(2, 2, 2); v != 0 || s != 2 {
		t.Errorf("test3 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(15, 15, 15); v != 6 || s != 1 {
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

	//fmt.Println(grid)
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

