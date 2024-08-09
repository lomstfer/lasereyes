package msgfromclient

import (
	"wzrds/common/player"
)

type ConnectMe struct {
	Name string
}

type MoveInput struct {
	Input []player.Input
}
