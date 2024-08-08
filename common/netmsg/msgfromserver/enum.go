package msgfromserver

type MsgType int

const (
	MsgTypeUndefined MsgType = iota
	MsgTypeDisconnectedClient
	MsgTypeAddSelfPlayer
	MsgTypeAddOtherPlayer
	MsgTypeUpdatePlayers
)
