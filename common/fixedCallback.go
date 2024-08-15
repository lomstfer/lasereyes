package common

import "time"

type FixedCallback struct {
	deltaSeconds   float64
	accumulator    float64
	lastUpdateTime time.Time
}

func NewFixedCallback(deltaSeconds float64) *FixedCallback {
	fc := &FixedCallback{}
	fc.deltaSeconds = deltaSeconds
	fc.lastUpdateTime = time.Now()
	return fc
}

// Calls the callback function if enough time has been accumulated
func (fc *FixedCallback) Update(callback func()) {
	now := time.Now()
	fc.accumulator += now.Sub(fc.lastUpdateTime).Seconds()
	fc.lastUpdateTime = now
	for fc.accumulator >= fc.deltaSeconds {
		fc.accumulator -= fc.deltaSeconds
		callback()
	}
}
