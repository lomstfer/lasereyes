package internal

import (
	"fmt"
	"wzrds/common/constants"
	"wzrds/common/netmsg/msgfromclient"
	"wzrds/common/pkg/vec2"
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

func (gs *GameServer) HandlePlayerInput(playerId uint, serverTimeNow float64, input msgfromclient.Input) {
	p := gs.Players[playerId]
	for _, i := range input.Move.MoveInputs {
		if i.HasInput() {
			p.QueuedInputs = append(p.QueuedInputs, InputServerSide{Input: i})
		}
	}

	// sort.Slice(p.QueuedInputs, func(i, j int) bool {
	// 	return p.QueuedInputs[i].Id < p.QueuedInputs[j].Id
	// })

	if !input.Shoot.DidShoot {
		return
	}
	playerPositionCopy := p.Data.Position
	for _, i := range input.Move.MoveInputs {
		if i.HasInput() {
			player.SimulateInput(&playerPositionCopy, i, constants.SimulationTickRate)
		}
	}
	// playerPositionCopy is now what the client saw

	// todo: check radius from shoot position to client player

	// rewind other players (constants.ServerBroadcastRate * 2 == interpolation rewind)
	timeOfClientView := input.Shoot.Time - constants.ServerBroadcastRate*2
	for _, pOrg := range gs.Players {
		// dont let you shoot yourself
		if pOrg.Data.Id == playerId {
			continue
		}

		pCop := &PlayerCopyForRewind{Data: pOrg.Data, HistoryForRewind: pOrg.HistoryForRewind}
		for i := len(pCop.HistoryForRewind) - 2; i >= 0; i-- {
			h0 := pCop.HistoryForRewind[i]
			h := pCop.HistoryForRewind[i+1]
			if h0.Time <= timeOfClientView && timeOfClientView <= h.Time {
				p0 := h0.Position
				p1 := h.Position
				t0 := h0.Time
				t1 := h.Time
				pCop.Data.Position = vec2.Lerp(p0, p1, (timeOfClientView-t0)/(t1-t0))
				break
			}
		}

		shootPos := input.Shoot.Position
		pPos := pCop.Data.Position
		if shootPos.X >= pPos.X && shootPos.X <= pPos.X+constants.PlayerWidthAndHeight &&
			shootPos.Y >= pPos.Y && shootPos.Y <= pPos.Y+constants.PlayerWidthAndHeight {
			gs.PlayerWasShot(pCop.Data.Id, playerId)
		}
	}
}

func (gs *GameServer) Simulate(deltaTime float64, serverTime float64) {
	for _, p := range gs.Players {
		p.HistoryForRewind = append(p.HistoryForRewind, player.Snapshot{Time: serverTime, Position: p.Data.Position})
		{
			i := 0
			for serverTime-p.HistoryForRewind[i].Time > constants.ServerBroadcastRate*2*10 {
				i += 1
			}
			p.HistoryForRewind = p.HistoryForRewind[i:]
		}
		if len(p.QueuedInputs) > 0 {
			input := p.QueuedInputs[0]
			player.SimulateInput(&p.Data.Position, input.Input, deltaTime)
			p.LastAuthorizedInputId = input.Input.Id
			p.QueuedInputs = p.QueuedInputs[1:]
			gs.PlayersThatMoved[p.Data.Id] = true

		}
	}
}

func (gs *GameServer) PlayerWasShot(wasShot uint, shooter uint) {
	fmt.Println(wasShot, "was shot by:", shooter)
}
