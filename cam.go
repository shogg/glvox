package glvox

import (
	"math"
)

type Cam struct {
	Pos, Dir, Up, Left Vec3
}

func NewCam() *Cam {
	c := new(Cam)
	c.Dir = Vec3 { 0.0, 0.0, -1.0 }
	c.Up = Vec3 { 0.0, 1.0, 0.0 }
	c.Left = Vec3 {-1.0, 0.0, 0.0 }
	return c
}

// forward/backward
func (c *Cam) Move(a float32) {
	c.Pos = c.Pos.Plus(c.Dir.Mul(a))
}

// left/right
func (c *Cam) Strafe(a float32) {
	strafe := c.Dir.Cross(c.Up)
	c.Pos = c.Pos.Plus(strafe.Mul(a))
}

// up/down
func (c *Cam) Updown(a float32) {
	c.Pos = c.Pos.Plus(c.Up.Mul(a))
}

// Rotate up/down
func (c *Cam) Pitch(a float32) {

	sinA := float32(math.Sin(float64(a)))
	cosA := float32(math.Cos(float64(a)))

	c.Dir = Vec3 {
		c.Dir.X,
		cosA*c.Dir.Y - sinA*c.Dir.Z,
		sinA*c.Dir.Y + cosA*c.Dir.Z,
	}

	c.Up = Vec3 {
		c.Up.X,
		cosA*c.Up.Y - sinA*c.Up.Z,
		sinA*c.Up.Y + cosA*c.Up.Z,
	}

	c.Left = Vec3 {
		c.Left.X,
		cosA*c.Left.Y - sinA*c.Left.Z,
		sinA*c.Left.Y + cosA*c.Left.Z,
	}
}

// Rotate left/right
func (c *Cam) Yaw(a float32) {

	sinA := float32(math.Sin(float64(a)))
	cosA := float32(math.Cos(float64(a)))

	c.Dir = Vec3 {
		sinA*c.Dir.Z + cosA*c.Dir.X,
		c.Dir.Y,
		cosA*c.Dir.Z - sinA*c.Dir.X,
	}

	c.Up = Vec3 {
		sinA*c.Up.Z + cosA*c.Up.X,
		c.Up.Y,
		cosA*c.Up.Z - sinA*c.Up.X,
	}

	c.Left = Vec3 {
		sinA*c.Left.Z + cosA*c.Left.X,
		c.Left.Y,
		cosA*c.Left.Z - sinA*c.Left.X,
	}
}

// Rotate z
func (c *Cam) Roll(a float32) {

	sinA := float32(math.Sin(float64(a)))
	cosA := float32(math.Cos(float64(a)))

	c.Dir = Vec3 {
		cosA*c.Dir.X - sinA*c.Dir.Y,
		sinA*c.Dir.X + cosA*c.Dir.Y,
		c.Dir.Z,
	}

	c.Up = Vec3 {
		cosA*c.Up.X - sinA*c.Up.Y,
		sinA*c.Up.X + cosA*c.Up.Y,
		c.Up.Z,
	}

	c.Left = Vec3 {
		cosA*c.Left.X - sinA*c.Left.Y,
		sinA*c.Left.X + cosA*c.Left.Y,
		c.Left.Z,
	}
}
