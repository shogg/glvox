package glvox

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float32
}

func (a Vec3) Cross(b Vec3) Vec3 {
	return Vec3 {
			a.Y*b.Z - a.Z*b.Y,
			a.Z*b.X - a.X*b.Z,
			a.X*b.Y - a.Y*b.X,
		}
}

func (a Vec3) Dot(b Vec3) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vec3) Plus(b Vec3) Vec3 {
	return Vec3 { a.X+b.X, a.Y+b.Y, a.Z+b.Z }
}

func (a Vec3) Minus(b Vec3) Vec3 {
	return Vec3 { a.X-b.X, a.Y-b.Y, a.Z-b.Z }
}

func (a Vec3) Mul(f float32) Vec3 {
	return Vec3 { a.X*f, a.Y*f, a.Z*f }
}

func (a Vec3) Norm() float32 {
	return float32(math.Sqrt(float64(a.Dot(a))))
}

func (a Vec3) Normalize() Vec3 {
	return a.Mul(1.0/a.Norm())
}
