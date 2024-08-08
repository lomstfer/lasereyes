package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"wzrds/common"
	"wzrds/common/netmsg"
	"wzrds/common/netmsg/msgfromclient"
	"wzrds/common/netmsg/msgfromserver"
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"
	"wzrds/server/internal"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	netServer := internal.NewNetworkServer()
	gameServer := internal.NewGameServer()

	simulationCallback := common.NewFixedCallback(1.0 / 60.0)

	broadcastGameCallback := common.NewFixedCallback(1.0 / 10.0)

	go func() {
		for {
			eventPeerId, eventStruct := netServer.CheckForEvents()
			switch msg := eventStruct.(type) {
			case msgfromclient.ConnectMe:
				data := player.PlayerSpawnData{Name: msg.Name, Id: eventPeerId}
				gameServer.AddPlayer(data)
				{
					addSelfStruct := msgfromserver.AddSelfPlayer{Data: data}
					addSelfBytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromserver.MsgTypeAddSelfPlayer), addSelfStruct)
					netServer.SendTo(eventPeerId, addSelfBytes, true)
				}
				{
					addOtherStruct := msgfromserver.AddOtherPlayer{Data: data}
					addOtherBytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromserver.MsgTypeAddOtherPlayer), addOtherStruct)
					netServer.SendToAllExcept(eventPeerId, addOtherBytes, true)
				}
			case msgfromclient.MoveInput:
				gameServer.HandlePlayerInput(eventPeerId, msg.Input)
			}

			simulationCallback.Update(func() {
				gameServer.Simulate(1.0 / 60.0)
			})

			broadcastGameCallback.Update(func() {
				playersToUpdate := make(map[uint]vec2.Vec2, 0)
				for id, p := range gameServer.Players {
					playersToUpdate[id] = p.Position
				}
				s := msgfromserver.UpdatePlayers{Players: playersToUpdate}
				bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromserver.MsgTypeUpdatePlayers), s)
				netServer.SendToAll(bytes, false)
			})

			time.Sleep(time.Millisecond * 1)
		}
	}()

	<-sigChan
	netServer.Stop()
}
