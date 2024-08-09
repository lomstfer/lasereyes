package internal

import (
	"fmt"
	"wzrds/common/player"
)

type GameServer struct {
	Players          map[uint]*Player
	PlayersThatMoved map[uint]bool
	Tick             uint64
}

func NewGameServer() *GameServer {
	gs := &GameServer{}
	gs.Players = make(map[uint]*Player, 0)
	gs.PlayersThatMoved = make(map[uint]bool)

	return gs
}

func (gs *GameServer) AddPlayer(p player.CommonData) {
	gs.Players[p.Id] = &Player{Data: p}
}

func (gs *GameServer) RemovePlayer(id uint) {
	delete(gs.Players, id)
}

func (gs *GameServer) HandlePlayerInput(playerId uint, inputs []player.Input) {
	fmt.Println(len(inputs))
	p := gs.Players[playerId]
	for _, i := range inputs {
		if i.HasInput() {
			p.QueuedInputs = append(p.QueuedInputs, inputs...)
			gs.PlayersThatMoved[playerId] = true
			break
		}
	}
}

func (gs *GameServer) Simulate(deltaTime float64) {
	for _, p := range gs.Players {
		for len(p.QueuedInputs) > 0 {
			input := p.QueuedInputs[0]
			player.SimulateInput(&p.Data, input, deltaTime)
			p.LastAuthorizedInputId = input.Id
			p.QueuedInputs = p.QueuedInputs[1:]
		}
	}
	gs.Tick += 1
}
