package common

import "time"

type FixedCallback struct {
	DeltaSeconds   float64
	Accumulator    float64
	lastUpdateTime time.Time
}

func NewFixedCallback(deltaSeconds float64) *FixedCallback {
	fc := &FixedCallback{}
	fc.DeltaSeconds = deltaSeconds
	fc.lastUpdateTime = time.Now()
	return fc
}

// Calls the callback function if enough time has been accumulated
func (fc *FixedCallback) Update(callback func()) {
	now := time.Now()
	fc.Accumulator += now.Sub(fc.lastUpdateTime).Seconds()
	fc.lastUpdateTime = now
	for fc.Accumulator >= fc.DeltaSeconds {
		fc.Accumulator -= fc.DeltaSeconds
		callback()
	}
}
