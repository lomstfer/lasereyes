package game

import (
	"image/color"
	"time"
	"wzrds/client/internal"
	"wzrds/client/internal/network"
	"wzrds/common/commonutils"
	"wzrds/common/netmsg"
	"wzrds/common/netmsg/msgfromclient"
	"wzrds/common/netmsg/msgfromserver"
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Draw(screen *ebiten.Image) {
	g.HandleNetworkEvents()

	if !g.finishedAssetLoading || !g.timeSyncer.FinishedSync {
		screen.Fill(color.NRGBA{255, 0, 0, 255})
		return
	}

	if g.selfPlayer == nil {
		return
	}

	g.RealUpdate()

	text.Draw(screen, "hello", *g.fontFace, 0, 24, color.NRGBA{255, 255, 255, 255})
	DrawPlayer(screen, g.playerImage, g.selfPlayer.Data)
	for _, p := range g.otherPlayers {
		DrawPlayer(screen, g.playerImage, p.Data)
	}
}

func (g *Game) HandleNetworkEvents() {
	netMessage := g.netClient.CheckForEvents()
	switch msg := netMessage.(type) {
	case network.Connected:
		go g.syncTime(time.Millisecond * 100)

	case network.Disconnected:
		g.cleanClose = true

	case msgfromserver.AddSelfPlayer:
		g.selfPlayer = internal.NewSelfPlayer(msg.Data)

	case msgfromserver.AddOtherPlayer:
		g.otherPlayers[msg.Data.Id] = &internal.Player{Data: msg.Data}

	case msgfromserver.UpdatePlayers:
		for id, snapshot := range msg.IdsToSnapshots {
			if id == g.selfPlayer.Data.Id {
				continue
			}
			g.otherPlayers[id].SnapshotsForInterp = append(g.otherPlayers[id].SnapshotsForInterp, snapshot)
		}

	case msgfromserver.UpdateSelf:
		g.selfPlayer.HandleServerUpdate(msg.LastAuthorizedInputId, msg.Snapshot)

	case msgfromserver.TimeAnswer:
		g.timeSyncer.OnTimeAnswer(commonutils.GetCurrentTimeAsFloat(), msg.Request.TimeSent, msg.TimeReceived)
	}
}

func (g *Game) RealUpdate() {
	g.time = commonutils.GetCurrentTimeAsFloat() - g.startTime
	timeNow := commonutils.GetCurrentTimeAsFloat()
	// dt := timeNow - g.lastUpdateTime
	g.lastUpdateTime = timeNow

	g.getInputCallback.Update(func() {
		inputVec := vec2.Vec2{}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			inputVec.Y -= 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			inputVec.Y += 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			inputVec.X -= 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			inputVec.X += 1
		}
		input := player.Input{
			Up:    inputVec.Y == -1,
			Down:  inputVec.Y == 1,
			Left:  inputVec.X == -1,
			Right: inputVec.X == 1,
		}
		if input.HasInput() {
			g.selfPlayer.AddInput(input)
		}
	})

	g.sendInputCallback.Update(func() {
		packetStruct := msgfromclient.MoveInput{Input: g.selfPlayer.InputsToSend}
		bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromclient.MsgTypeMoveInput), packetStruct)
		g.netClient.SendToServer(bytes, true)
		g.selfPlayer.OnSendInputs()
	})

	for _, p := range g.otherPlayers {
		p.LerpBetweenSnapshots(g.timeSyncer.ServerTime())
	}
}

func DrawPlayer(screen *ebiten.Image, playerImage *ebiten.Image, pData player.CommonData) {
	geo := ebiten.GeoM{}
	geo.Translate(pData.Position.X, pData.Position.Y)
	screen.DrawImage(playerImage, &ebiten.DrawImageOptions{GeoM: geo})
}
