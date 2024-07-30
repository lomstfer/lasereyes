package network

import (
	"fmt"
	"os"
	"wzrds/common/utils"

	"github.com/codecat/go-enet"
)

type NetworkServer struct {
	enetServerHost enet.Host
}

func NewNetworkServer() *NetworkServer {
	enet.Initialize()
	ns := &NetworkServer{}

	var err error
	ns.enetServerHost, err = enet.NewHost(enet.NewListenAddress(8095), 32, 1, 0, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return ns
}

func (ns *NetworkServer) CheckForEvents() interface{} {
	event := ns.enetServerHost.Service(1000)

	switch event.GetType() {
	case enet.EventConnect:
		id := event.GetPeer().GetConnectId()
		fmt.Println("connected id:", id)
		event.GetPeer().SetData(utils.UintToByteArray(id))

	case enet.EventDisconnect:
		id := utils.ByteArrayToUint(event.GetPeer().GetData())
		fmt.Println("disconnected id:", id)

	case enet.EventReceive:
		fmt.Println("recieved from", event.GetPeer().GetConnectId())
		packet := event.GetPacket()
		fmt.Println("\tdata:", packet.GetData())
		packet.Destroy()
	}

	return nil
}

func (nc *NetworkServer) Stop() {
	nc.enetServerHost.Destroy()
	enet.Deinitialize()
}
