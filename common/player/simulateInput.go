package player

import (
	"wzrds/common/pkg/vec2"
)

func SimulateInput(playerData *CommonData, input Input, deltaTime float64) {
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
	inputVec.Normalize()
	inputVec.Mul(200 * deltaTime)
	playerData.Position.Add(inputVec)
}
