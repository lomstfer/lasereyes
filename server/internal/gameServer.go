package internal

import (
	"wzrds/common/player"
)

type GameServer struct {
	Players map[uint]*Player
}

func NewGameServer() *GameServer {
	gs := &GameServer{}
	gs.Players = make(map[uint]*Player, 0)

	return gs
}

func (gs *GameServer) AddPlayer(p player.PlayerSpawnData) {
	gs.Players[p.Id] = &Player{Id: p.Id, Name: p.Name, Position: p.Position, Velocity: p.Velocity}
}

func (gs *GameServer) RemovePlayer(id uint) {
	delete(gs.Players, id)
}

func (gs *GameServer) HandlePlayerInput(id uint, inputs []player.PlayerInput) {
	p := gs.Players[id]
	p.QueuedInputs = append(p.QueuedInputs, inputs...)
}

func (gs *GameServer) Simulate(deltaTime float64) {
	for _, p := range gs.Players {
		if len(p.QueuedInputs) > 0 {
			SimulateInput(p, p.QueuedInputs[0], deltaTime)
			p.QueuedInputs = p.QueuedInputs[1:]
		}
	}
}
