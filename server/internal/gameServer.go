package internal

import (
	"wzrds/common/player"
)

type GameServer struct {
	Players          map[uint]*Player
	PlayersThatMoved map[uint]bool
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

func (gs *GameServer) HandlePlayerInput(playerId uint, inputs []player.Input, serverTime float64) {
	p := gs.Players[playerId]
	for _, i := range inputs {
		if i.HasInput() {
			p.QueuedInputs = append(p.QueuedInputs, InputServerSide{Input: i})
			gs.PlayersThatMoved[playerId] = true
		}
	}

	// sort.Slice(p.QueuedInputs, func(i, j int) bool {
	// 	return p.QueuedInputs[i].Id < p.QueuedInputs[j].Id
	// })
}

func (gs *GameServer) Simulate(deltaTime float64, serverTime float64) {
	for _, p := range gs.Players {
		if len(p.QueuedInputs) > 0 {
			input := p.QueuedInputs[0]
			player.SimulateInput(&p.Data, input.Input, deltaTime)
			p.LastAuthorizedInputId = input.Input.Id
			p.QueuedInputs = p.QueuedInputs[1:]
		}
	}
}
