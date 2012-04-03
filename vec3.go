package glvox

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float32
}

type Mat3 [9]float32

type Quat struct {
	R float32
	Vec3
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

func (a Vec3) Mulv(b Vec3) Vec3 {
	return Vec3 { a.X*b.X, a.Y*b.Y, a.Z*b.Z }
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

func (a Vec3) Floor() Vec3 {
	return Vec3 {
		float32(math.Floor(float64(a.X))),
		float32(math.Floor(float64(a.Y))),
		float32(math.Floor(float64(a.Z))) }
}

func (a Vec3) Ceil() Vec3 {
	return Vec3 {
		float32(math.Ceil(float64(a.X))),
		float32(math.Ceil(float64(a.Y))),
		float32(math.Ceil(float64(a.Z))) }
}

func (a Vec3) Nexttoward(b Vec3) Vec3 {
	return Vec3 {
		float32(math.Nextafter(float64(a.X), float64(b.X))),
		float32(math.Nextafter(float64(a.Y), float64(b.Y))),
		float32(math.Nextafter(float64(a.Z), float64(b.Z))) }
}

func (a Vec3) Clamp(l, h Vec3) Vec3 {

	x := a.X
	if x < l.X { x = l.X }; if x > h.X { x = h.X }

	y := a.Y
	if y < l.Y { y = l.Y }; if y > h.Y { y = h.Y }

	z := a.Z
	if z < l.Z { z = l.Z }; if z > h.Z { z = h.Z }

	return Vec3 { x, y, z }
}

func NewQuat(r float32, v Vec3) (q Quat) {

	cosR := float32(math.Cos(float64(r/2)))
	sinR := float32(math.Sin(float64(r/2)))

	q = Quat { cosR, v.Mul(sinR) }
	return q
}

func (q Quat) Mat3() (m Mat3) {

	n := q.R*q.R + q.X*q.X + q.Y*q.Y + q.Z*q.Z
	s := float32(0.0); if n > 0.0 { s = 2.0 / n }

	x2, y2, z2 := q.X*s,  q.Y*s,  q.Z*s

	rx, ry, rz := q.R*x2, q.R*y2, q.R*z2
	xx, xy, xz := q.X*x2, q.X*y2, q.X*z2
	yy, yz, zz := q.Y*y2, q.Y*z2, q.Z*z2

	m[0], m[1], m[2] = 1.0 - yy - zz,	xy - rz,		xz + ry
	m[3], m[4], m[5] = xy + rz,			1.0 - xx - zz,	yz - rx
	m[6], m[7], m[8] = xz - ry,			yz + rx,		1.0 - xx - yy

	return m
}

func (q Quat) Mul(r Quat) (s Quat) {
	s.R = q.R*r.R - q.X*r.X - q.Y*r.Y - q.Z*r.Z;
	s.X = q.R*r.X + r.R*q.X + q.Y*r.Z - q.Z*r.Y;
	s.Y = q.R*r.Y + r.R*q.Y + q.Z*r.X - q.X*r.Z;
	s.Z = q.R*r.Z + r.R*q.Z + q.X*r.Y - q.Y*r.X;
	return s;
}

func (m Mat3) Mul(v Vec3) (w Vec3) {
	w = Vec3 {
		m[0]*v.X + m[1]*v.Y + m[2]*v.Z,
		m[3]*v.X + m[4]*v.Y + m[5]*v.Z,
		m[6]*v.X + m[7]*v.Y + m[8]*v.Z,
	}
	return w
}
