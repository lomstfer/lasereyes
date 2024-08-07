package msgfromclient

type MsgType int

const (
	MsgTypeUndefined MsgType = iota
	MsgTypeConnectionInfo
	MsgTypeMoveInput
)
