package player

import (
	"time"
	"wzrds/common/pkg/vec2"
)

type PlayerSnapshot struct {
	Time     time.Time
	Position vec2.Vec2
}
