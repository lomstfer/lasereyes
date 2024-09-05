package msgfromserver

type MsgType byte

const (
	MsgTypeUndefined MsgType = iota
	MsgTypeDisconnectedClient
	MsgTypeAddSelfPlayer
	MsgTypeAddOtherPlayer
	MsgTypeRemoveOtherPlayer
	MsgTypeUpdatePlayers
	MsgTypeUpdateSelf
	MsgTypeTimeAnswer
	MsgTypePlayerTakeDamage
)
