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

func (a Vec3) Mul(f float32) Vec3 {
	return Vec3 { a.X*f, a.Y*f, a.Z*f }
}

func (a Vec3) Norm() float32 {
	return float32(math.Sqrt(float64(a.Dot(a))))
}

func (a Vec3) Normalize() Vec3 {
	return a.Mul(1.0/a.Norm())
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

	xs, ys, zs := q.X*s,  q.Y*s,  q.Z*s

	wx, wy, wz := q.R*xs, q.R*ys, q.R*zs
	xx, xy, xz := q.X*xs, q.X*ys, q.X*zs
	yy, yz, zz := q.Y*ys, q.Y*zs, q.Z*zs

	m[0], m[1], m[2] = 1.0 - yy - zz,	xy - wz,		xz + wy
	m[3], m[4], m[5] = xy + wz,			1.0 - xx - zz,	yz - wx
	m[6], m[7], m[8] = xz - wy,			yz + wx,		1.0 - xx - yy

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
