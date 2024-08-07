package network

import (
	"fmt"
	"os"

	"wzrds/common/netmsg"
	"wzrds/common/netmsg/msgfromserver"

	"github.com/codecat/go-enet"
)

type NetworkClient struct {
	enetClientHost enet.Host
	enetServerPeer enet.Peer
}

func NewNetworkClient() *NetworkClient {
	enet.Initialize()

	nc := &NetworkClient{}

	var err error
	nc.enetClientHost, err = enet.NewHost(nil, 1, 1, 0, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nc.enetServerPeer, err = nc.enetClientHost.Connect(enet.NewAddress("127.0.0.1", 8095), 1, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nc
}

func (nc *NetworkClient) CheckForEvents() interface{} {
	event := nc.enetClientHost.Service(0)

	switch event.GetType() {
	case enet.EventConnect:
		fmt.Println("connected")
		nc.enetServerPeer.SendString("hello", 0, enet.PacketFlagReliable)
	case enet.EventDisconnect:
		fmt.Println("disconnected", event.GetPeer())
		return msgfromserver.DisconnectSelf{}
	case enet.EventReceive:
		packet := event.GetPacket()
		bytes := packet.GetData()
		packet.Destroy()
		id := bytes[0]
		bytes = bytes[1:]
		switch id {
		case byte(msgfromserver.MsgTypeDisconnectedClient):
			s := netmsg.GetStructFromBytes[msgfromserver.DisconnectedClient](bytes)
			return s
		}
		return nil
	}

	return nil
}

func (nc *NetworkClient) SendToServer(msg []byte, reliable bool) {
	flag := enet.PacketFlagReliable
	if !reliable {
		flag = enet.PacketFlagUnsequenced
	}
	nc.enetServerPeer.SendBytes(msg, 0, flag)
}

func (nc *NetworkClient) StartDisconnect() {
	nc.enetServerPeer.Disconnect(0)
}

func (nc *NetworkClient) CleanUp() {
	nc.enetClientHost.Destroy()
	enet.Deinitialize()
}
