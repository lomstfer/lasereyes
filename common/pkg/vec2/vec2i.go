package vec2

import (
	"math"
)

type Vec2i struct {
	X int
	Y int
}

func NewVec2i(x int, y int) Vec2i {
	return Vec2i{x, y}
}

func NewVec2iBoth(xy int) Vec2i {
	return Vec2i{xy, xy}
}

func (v Vec2i) Add(v2 Vec2i) Vec2i {
	v.X += v2.X
	v.Y += v2.Y
	return v
}

func (v Vec2i) Sub(v2 Vec2i) Vec2i {
	v.X -= v2.X
	v.Y -= v2.Y
	return v
}

func (v Vec2i) Mul(s int) Vec2i {
	v.X *= s
	v.Y *= s
	return v
}

func (v Vec2i) Div(s int) Vec2i {
	v.X /= s
	v.Y /= s
	return v
}

func (v Vec2i) Normalized() Vec2 {
	length := v.Length()
	if length == 0 {
		return NewVec2(float64(v.X), float64(v.Y))
	}
	return NewVec2(float64(v.X)/length, float64(v.Y)/length)
}

func (v Vec2i) Length() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}

func (v Vec2i) LengthSquared() int {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec2i) Clamped(vMin Vec2i, vMax Vec2i) Vec2i {
	v.X = int(math.Min(math.Max(float64(v.X), float64(vMin.X)), float64(vMax.X)))
	v.Y = int(math.Min(math.Max(float64(v.Y), float64(vMin.Y)), float64(vMax.Y)))
	return v
}
