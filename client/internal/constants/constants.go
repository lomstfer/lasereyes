package constants

import (
	"wzrds/common/commonconstants"
)

const WaitForCleanCloseTime = 2
const SendInputRate = 1.0 / 20.0
const TimesToSyncClock = 10
const ServerTimeSyncDeltaMS = 100

const PupilSize = commonconstants.PixelScale * 4
const PupilMaxDistanceFromEye = commonconstants.PlayerSize / 7
const MouseDistanceFromPupilForMax = 50

const LaserBeamViewTime = 0.2

const CameraSpeed = 5
