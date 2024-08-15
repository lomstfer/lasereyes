package internal

import (
	"wzrds/common/player"
)

type Player struct {
	Data                  player.CommonData
	QueuedInputs          []InputServerSide
	LastAuthorizedInputId uint32
}

type InputServerSide struct {
	Input player.Input
}
