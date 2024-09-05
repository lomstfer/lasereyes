package player

import "wzrds/common/pkg/vec2"

type CommonData struct {
	Id       uint
	Name     string
	Position vec2.Vec2
	Health   float32
	Dead     bool
}
