package player

import (
	"wzrds/common/pkg/vec2"
)

func SimulateInput(playerDataPosition *vec2.Vec2, input MoveInput, deltaTime float64) {
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
	inputVec = inputVec. /* .Normalized() */ Mul(40 * deltaTime)
	*playerDataPosition = playerDataPosition.Add(inputVec)
}
