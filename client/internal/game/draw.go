package game

import (
	"image/color"
	"time"
	"wzrds/common/netmsg"
	"wzrds/common/netmsg/msgfromclient"
	"wzrds/common/netmsg/msgfromserver"
	"wzrds/common/pkg/vec2"
	"wzrds/common/player"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.finishedAssetLoading {
		screen.Fill(color.NRGBA{255, 0, 0, 255})
		return
	}

	netMessage := g.netClient.CheckForEvents()
	switch msg := netMessage.(type) {
	case msgfromserver.DisconnectSelf:
		g.cleanClose = true
	case msgfromserver.AddSelfPlayer:
		g.selfPlayer = msg.Data
	case msgfromserver.AddOtherPlayer:
		g.otherPlayers[msg.Data.Id] = msg.Data
	case msgfromserver.UpdatePlayers:
		for id, pos := range msg.Players {
			if id == g.selfPlayer.Id {
				g.selfPlayer.Position = pos
			} else {
				player := g.otherPlayers[id]
				player.Position = pos
				g.otherPlayers[id] = player
			}
		}
	}

	g.time = time.Since(g.startTime).Seconds()
	timeNow := time.Now()
	// dt := timeNow.Sub(g.lastUpdateTime).Seconds()
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
		input := player.PlayerInput{
			Up:    inputVec.Y == -1,
			Down:  inputVec.Y == 1,
			Left:  inputVec.X == -1,
			Right: inputVec.X == 1,
		}
		g.accumulatedPlayerInputs = append(g.accumulatedPlayerInputs, input)
	})

	g.sendInputCallback.Update(func() {
		packetStruct := msgfromclient.MoveInput{Input: g.accumulatedPlayerInputs}
		g.accumulatedPlayerInputs = make([]player.PlayerInput, 0)
		bytes := netmsg.GetBytesFromIdAndStruct(byte(msgfromclient.MsgTypeMoveInput), packetStruct)
		g.netClient.SendToServer(bytes, true)
	})

	text.Draw(screen, "hello", *g.fontFace, 0, 24, color.NRGBA{255, 255, 255, 255})
	DrawPlayer(screen, g.playerImage, g.selfPlayer)
	for _, p := range g.otherPlayers {
		DrawPlayer(screen, g.playerImage, p)
	}
}

func DrawPlayer(screen *ebiten.Image, playerImage *ebiten.Image, pData player.PlayerSpawnData) {
	geo := ebiten.GeoM{}
	geo.Translate(pData.Position.X, pData.Position.Y)
	screen.DrawImage(playerImage, &ebiten.DrawImageOptions{GeoM: geo})
}
