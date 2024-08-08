package internal

import (
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"
)

func SimulateInput(player *Player, input player.PlayerInput, deltaTime float64) {
	inputVec := vec2.NewVec2(0, 0)
	if input.Up {
		inputVec.Y -= 1
	}
	if input.Down {
		inputVec.Y += 1
	}
	if input.Left {
		inputVec.X -= 1
	}
	if input.Right {
		inputVec.X += 1
	}
	inputVec.Mul(deltaTime)
	inputVec.Mul(100)
	player.Position.Add(inputVec)
}
