package vec2

import (
	"math"
)

type Vec2 struct {
	X float64
	Y float64
}

func NewVec2(x float64, y float64) Vec2 {
	return Vec2{x, y}
}

func NewVec2Both(xy float64) Vec2 {
	return Vec2{xy, xy}
}

func (v Vec2) Add(v2 Vec2) Vec2 {
	v.X += v2.X
	v.Y += v2.Y
	return v
}

func (v Vec2) Sub(v2 Vec2) Vec2 {
	v.X -= v2.X
	v.Y -= v2.Y
	return v
}

func (v Vec2) Mul(s float64) Vec2 {
	v.X *= s
	v.Y *= s
	return v
}

func (v Vec2) Div(s float64) Vec2 {
	v.X /= s
	v.Y /= s
	return v
}

func (v Vec2) Normalized() Vec2 {
	length := v.Length()
	if length == 0 {
		return v
	}
	v.X /= length
	v.Y /= length
	return v
}

func (v Vec2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vec2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func Lerp(v1 Vec2, v2 Vec2, t float64) Vec2 {
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	return Vec2{
		v1.X + (v2.X-v1.X)*t,
		v1.Y + (v2.Y-v1.Y)*t,
	}
}

func (v Vec2) Clamped(vMin Vec2, vMax Vec2) Vec2 {
	v.X = math.Min(math.Max(v.X, vMin.X), vMax.X)
	v.Y = math.Min(math.Max(v.Y, vMin.Y), vMax.Y)
	return v
}

func (v Vec2) LengthClamped(min float64, max float64) Vec2 {
	if v.Length() < min {
		return v.Normalized().Mul(min)
	}
	if v.Length() > max {
		return v.Normalized().Mul(max)
	}
	return v
}
