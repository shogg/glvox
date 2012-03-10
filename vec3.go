package glvox

import (
	"math"
)

type vec3 struct {
	x, y, z float32
}

func (a vec3) Dot(b vec3) float32 {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

func (a vec3) Plus(b vec3) vec3 {
	return vec3 { a.x+b.x, a.y+b.y, a.z+b.z }
}

func (a vec3) Minus(b vec3) vec3 {
	return vec3 { a.x-b.x, a.y-b.y, a.z-b.z }
}

func (a vec3) Mul(f float32) vec3 {
	return vec3 { a.x*f, a.y*f, a.z*f }
}

func (a vec3) Norm() float32 {
	return float32(math.Sqrt(float64(a.Dot(a))))
}

func (a vec3) Normalize() vec3 {
	return a.Mul(1.0/a.Norm())
}
