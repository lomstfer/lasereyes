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
				newPlayerData := player.CommonData{Name: msg.Name, Id: eventPeerId}

				// add new player to new player
				{
					addSelfStruct := msgfromserver.AddSelfPlayer{Data: newPlayerData}
					addSelfBytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromserver.MsgTypeAddSelfPlayer), addSelfStruct)
					netServer.SendTo(eventPeerId, addSelfBytes, true)
				}

				// add old players to new player
				for _, p := range gameServer.Players {
					addOtherStruct := msgfromserver.AddOtherPlayer{Data: p.Data}
					addOtherBytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromserver.MsgTypeAddOtherPlayer), addOtherStruct)
					netServer.SendTo(eventPeerId, addOtherBytes, true)
				}

				// add new player to old players
				{
					addOtherStruct := msgfromserver.AddOtherPlayer{Data: newPlayerData}
					addOtherBytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromserver.MsgTypeAddOtherPlayer), addOtherStruct)
					netServer.SendToAllExcept(eventPeerId, addOtherBytes, true)
				}

				gameServer.AddPlayer(newPlayerData)

			case msgfromclient.MoveInput:
				gameServer.HandlePlayerInput(eventPeerId, msg.Input)
			}

			simulationCallback.Update(func() {
				gameServer.Simulate(1.0 / 60.0)
			})

			broadcastGameCallback.Update(func() {
				playersToUpdate := make(map[uint]player.Snapshot, 0)
				for id := range gameServer.PlayersThatMoved {
					p := gameServer.Players[id]
					snapshot := player.Snapshot{Time: time.Now(), Position: p.Data.Position}
					playersToUpdate[id] = snapshot

					{
						s := msgfromserver.UpdateSelf{LastAuthorizedInputId: p.LastAuthorizedInputId, Snapshot: snapshot}
						bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromserver.MsgTypeUpdateSelf), s)
						netServer.SendTo(id, bytes, false)
					}
				}
				for k := range gameServer.PlayersThatMoved {
					delete(gameServer.PlayersThatMoved, k)
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
