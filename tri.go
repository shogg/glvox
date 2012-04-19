package glvox

type Tri [3]Vec3

// http://www.geometrictools.com/Documentation/DistancePoint3Triangle3.pdf
func (tri *Tri) SqrDistance(p Vec3) float32 {

	B := tri[0]
	E0 := tri[1].Minus(B)
	E1 := tri[2].Minus(B)
	D := B.Minus(p)

	a := E0.Dot(E0)
	b := E0.Dot(E1)
	c := E1.Dot(E1)
	d := E0.Dot(D)
	e := E1.Dot(D)
	f := D.Dot(D)

	det := a*c - b*b
	s := b*e - c*d
	t := b*d - a*e

	if s+t <= det {
		if s < 0.0 {
			if t < 0.0 {
				// region 4

				// Grad(Q) = 2(as+bt+d,bs+ct+e)
				// (1,0)*Grad(Q(0,0)) = (1,0)*(d,e) = d
				// (0,1)*Grad(Q(0,0)) = (0,1)*(d,e) = e
				// min on edge t=0 if (0,1)*Grad(Q(0,0)) < 0 )
				// min on edge s=0 otherwise

				if d < 0.0 {
					t = 0.0
					if -d >= a {
						s = 1.0
					} else {
						s = -d / a
					}
				} else {
					s = 0.0
					if e >= 0.0 {
						t = 0.0
					} else if -e >= c {
						t = 1.0
					} else {
						t = -e / c
					}
				}
			} else {
				// region 3

				// F(t) = Q(0,t) = ct^2 + 2et + f
				// F’(t)/2 = ct+e
				// F’(T) = 0 when T = -e/c

				s = 0.0
				if e >= 0.0 {
					t = 0.0
				} else if -e >= c {
					t = 1.0
				} else {
					t = -e / c
				}
			}
		} else if t < 0.0 {
			// region 5

			t = 0.0
			if d >= 0.0 {
				s = 0.0
			} else if -d >= a {
				s = 1.0
			} else {
				s = -d / a
			}
		} else {
			// region 0

			invDet := 1.0 / det
			s *= invDet
			t *= invDet
		}
	} else {
		if s < 0.0 {
			// region 2

			// Grad(Q) = 2(as+bt+d,bs+ct+e)
			// (0,-1)*Grad(Q(0,1)) = (0,-1)*(b+d,c+e) = -(c+e)
			// (1,-1)*Grad(Q(0,1)) = (1,-1)*(b+d,c+e) = (b+d)-(c+e)
			// min on edge s+t=1 if (1,-1)*Grad(Q(0,1)) < 0 )
			// min on edge s=0 otherwise

			tmp0 := b + d
			tmp1 := c + e
			if tmp1 > tmp0 {
				numer := tmp1 - tmp0
				denom := a - 2.0*b + c
				if numer >= denom {
					s = 1.0
				} else {
					s = numer / denom
				}
				t = 1.0 - s
			} else {
				s = 0.0
				if tmp1 <= 0.0 {
					t = 1.0
				} else if e >= 0.0 {
					t = 0.0
				} else {
					t = -e / c
				}
			}
		} else if t < 0.0 {
			// region 6

			// Grad(Q) = 2(as+bt+d,bs+ct+e)
			// (-1,0)*Grad(Q(1,0)) = (-1,0)*(a+d,b+e) = -(a+d)
			// (-1,1)*Grad(Q(1,0)) = (-1,1)*(a+d,b+e) = -(a+d)+(b+e)
			// min on edge s+t=1 if (1,-1)*Grad(Q(0,1)) < 0 )
			// min on edge s=0 otherwise

			tmp0 := a + d
			tmp1 := b + e
			if tmp1 > tmp0 {
				numer := tmp1 + tmp0
				denom := a - 2.0*b + c
				s = 1.0
				if numer >= denom {
					s = numer / denom
				}
				t = 1.0 - s
			} else {
				t = 0.0
				if tmp1 <= 0.0 {
					s = 1.0
				} else if e >= 0.0 {
					s = 0.0
				} else {
					s = -e / c
				}
			}
		} else {
			// region 1

			// F(s) = Q(s,1-s) = (a-2b+c)s^2 + 2(b-c+d-e)s + (c+2e+f)
			// F’(s)/2 = (a-2b+c)s + (b-c+d-e)
			// F’(S) = 0 when S = (c+e-b-d)/(a-2b+c)
			// a-2b+c = |E0-E1|^2 > 0, so only sign of c+e-b-d need be considered

			numer := c + e - b - d
			if numer <= 0.0 {
				s = 0.0
			} else {
				denom := a - 2.0*b + c
				if numer >= denom {
					s = 1.0
				} else {
					s = numer / denom
				}
			}
			t = 1.0 - s
		}
	}

	dist := a*s*s + 2.0*b*s*t + c*t*t + 2.0*d*s + 2.0*e*t + f
	return dist
}
