package network

import "wzrds/common/utils"

type TimeSyncer struct {
	syncsToDo   int
	syncDeltaMS int

	syncAnswersReceived int

	smallestLatency float64
	timeDiff        float64

	FinishedSync bool
}

func NewTimeSyncer(syncsToDo int) *TimeSyncer {
	ts := &TimeSyncer{}
	ts.syncsToDo = syncsToDo
	ts.smallestLatency = -1

	return ts
}

func (ts *TimeSyncer) OnTimeAnswer(currentTime float64, timeSentFromClient float64, timeSentFromServer float64) {
	roundTrip := currentTime - timeSentFromClient
	oneWayLatency := roundTrip / 2.0
	predictedServerTime := timeSentFromServer + oneWayLatency

	if ts.smallestLatency == -1 || oneWayLatency < ts.smallestLatency {
		ts.smallestLatency = oneWayLatency
		ts.timeDiff = currentTime - predictedServerTime
	}

	ts.syncAnswersReceived += 1
	if ts.syncAnswersReceived >= ts.syncsToDo {
		ts.FinishedSync = true
	}
}

func (ts *TimeSyncer) GetServerTime() float64 {
	return utils.GetCurrentTimeAsFloat() - ts.timeDiff
}
