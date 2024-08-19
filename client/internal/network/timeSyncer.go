package network

import "wzrds/common/commonutils"

type TimeSyncer struct {
	FinishedSync       bool
	minimumSyncAnswers int

	syncAnswersReceived int

	smallestLatency float64
	timeDiff        float64
}

func NewTimeSyncer(minimumSyncAnswers int) *TimeSyncer {
	ts := &TimeSyncer{}
	ts.minimumSyncAnswers = minimumSyncAnswers
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
	if ts.syncAnswersReceived >= ts.minimumSyncAnswers {
		ts.FinishedSync = true
	}
}

func (ts *TimeSyncer) ServerTime() float64 {
	return commonutils.GetUnixTimeAsFloat() - ts.timeDiff
}
