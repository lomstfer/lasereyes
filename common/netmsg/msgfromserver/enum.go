package msgfromserver

type MsgType byte

const (
	MsgTypeUndefined MsgType = iota
	MsgTypeDisconnectedClient
	MsgTypeAddSelfPlayer
	MsgTypeAddOtherPlayer
	MsgTypeUpdatePlayers
	MsgTypeUpdateSelf
	MsgTypeTimeAnswer
)
