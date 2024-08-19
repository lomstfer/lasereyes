package internal

import (
	"wzrds/common/constants"
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"
)

type SelfPlayer struct {
	Data               player.CommonData
	InputsToSend       []player.Input
	OldPosition        vec2.Vec2
	RenderPosition     vec2.Vec2
	unauthorizedInputs []player.Input
	inputIdCounter     uint32
}

func NewSelfPlayer(data player.CommonData) *SelfPlayer {
	sp := &SelfPlayer{Data: data}
	sp.OldPosition = sp.Data.Position
	sp.RenderPosition = sp.Data.Position
	sp.InputsToSend = make([]player.Input, 0)
	sp.unauthorizedInputs = make([]player.Input, 0)
	return sp
}

func (sp *SelfPlayer) CheckMoveInput(inputVec vec2.Vec2, localTime float64) {
	input := player.Input{
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

func (sp *SelfPlayer) AddInput(input player.Input) {
	sp.inputIdCounter += 1
	input.Id = sp.inputIdCounter
	sp.unauthorizedInputs = append(sp.unauthorizedInputs, input)
	sp.InputsToSend = append(sp.InputsToSend, input)

	player.SimulateInput(&sp.Data, input, constants.SimulationTickRate)
}

func (sp *SelfPlayer) OnSendInputs() {
	sp.InputsToSend = sp.InputsToSend[:0]
}

func (sp *SelfPlayer) HandleServerUpdate(lastAuthorizedInputId uint32, snapshot player.Snapshot) {
	for i, inp := range sp.unauthorizedInputs {
		if inp.Id >= lastAuthorizedInputId {
			sp.unauthorizedInputs = sp.unauthorizedInputs[i+1:]
			break
		}
	}
	authorizedPosition := snapshot.Position
	sp.Data.Position = authorizedPosition
	for _, inp := range sp.unauthorizedInputs {
		player.SimulateInput(&sp.Data, inp, constants.SimulationTickRate)
	}
}

func (sp *SelfPlayer) UpdateRenderPosition(alpha float64) {
	p0 := sp.OldPosition
	p1 := sp.Data.Position
	// todo: check if delta position is greater than movement and snap if so
	sp.RenderPosition = vec2.Lerp(p0, p1, alpha)
}
