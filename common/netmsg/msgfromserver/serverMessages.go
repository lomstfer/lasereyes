package msgfromserver

import (
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"
)

type DisconnectedClient struct {
	Id uint
}

type DisconnectSelf struct{}

type AddSelfPlayer struct {
	Data player.PlayerSpawnData
}

type AddOtherPlayer struct {
	Data player.PlayerSpawnData
}

type UpdatePlayers struct {
	Players map[uint]vec2.Vec2
}
