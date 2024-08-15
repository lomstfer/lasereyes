package internal

import (
	"sort"
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"
)

type Player struct {
	Data               player.CommonData
	SnapshotsForInterp []player.Snapshot
}

func (p *Player) LerpBetweenSnapshots(syncedServerTime float64) {
	sort.Slice(p.SnapshotsForInterp, func(i, j int) bool {
		return p.SnapshotsForInterp[i].Time < p.SnapshotsForInterp[j].Time
	})

	renderingTime := syncedServerTime - 0.2

	if len(p.SnapshotsForInterp) < 2 {
		if len(p.SnapshotsForInterp) == 1 {
			p.SnapshotsForInterp[0].Time = renderingTime
		}
		return
	}

	// remove old snapshots
	for len(p.SnapshotsForInterp) >= 2 && p.SnapshotsForInterp[1].Time < renderingTime {
		p.SnapshotsForInterp = p.SnapshotsForInterp[1:]
	}
	if len(p.SnapshotsForInterp) < 2 {
		return
	}

	s0 := p.SnapshotsForInterp[0]
	s1 := p.SnapshotsForInterp[1]
	p0 := s0.Position
	p1 := s1.Position
	t0 := s0.Time
	t1 := s1.Time

	inBetween := t0 < renderingTime && renderingTime < t1
	if !inBetween {
		return
	}

	t := (renderingTime - t0) / (t1 - t0)

	p.Data.Position = vec2.Lerp(p0, p1, t)
}
