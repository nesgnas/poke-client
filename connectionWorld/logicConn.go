package connectionWorld

import (
	"fmt"
	"net"
	"strings"
)

var (
	PrePokerPosX int = 0
	PrePokerPosY int = 0
)

func initiateConnectionWorld(conn net.Conn) {
	subscribeMessage := "SUBSCRIBE hi"
	_, err := conn.Write([]byte(subscribeMessage + "\n"))
	if err != nil {
		fmt.Printf("Error sending subscription to server: %v\n", err)
		return
	}
	CHANNEL_HI = true
	fmt.Println("Subscribed to channel 'hi'")
}

func checkChannelAction(mess string) {
	parts := strings.SplitN(mess, " ", 2)

	fmt.Println(parts)
}
