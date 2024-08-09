package internal

import (
	"wzrds/common/player"
)

type Player struct {
	Data               player.CommonData
	SnapshotsForInterp []player.Snapshot
}
