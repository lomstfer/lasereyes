package msgfromclient

type MsgType byte

const (
	MsgTypeUndefined MsgType = iota
	MsgTypeConnectMe
	MsgTypeInput
	MsgTypeTimeRequest
	MsgTypeUpdateFacingDirection
)
