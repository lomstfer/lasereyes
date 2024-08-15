package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wzrds/common"
	"wzrds/common/commonutils"
	"wzrds/common/netmsg"
	"wzrds/common/netmsg/msgfromclient"
	"wzrds/common/netmsg/msgfromserver"
	"wzrds/common/player"
	"wzrds/server/internal"
	"wzrds/server/internal/network"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	netServer := network.NewNetworkServer()
	gameServer := internal.NewGameServer()

	simulationCallback := common.NewFixedCallback(1.0 / 60.0)
	broadcastGameCallback := common.NewFixedCallback(1.0 / 10.0)

	startedTime := commonutils.GetCurrentTimeAsFloat()

	go func() {
		for {
			serverTime := commonutils.GetCurrentTimeAsFloat() - startedTime

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

			case network.ClientDisconnected:
				gameServer.RemovePlayer(eventPeerId)

			case msgfromclient.MoveInput:
				gameServer.HandlePlayerInput(eventPeerId, msg.Input, serverTime)

			case msgfromclient.TimeRequest:
				s := msgfromserver.TimeAnswer{Request: msg, TimeReceived: serverTime}
				bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromserver.MsgTypeTimeAnswer), s)
				netServer.SendTo(eventPeerId, bytes, false)
			}

			simulationCallback.Update(func() {
				gameServer.Simulate(1.0/60.0, serverTime)
			})

			broadcastGameCallback.Update(func() {

				playersToUpdate := make(map[uint]player.Snapshot, 0)
				for id := range gameServer.Players {
					p := gameServer.Players[id]
					snapshot := player.Snapshot{Time: serverTime, Position: p.Data.Position}
					playersToUpdate[id] = snapshot

					fmt.Println(len(p.QueuedInputs), p.Data.Position.Sub(p.LastPos).Length())
					p.LastPos = p.Data.Position

					{
						s := msgfromserver.UpdateSelf{LastAuthorizedInputId: p.LastAuthorizedInputId, Snapshot: snapshot}
						bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromserver.MsgTypeUpdateSelf), s)
						netServer.SendTo(id, bytes, false)
					}
				}
				for k := range gameServer.PlayersThatMoved {
					delete(gameServer.PlayersThatMoved, k)
				}

				s := msgfromserver.UpdatePlayers{IdsToSnapshots: playersToUpdate}
				bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromserver.MsgTypeUpdatePlayers), s)
				netServer.SendToAll(bytes, false)
			})

			time.Sleep(time.Millisecond * 1)
		}
	}()

	<-sigChan
	netServer.Stop()
}
