package internal

import (
	"wzrds/common/player"
)

type Player struct {
	Data                  player.CommonData
	QueuedInputs          []InputServerSide
	LastAuthorizedInputId uint32
	HistoryForRewind      []player.Snapshot
}

type InputServerSide struct {
	Input player.MoveInput
}

type PlayerCopyForRewind struct {
	Data             player.CommonData
	HistoryForRewind []player.Snapshot
}
