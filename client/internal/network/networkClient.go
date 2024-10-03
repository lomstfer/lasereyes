package network

import (
	"fmt"
	"os"
<<<<<<< HEAD
	"strconv"
	"strings"
=======
>>>>>>> b53fc98f772f7b9cccc62013e8daad43dd9f9f74
	"wzrds/common/netmsg/msgfromclient"

	"wzrds/common/netmsg"
	"wzrds/common/netmsg/msgfromserver"

	"github.com/codecat/go-enet"
)

type NetworkClient struct {
	enetClientHost enet.Host
	enetServerPeer enet.Peer
	Connected      bool
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

<<<<<<< HEAD
	ipFile, err := os.ReadFile("custom_server_info.txt")

	defaultIpStr := "127.0.0.1"
	defaultPort := uint16(5005)
	ipStr := defaultIpStr
	port := defaultPort

	if err != nil {
		fmt.Println("Did not find custom_server_info.txt, connects to the default ip: 127.0.0.1:5005")
	} else if len(ipFile) != 0 {
		fileInput := strings.Split(string(ipFile), ":")
		if len(fileInput) != 2 {
			fmt.Println("Error loading custom_server_info.txt data, ip and port should be in the format of ip:port. Connects to the default ip: 127.0.0.1:5005")
		} else {
			ipStr = fileInput[0]
			port64, err := strconv.ParseUint(fileInput[1], 10, 16)
			if err != nil {
				fmt.Println("Error loading custom_server_info.txt data, ip and port should be in the format of ip:port. Connects to the default ip: 127.0.0.1:5005")
				os.Exit(1)
			}
			port = uint16(port64)
		}
	}

	nc.enetServerPeer, err = nc.enetClientHost.Connect(enet.NewAddress(ipStr, port), 1, 0)
=======
	nc.enetServerPeer, err = nc.enetClientHost.Connect(enet.NewAddress("127.0.0.1", 5005), 1, 0)
>>>>>>> b53fc98f772f7b9cccc62013e8daad43dd9f9f74
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
<<<<<<< HEAD
		nc.Connected = true
=======
		fmt.Println("connected")
>>>>>>> b53fc98f772f7b9cccc62013e8daad43dd9f9f74
		return Connected{}

	case enet.EventDisconnect:
		nc.Connected = false
		return Disconnected{}

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
		case byte(msgfromserver.MsgTypeAddSelfPlayer):
			s := netmsg.GetStructFromBytes[msgfromserver.AddSelfPlayer](bytes)
			return s
		case byte(msgfromserver.MsgTypeAddOtherPlayer):
			s := netmsg.GetStructFromBytes[msgfromserver.AddOtherPlayer](bytes)
			return s
		case byte(msgfromserver.MsgTypeRemoveOtherPlayer):
			s := netmsg.GetStructFromBytes[msgfromserver.RemoveOtherPlayer](bytes)
			return s
		case byte(msgfromserver.MsgTypeUpdatePlayers):
			s := netmsg.GetStructFromBytes[msgfromserver.UpdatePlayers](bytes)
			return s
		case byte(msgfromserver.MsgTypeUpdateSelf):
			s := netmsg.GetStructFromBytes[msgfromserver.UpdateSelf](bytes)
			return s
		case byte(msgfromserver.MsgTypeTimeAnswer):
			s := netmsg.GetStructFromBytes[msgfromserver.TimeAnswer](bytes)
			return s
		case byte(msgfromserver.MsgTypePlayerTakeDamage):
			s := netmsg.GetStructFromBytes[msgfromserver.PlayerTakeDamage](bytes)
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

func (nc *NetworkClient) SendConnectMe(selfName string) {
	s := msgfromclient.ConnectMe{Name: selfName}
	bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromclient.MsgTypeConnectMe), s)
	nc.SendToServer(bytes, true)
}
