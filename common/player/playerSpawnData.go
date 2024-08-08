package player

import "wzrds/common/pkg/vec2"

type PlayerSpawnData struct {
	Id       uint
	Name     string
	Position vec2.Vec2
	Velocity vec2.Vec2
}
