package intruction

import (
	"client/connectionWorld"
	"fmt"
	"net"
)

func NoticePosition(conn net.Conn) {

	param := "PUBLISH"
	channel := "hi"
	address := connectionWorld.LocalAddress
	playerX := connectionWorld.PlayerX
	playerY := connectionWorld.PlayerY

	message := fmt.Sprintf("%s %s %s %d %d", param, channel, address, playerX, playerY)

	connectionWorld.SendMessage(conn, message)
	connectionWorld.ListenForMessages(conn)
}
