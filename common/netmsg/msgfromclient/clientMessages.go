package msgfromclient

import (
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"
)

type ConnectMe struct {
	Name string
}

type MoveInput struct {
	MoveInputs []player.MoveInput
}
type ShootInput struct {
	DidShoot bool
	Time     float64
	Position vec2.Vec2
}
type Input struct {
	Shoot ShootInput
	Move  MoveInput
}

type TimeRequest struct {
	TimeSent float64
}
