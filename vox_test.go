package glvox

import (
//	"github.com/shogg/glvox"
	"fmt"
	"testing"
)

var (
	data = []float32 {
		1.0, 0.0, 0.0, 0.0, // 0..3
		0.0, 0.0, 6.0, 7.0,	// 4..7
		0.0, 0.0, 0.0, 0.0, // 8..11
		0.0, 0.0, 0.0, 0.0, // 12..15

		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0,

		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0,

		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
)

func buildOctree() *Octree {

	oct := NewOctree()
	oct.Dim(4, 4, 4)
	oct.data = data

	for z := int32(0); z < 4; z++ {
		for y := int32(0); y < 4; y++ {
			for x := int32(0); x < 4; x++ {
				oct.Set(x, y, z, 1)

				if x == 0 && y == 0 && z == 0 {
					oct.Set(0, 0, 0, 6)
				}

				if x == 3 && y == 3 && z == 3 {
					oct.Set(3, 3, 3, 7)
				}
			}
		}
	}

	return oct
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
}

func TestOctree(t *testing.T) {

	oct := buildOctree()

	for i := 0; i < len(oct.Index) / 9; i++ {
		fmt.Printf("%d: %v %v\n", i,
			oct.Index[i*9:i*9 + 8],
			oct.Index[i*9 + 8:i*9 + 9])
	}

	if v, s := oct.Get(0, 0, 0); v != 6 || s != 1 {
		t.Errorf("test1 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(1, 1, 1); v != 1 || s != 1 {
		t.Errorf("test2 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(2, 2, 2); v != 1 || s != 1 {
		t.Errorf("test3 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(3, 3, 3); v != 7 || s != 1 {
		t.Errorf("test4 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(0, 0, 3); v != 1 || s != 2 {
		t.Errorf("test5 v=%d d=%d", v, s)
	}
	if v, s := oct.Get(0, 3, 0); v != 1 || s != 2 {
		t.Errorf("test6 v=%d d=%d", v, s)
	}
}

func TestGridTrace(t *testing.T) {

	g := Grid{data, 4, 4, 4}

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
