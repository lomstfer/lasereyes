package vec2

import "math"

type Vec2 struct {
	X float64
	Y float64
}

func NewVec2(x float64, y float64) Vec2 {
	return Vec2{x, y}
}

func (v *Vec2) Add(v2 Vec2) {
	v.X += v2.X
	v.Y += v2.Y
}

func (v *Vec2) Sub(v2 Vec2) {
	v.X -= v2.X
	v.Y -= v2.Y
}

func (v *Vec2) Mul(s float64) {
	v.X *= s
	v.Y *= s
}

func (v *Vec2) Div(s float64) {
	v.X /= s
	v.Y /= s
}

func (v *Vec2) Length() float64 {
	return float64(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v *Vec2) SqrLength() float64 {
	return v.X*v.X + v.Y*v.Y
}
