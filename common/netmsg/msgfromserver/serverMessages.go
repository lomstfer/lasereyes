package msgfromserver

import (
	"wzrds/common/player"
)

type DisconnectedClient struct {
	Id uint
}

type DisconnectSelf struct{}

type AddSelfPlayer struct {
	Data player.CommonData
}

type AddOtherPlayer struct {
	Data player.CommonData
}

type UpdatePlayers struct {
	Players map[uint]player.Snapshot
}

type UpdateSelf struct {
	LastAuthorizedInputId uint32
	Snapshot              player.Snapshot
}
