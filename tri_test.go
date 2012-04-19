package glvox_test

import (
	"github.com/shogg/glvox"
	"testing"
)

var (
	tri = &glvox.Tri{
		{0.0, 0.0, 0.0},
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
	}

	pp = glvox.Vec3{1.0, 1.0, 1.0}.Normalize()
	pm = glvox.Vec3{1.5, -1.0, 1.0}.Normalize()
	mp = glvox.Vec3{-1.0, 1.5, 1.0}.Normalize()
	mm = glvox.Vec3{-1.0, -1.0, 1.0}.Normalize()

	halfHypo = tri[1].Plus(tri[2].Minus(tri[1]).Mul(0.5))

	points = []glvox.Vec3{
		{0.3, 0.3, -1.0},  // region 0
		{0.3, 0.3, 1.0},   // region 0
		tri[0].Plus(mm),   // region 4
		{0.0, -1.0, 0.0},  // region 5
		{0.5, -1.0, 0.0},  // region 5
		{1.0, -1.0, 0.0},  // region 5
		tri[1].Plus(pm),   // region 6
		{2.0, 0.0, 0.0},   // region 1
		tri[1].Plus(pp),   // region 1
		halfHypo.Plus(pp), // region 1
		tri[2].Plus(pp),   // region 1
		{0.0, 2.0, 0.0},   // region 1
		tri[2].Plus(mp),   // region 2
		{-1.0, 1.0, 0.0},  // region 3
		{-1.0, 0.5, 0.0},  // region 3
		{-1.0, 0.0, 0.0},  // region 3
	}
)

func TestSqrDistance(t *testing.T) {

	eps := float32(0.000001)

	for _, p := range points {
		dist := tri.SqrDistance(p)
		if dist < 1.0-eps || dist > 1.0+eps {
			t.Error("distance 1.0 expected, was", dist)
		}
	}
}
