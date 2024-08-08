package internal

import (
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"
)

type Player struct {
	Id           uint
	Name         string
	Position     vec2.Vec2
	Velocity     vec2.Vec2
	QueuedInputs []player.PlayerInput
}
