package glvox_test

import (
	"github.com/shogg/glvox"
	"testing"
)

var (

	/*

		\ 2 |
		 \  |
		  \ |
		   \|
		   [2]
		    |\
		    | \
		    |  \   1
		 3  |   \
		    |    \
		    |  0  \
		    |      \
		---[0]-----[1]----
		    |        \ 6
		 4  |    5    \

	*/

	tri = &glvox.Tri{
		{0.0, 0.0, 0.0},
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
	}

	xm = glvox.Vec3{-1.0, 0.0, 0.0}
	xp = glvox.Vec3{1.0, 0.0, 0.0}
	ym = glvox.Vec3{0.0, -1.0, 0.0}
	yp = glvox.Vec3{0.0, 1.0, 0.0}
	zm = glvox.Vec3{0.0, 0.0, -1.0}
	zp = glvox.Vec3{0.0, 0.0, 1.0}

	pp = glvox.Vec3{1.0, 1.0, 0.0}
	pm = glvox.Vec3{1.5, -1.0, 0.0}
	mp = glvox.Vec3{-1.0, 1.5, 0.0}
	mm = glvox.Vec3{-1.0, -1.0, 0.0}

	points = []glvox.Vec3{
		{0.3, 0.3, 0.0}, // region 0
		tri[0],          // region 4
		tri[0],          // region 5
		tri[1],          // region 5
		tri[1],          // region 6
		tri[1],          // region 1
		tri[1],          // region 1
		tri[2],          // region 1
		tri[2],          // region 1
		tri[2],          // region 2
		tri[2],          // region 3
		tri[0],          // region 3
	}

	dirs = []glvox.Vec3{
		{0.0, 0.0, 0.0}, // region 0
		mm,              // region 4
		ym,              // region 5
		ym,              // region 5
		pm,              // region 6
		xp,              // region 1
		pp,              // region 1
		pp,              // region 1
		yp,              // region 1
		mp,              // region 2
		xm,              // region 3
		xm,              // region 3
	}
)

func TestSqrDistance(t *testing.T) {

	eps := float32(0.000001)

	for i, p := range points {
		dist := tri.SqrDistance(p.Plus(dirs[i].Plus(zp).Normalize()))
		if dist < 1.0-eps || dist > 1.0+eps {
			t.Error("distance 1.0 expected, was", dist)
		}
	}

	for i, p := range points {
		dist := tri.SqrDistance(p.Plus(dirs[i].Plus(zm).Normalize()))
		if dist < -1.0-eps || dist > -1.0+eps {
			t.Error("distance -1.0 expected, was", dist)
		}
	}
}
