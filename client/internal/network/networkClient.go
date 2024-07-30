package network

import (
	"fmt"
	"os"

	"wzrds/common"

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
		return common.ServerDisconnectedClient{}
	case enet.EventReceive:
		fmt.Println("recieved")
		packet := event.GetPacket()
		fmt.Println("\tdata:", packet.GetData())
		packet.Destroy()
		return common.ClientConnectionInfo{Name: "Peter Griffin"}
	}

	return nil
}

func (nc *NetworkClient) StartDisconnect() {
	nc.enetServerPeer.Disconnect(0)
}

func (nc *NetworkClient) CleanUp() {
	nc.enetClientHost.Destroy()
	enet.Deinitialize()
}
