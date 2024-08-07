package msgfromclient

import "wzrds/common"

type ConnectionInfo struct {
	Name string
}

type MoveInput struct {
	Input []common.PlayerInput
}
