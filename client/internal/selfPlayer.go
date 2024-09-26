package internal

import (
	"wzrds/client/internal/constants"
	"wzrds/common/commonconstants"
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"
)

type SelfPlayer struct {
	Data               player.CommonData
	InputsToSend       []player.MoveInput
	OldPosition        vec2.Vec2
	SmoothedPosition   vec2.Vec2
	unauthorizedInputs []player.MoveInput
	inputIdCounter     uint32
}

func NewSelfPlayer(data player.CommonData) *SelfPlayer {
	sp := &SelfPlayer{Data: data}
	sp.OldPosition = sp.Data.Position
	sp.SmoothedPosition = sp.Data.Position
	sp.InputsToSend = make([]player.MoveInput, 0)
	sp.unauthorizedInputs = make([]player.MoveInput, 0)
	return sp
}

func (sp *SelfPlayer) CalculateFacingVec(mousePositionWorld vec2.Vec2) {
	rel := mousePositionWorld.Sub(sp.SmoothedPosition)
	sp.Data.PupilDistDir01 = rel.Div(constants.MouseDistanceFromPupilForMax).LengthClamped(0, 1)
}

func (sp *SelfPlayer) CheckMoveInput(inputVec vec2.Vec2) {
	input := player.MoveInput{
		Up:    inputVec.Y == -1,
		Down:  inputVec.Y == 1,
		Left:  inputVec.X == -1,
		Right: inputVec.X == 1,
	}
	sp.OldPosition = sp.Data.Position
	if input.HasInput() {
		sp.AddInput(input)
	}
}

func (sp *SelfPlayer) AddInput(input player.MoveInput) {
	sp.inputIdCounter += 1
	input.Id = sp.inputIdCounter
	sp.unauthorizedInputs = append(sp.unauthorizedInputs, input)
	sp.InputsToSend = append(sp.InputsToSend, input)

	player.SimulateInput(&sp.Data.Position, input, commonconstants.SimulationTickRate)
}

func (sp *SelfPlayer) OnSendInputs() {
	sp.InputsToSend = sp.InputsToSend[:0]
}

// client accumulates inputs and updates its position, stores the inputs that have only been simulated on client
// client sends inputs
// server updates clients position and sends last updated id and position
// client checks its array of inputs only simulated on client and removes the ones less than or equal to the id from the server
// client sets its position to the server position

func (sp *SelfPlayer) HandleServerUpdate(lastAuthorizedInputId uint32, snapshot player.Snapshot) {
	for i, inp := range sp.unauthorizedInputs {
		if inp.Id == lastAuthorizedInputId {
			sp.unauthorizedInputs = sp.unauthorizedInputs[i+1:]
			break
		} else if inp.Id > lastAuthorizedInputId {
			sp.unauthorizedInputs = sp.unauthorizedInputs[i:]
			break
		}
	}
	authorizedPosition := snapshot.Position
	sp.Data.Position = authorizedPosition
	for _, inp := range sp.unauthorizedInputs {
		player.SimulateInput(&sp.Data.Position, inp, commonconstants.SimulationTickRate)
	}

}

func (sp *SelfPlayer) UpdateSmoothPosition(alpha float64) {
	p0 := sp.OldPosition
	p1 := sp.Data.Position
	// todo: check if delta position is greater than movement and snap if so
	sp.SmoothedPosition = vec2.Lerp(p0, p1, alpha)
}
