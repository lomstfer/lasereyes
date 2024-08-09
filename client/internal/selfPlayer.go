package internal

import (
	"wzrds/common/player"
)

type SelfPlayer struct {
	Data               player.CommonData
	InputsToSend       []player.Input
	unauthorizedInputs []player.Input
	inputIdCounter     uint32
}

func NewSelfPlayer(data player.CommonData) *SelfPlayer {
	sp := &SelfPlayer{Data: data}
	sp.InputsToSend = make([]player.Input, 0)
	sp.unauthorizedInputs = make([]player.Input, 0)
	return sp
}

func (sp *SelfPlayer) AddInput(input player.Input) {
	sp.inputIdCounter += 1
	input.Id = sp.inputIdCounter
	sp.unauthorizedInputs = append(sp.unauthorizedInputs, input)
	sp.InputsToSend = append(sp.InputsToSend, input)

	player.SimulateInput(&sp.Data, input, 1.0/60.0)
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
		player.SimulateInput(&sp.Data, inp, 1.0/60.0)
	}
}
