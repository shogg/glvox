package glvox

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
	q := NewQuat(a, c.Left)
	m := q.Mat3()
	c.Dir = m.Mul(c.Dir)
	c.Up = m.Mul(c.Up)
}

// Rotate left/right
func (c *Cam) Yaw(a float32) {
	q := NewQuat(a, c.Up)
	m := q.Mat3()
	c.Dir = m.Mul(c.Dir)
	c.Left = m.Mul(c.Left)
}

// Rotate z
func (c *Cam) Roll(a float32) {
	q := NewQuat(a, c.Dir)
	m := q.Mat3()
	c.Left = m.Mul(c.Left)
	c.Up = m.Mul(c.Up)
}
