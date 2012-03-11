package glvox

import (
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

	index = []int32 {
	// 4x4x4
	//  x  x  x  x  X  X  X  X
	//  y  y  Y  Y  y  y  Y  Y
	//  z  Z  z  Z  z  Z  z  Z
		1, 0, 0, 0, 0, 0, 0, 2, // 0

	// 2x2x2 xyz
		6, 0, 0, 0, 0, 0, 0, 0, // 1 * 8

	// 2x2x2 XYZ
		0, 0, 0, 0, 0, 0, 0, 7, // 2 * 8
	}
)

func TestOctreeTrace(t *testing.T) {

	oct := Octree{index, data, 4}

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

	oct := Octree{index, data, 4}
	if a, s := oct.find(0, 0, 0); a != 6.0 || s != 1 {
		t.Errorf("test1 a=%f d=%d", a, s)
	}
	if a, s := oct.find(1, 1, 1); a != 0.0 || s != 1 {
		t.Errorf("test2 a=%f d=%d", a, s)
	}
	if a, s := oct.find(2, 2, 2); a != 0.0 || s != 1 {
		t.Errorf("test3 a=%f d=%d", a, s)
	}
	if a, s := oct.find(3, 3, 3); a != 7.0 || s != 1 {
		t.Errorf("test4 a=%f d=%d", a, s)
	}
	if a, s := oct.find(0, 0, 3); a != 0.0 || s != 2 {
		t.Errorf("test5 a=%f d=%d", a, s)
	}
	if a, s := oct.find(0, 3, 0); a != 0.0 || s != 2 {
		t.Errorf("test5 a=%f d=%d", a, s)
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

func TestReadBinvox(t *testing.T) {

	voxels := NewOctree()
	err := ReadBinvox("skull.binvox", voxels)
	if err != nil {
		t.Error(err)
	}

	if voxels.WHD != 256 {
		t.Error("dimension 256 expected, was", voxels.WHD)
	}

	fmt.Println("index size", len(voxels.index))
}

// TODO - BenchmarkOctreeTrace, BenchmarkGridTrace
