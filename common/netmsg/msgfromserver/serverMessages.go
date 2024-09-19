package msgfromserver

import (
	"wzrds/common/netmsg/msgfromclient"
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

type RemoveOtherPlayer struct {
	Id uint
}

type UpdatePlayers struct {
	IdsToSnapshots map[uint]player.Snapshot
}

type UpdateSelf struct {
	LastAuthorizedInputId uint32
	Snapshot              player.Snapshot
}

type TimeAnswer struct {
	Request      msgfromclient.TimeRequest
	TimeReceived float64
}

type PlayerTakeDamage struct {
	PlayerId        uint
	Damage          float32
	CausingDamageId uint
}
