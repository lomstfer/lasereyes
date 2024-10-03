package internal

import (
	"image/color"
	"math"
	"wzrds/client/internal/constants"
	"wzrds/common/pkg/vec2"

	"github.com/hajimehoshi/ebiten/v2"
)

type LaserBeam struct {
	OwnerId          uint
	TimeInstantiated float64
	TargetId         uint
}

func (lb *LaserBeam) Draw(screen *ebiten.Image, ownerPosition vec2.Vec2, targetPosition vec2.Vec2, laserBeamImage *ebiten.Image, timeNow float64, cameraTranslation vec2.Vec2) {
	size := 3.0
	{
		diff := targetPosition.Sub(ownerPosition)
		timeFrac := (timeNow - lb.TimeInstantiated) / constants.LaserBeamViewTime * 2
		length := diff.Length() * math.Min(timeFrac, 1)
		angle := math.Atan2(diff.Y, diff.X)
		{
			geo := ebiten.GeoM{}
			geo.Scale(length, size)
			geo.Rotate(angle)
			geo.Translate(ownerPosition.X+size/2.0*math.Sin(angle), ownerPosition.Y-size/2.0*math.Cos(angle))
			cs := ebiten.ColorScale{}
			cs.ScaleWithColor(color.NRGBA{R: 255, A: 255})
			geo.Translate(cameraTranslation.X, cameraTranslation.Y)
			screen.DrawImage(laserBeamImage, &ebiten.DrawImageOptions{GeoM: geo, ColorScale: cs})
		}
	}
}
