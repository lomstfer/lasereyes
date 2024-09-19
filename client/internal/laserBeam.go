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
	TargetPosition   vec2.Vec2
}

func (lb *LaserBeam) Draw(screen *ebiten.Image, ownerPosition vec2.Vec2, laserBeamImage *ebiten.Image, timeNow float64) {
	size := 3.0
	{
		diff := lb.TargetPosition.Sub(ownerPosition)
		timeFrac := (timeNow - lb.TimeInstantiated) / constants.LaserBeamViewTime * 2
		length := diff.Length() * math.Min(timeFrac, 1)
		angle := math.Atan2(diff.Y, diff.X)
		{
			geo := ebiten.GeoM{}
			geo.Scale(float64(length), float64(size))
			geo.Rotate(angle)
			geo.Translate(ownerPosition.X+size/2.0*math.Sin(angle), ownerPosition.Y-size/2.0*math.Cos(angle))
			cs := ebiten.ColorScale{}
			cs.ScaleWithColor(color.NRGBA{255, 0, 0, 255})
			screen.DrawImage(laserBeamImage, &ebiten.DrawImageOptions{GeoM: geo, ColorScale: cs})
		}
	}
}
