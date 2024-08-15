package internal

import (
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"
)

type Player struct {
	Data                  player.CommonData
	QueuedInputs          []InputServerSide
	LastAuthorizedInputId uint32
	LastPos               vec2.Vec2
}

type InputServerSide struct {
	Input player.Input
}
