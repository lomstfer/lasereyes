package internal

import (
	"sort"
	"wzrds/common/commonconstants"
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawPlayer(data player.CommonData, screen *ebiten.Image, playerImage *ebiten.Image) {
	geo := ebiten.GeoM{}
	geo.Scale(3, 3)
	colorScale := ebiten.ColorScale{}
	if data.Dead {
		colorScale.Scale(0, 0, 0, 1)
	}
	geo.Translate(data.Position.X, data.Position.Y)
	screen.DrawImage(playerImage, &ebiten.DrawImageOptions{GeoM: geo, ColorScale: colorScale})
}

func DrawPlayerHealthBar(data player.CommonData, screen *ebiten.Image, healthBarBg *ebiten.Image, healthBarFg *ebiten.Image) {
	width := 60
	height := 5
	x := data.Position.X + commonconstants.PlayerWidthAndHeight/2 - float64(width)/2
	y := data.Position.Y - 10
	{
		geo := ebiten.GeoM{}
		geo.Scale(float64(width), float64(height))
		geo.Translate(x, y)
		screen.DrawImage(healthBarBg, &ebiten.DrawImageOptions{GeoM: geo})
	}
	{
		geo := ebiten.GeoM{}
		healthFraction := float64(data.Health) / 100
		if healthFraction > 1 {
			healthFraction = 1
		} else if healthFraction < 0 {
			healthFraction = 0
		}
		geo.Scale(healthFraction*float64(width), float64(height))
		geo.Translate(x, y)
		screen.DrawImage(healthBarFg, &ebiten.DrawImageOptions{GeoM: geo})
	}
}

type Player struct {
	Data               player.CommonData
	SnapshotsForInterp []player.Snapshot
}

func (p *Player) LerpBetweenSnapshots(syncedServerTime float64) {
	sort.Slice(p.SnapshotsForInterp, func(i, j int) bool {
		return p.SnapshotsForInterp[i].Time < p.SnapshotsForInterp[j].Time
	})

	renderingTime := syncedServerTime - commonconstants.ServerBroadcastRate*2

	if len(p.SnapshotsForInterp) < 2 {
		// avoids jump when the player just starts moving
		// if len(p.SnapshotsForInterp) == 1 {
		// 	p.SnapshotsForInterp[0].Time = renderingTime
		// }
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

	// if inBetween := t0 <= renderingTime && renderingTime <= t1; !inBetween {
	// 	return
	// }

	t := (renderingTime - t0) / (t1 - t0)

	p.Data.Position = vec2.Lerp(p0, p1, t)
}
