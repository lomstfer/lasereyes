package network

import (
	"fmt"
	"os"
	"wzrds/common/commonutils"
	"wzrds/common/netmsg"
	"wzrds/common/netmsg/msgfromclient"

	"github.com/codecat/go-enet"
)

type NetworkServer struct {
	enetServerHost enet.Host
	enetPeers      map[uint]enet.Peer
}

func NewNetworkServer() *NetworkServer {
	enet.Initialize()
	ns := &NetworkServer{}

	var err error
	ns.enetServerHost, err = enet.NewHost(enet.NewListenAddress(5005), 32, 1, 0, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ns.enetPeers = make(map[uint]enet.Peer)

	return ns
}

func (ns *NetworkServer) CheckForEvents() (uint, interface{}) {
	event := ns.enetServerHost.Service(0)

	switch event.GetType() {
	case enet.EventConnect:
		peer := event.GetPeer()
		id := peer.GetConnectId()
		peer.SetData(commonutils.UintToByteArray(id))
		ns.enetPeers[id] = peer
		fmt.Println("connected id:", id)

	case enet.EventDisconnect:
		peer := event.GetPeer()
		id := commonutils.ByteArrayToUint(peer.GetData())
		peer.SetData(nil)
		fmt.Println("disconnected id:", id)
		return id, ClientDisconnected{}

	case enet.EventReceive:
		peerId := event.GetPeer().GetConnectId()
		packet := event.GetPacket()
		bytes := packet.GetData()
		packet.Destroy()
		id := bytes[0]
		bytes = bytes[1:]
		switch id {
		case byte(msgfromclient.MsgTypeConnectMe):
			s := netmsg.GetStructFromBytes[msgfromclient.ConnectMe](bytes)
			return peerId, s

		case byte(msgfromclient.MsgTypeInput):
			s := netmsg.GetStructFromBytes[msgfromclient.Input](bytes)
			return peerId, s

		case byte(msgfromclient.MsgTypeTimeRequest):
			s := netmsg.GetStructFromBytes[msgfromclient.TimeRequest](bytes)
			return peerId, s
		}
	}

	return 0, nil
}

func (ns *NetworkServer) Stop() {
	fmt.Println("network server stop")
	ns.enetServerHost.Destroy()
	enet.Deinitialize()
}

func (ns *NetworkServer) SendTo(id uint, msg []byte, reliable bool) {
	flag := enet.PacketFlagReliable
	if !reliable {
		flag = enet.PacketFlagUnsequenced
	}
	ns.enetPeers[id].SendBytes(msg, 0, flag)
}

func (ns *NetworkServer) SendToAll(msg []byte, reliable bool) {
	flag := enet.PacketFlagReliable
	if !reliable {
		flag = enet.PacketFlagUnsequenced
	}
	ns.enetServerHost.BroadcastBytes(msg, 0, flag)
}

func (ns *NetworkServer) SendToAllExcept(id uint, msg []byte, reliable bool) {
	flag := enet.PacketFlagReliable
	if !reliable {
		flag = enet.PacketFlagUnsequenced
	}

	for epid, ep := range ns.enetPeers {
		if epid != id {
			ep.SendBytes(msg, 0, flag)
		}
	}
}
