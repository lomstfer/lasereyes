package internal

import (
	"wzrds/common/player"
)

type Player struct {
	Data                  player.CommonData
	QueuedInputs          []player.Input
	LastAuthorizedInputId uint32
}
