package internal

import (
	"image/color"
	"sort"
	"wzrds/client/internal/constants"
	"wzrds/common/commonconstants"
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func DrawPlayer(data player.CommonData, screen *ebiten.Image, eyeImage *ebiten.Image, pupilImage *ebiten.Image, cameraTranslation vec2.Vec2) {
	geo := ebiten.GeoM{}
	geo.Scale(commonconstants.PixelScale, commonconstants.PixelScale)

	pos := data.Position.Sub(vec2.NewVec2Both(commonconstants.PlayerSize / 2.0))
	geo.Translate(pos.X, pos.Y)

	colorScale := ebiten.ColorScale{}
	if data.Dead {
		colorScale.ScaleWithColor(color.NRGBA{150, 100, 0, 255})
	}
	geo.Translate(cameraTranslation.X, cameraTranslation.Y)
	screen.DrawImage(eyeImage, &ebiten.DrawImageOptions{GeoM: geo, ColorScale: colorScale})

	drawPlayerPupil(data, screen, pupilImage, colorScale, cameraTranslation)
}

func GetPupilPos(data player.CommonData) vec2.Vec2 {
	pos := data.Position.Add(data.PupilDistDir01.Mul(constants.PupilMaxDistanceFromEye))
	return pos
}

func drawPlayerPupil(data player.CommonData, screen *ebiten.Image, pupilImage *ebiten.Image, colorScale ebiten.ColorScale, cameraTranslation vec2.Vec2) {
	geo := ebiten.GeoM{}
	geo.Scale(commonconstants.PixelScale, commonconstants.PixelScale)
	colorScale.ScaleWithColor(data.Color)

	drawPos := GetPupilPos(data).Sub(vec2.NewVec2Both(constants.PupilSize / 2.0))

	geo.Translate(drawPos.X, drawPos.Y)
	geo.Translate(cameraTranslation.X, cameraTranslation.Y)

	screen.DrawImage(pupilImage, &ebiten.DrawImageOptions{GeoM: geo, ColorScale: colorScale})
}

func DrawPlayerHealthBar(data player.CommonData, screen *ebiten.Image, healthBarBg *ebiten.Image, healthBarFg *ebiten.Image, cameraTranslation vec2.Vec2) {
	width := 60
	height := 5
	x := data.Position.X - float64(width)/2
	y := data.Position.Y - commonconstants.PlayerSize/2.0 - 10
	{
		geo := ebiten.GeoM{}
		geo.Scale(float64(width), float64(height))
		geo.Translate(x, y)
		geo.Translate(cameraTranslation.X, cameraTranslation.Y)
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
		geo.Translate(cameraTranslation.X, cameraTranslation.Y)
		screen.DrawImage(healthBarFg, &ebiten.DrawImageOptions{GeoM: geo})
	}
}

func DrawPlayerName(data player.CommonData, screen *ebiten.Image, f *text.GoTextFace, cameraTranslation vec2.Vec2) {
	textToDraw := data.Name
	width, _ := text.Measure(textToDraw, f, 0)
	geo := ebiten.GeoM{}
	geo.Translate(data.Position.X, data.Position.Y)
	geo.Translate(-width/2, commonconstants.PlayerSize/2)
	geo.Translate(cameraTranslation.X, cameraTranslation.Y)
	text.Draw(screen, textToDraw, f, &text.DrawOptions{DrawImageOptions: ebiten.DrawImageOptions{GeoM: geo}})
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
		if len(p.SnapshotsForInterp) == 1 {
			p.SnapshotsForInterp[0].Time = renderingTime
		}
		return
	}

	// remove old snapshots
	for len(p.SnapshotsForInterp) >= 2 && p.SnapshotsForInterp[1].Time < renderingTime {
		p.SnapshotsForInterp = p.SnapshotsForInterp[1:]
	}
	if len(p.SnapshotsForInterp) > 0 && renderingTime > p.SnapshotsForInterp[len(p.SnapshotsForInterp)-1].Time {
		p.Data.Position = p.SnapshotsForInterp[0].Position
		p.Data.PupilDistDir01 = p.SnapshotsForInterp[0].PupilDistDir01
		return
	}

	s0 := p.SnapshotsForInterp[0]
	s1 := p.SnapshotsForInterp[1]
	t0 := s0.Time
	t1 := s1.Time
	t := (renderingTime - t0) / (t1 - t0)
	{
		p0 := s0.Position
		p1 := s1.Position
		p.Data.Position = vec2.Lerp(p0, p1, t)
	}
	{
		d0 := s0.PupilDistDir01
		d1 := s1.PupilDistDir01
		p.Data.PupilDistDir01 = vec2.Lerp(d0, d1, t)
	}
}
