package player

import (
	"image/color"
	"wzrds/common/pkg/vec2"
)

type CommonData struct {
	Id             uint
	Name           string
	Position       vec2.Vec2
	Health         float32
	Dead           bool
	PupilDistDir01 vec2.Vec2
	Color          color.NRGBA
}
