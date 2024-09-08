package internal

import (
	"fmt"
	"wzrds/common/commonconstants"
	"wzrds/common/netmsg/msgfromclient"
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"
	"wzrds/server/constants"
)

type GameServer struct {
	Players          map[uint]*Player
	PlayersThatMoved map[uint]bool
}

type PlayerInputOutcome struct {
	SomeoneWasShot bool
	ShooterId      uint
	WereShotIds    []uint
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
	delete(gs.PlayersThatMoved, id)
}

func (gs *GameServer) HandlePlayerUpdateFacingDir(playerId uint, dir vec2.Vec2) {
	gs.Players[playerId].Data.PupilDistDir01 = dir
	gs.PlayersThatMoved[playerId] = true
}

func (gs *GameServer) HandlePlayerInput(playerId uint, serverTimeNow float64, input msgfromclient.Input) *PlayerInputOutcome {
	p := gs.Players[playerId]
	if p.Data.Dead {
		return nil
	}
	for _, i := range input.Move.MoveInputs {
		if i.HasInput() {
			p.QueuedInputs = append(p.QueuedInputs, InputServerSide{Input: i})
		}
	}

	// sort.Slice(p.QueuedInputs, func(i, j int) bool {
	// 	return p.QueuedInputs[i].Id < p.QueuedInputs[j].Id
	// })

	if !input.Shoot.DidShoot {
		return nil
	}
	playerPositionCopy := p.Data.Position
	for _, i := range input.Move.MoveInputs {
		if i.HasInput() {
			player.SimulateInput(&playerPositionCopy, i, commonconstants.SimulationTickRate)
		}
	}
	// playerPositionCopy is now what the client saw

	// todo: check radius from shoot position to client player

	outcome := &PlayerInputOutcome{
		ShooterId:   playerId,
		WereShotIds: make([]uint, 0),
	}

	// rewind other players (constants.ServerBroadcastRate * 2 == interpolation rewind)
	timeOfClientView := input.Shoot.Time - commonconstants.ServerBroadcastRate*2
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
		pRectMin := pCop.Data.Position.Sub(vec2.NewVec2Both(commonconstants.PlayerSize / 2.0))
		pRectMax := pCop.Data.Position.Add(vec2.NewVec2Both(commonconstants.PlayerSize / 2.0))
		if shootPos.X >= pRectMin.X && shootPos.X <= pRectMax.X &&
			shootPos.Y >= pRectMin.Y && shootPos.Y <= pRectMax.Y {
			gs.PlayerWasShot(pCop.Data.Id, playerId)
			outcome.SomeoneWasShot = true
			outcome.WereShotIds = append(outcome.WereShotIds, pCop.Data.Id)
		}
	}

	return outcome
}

func (gs *GameServer) Simulate(deltaTime float64, serverTime float64) {
	for _, p := range gs.Players {
		if p.Data.Dead {
			continue
		}
		p.HistoryForRewind = append(p.HistoryForRewind, player.Snapshot{Time: serverTime, Position: p.Data.Position})
		{
			i := 0
			for serverTime-p.HistoryForRewind[i].Time > commonconstants.ServerBroadcastRate*2*10 {
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
	p := gs.Players[wasShot]
	p.Data.Health -= constants.Damage
	if p.Data.Health <= 0 {
		p.Data.Dead = true
	}
}
