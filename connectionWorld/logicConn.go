package connectionWorld

import (
	"client/battle"
	"client/bluePrint"
	"fmt"
	"net"
	"strings"
)

var (
	PrePokerPosX int = 0
	PrePokerPosY int = 0
	PartString []string
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
	PartString = strings.SplitN(mess, " ", 3)
	fmt.Println("checkChannelAction parts:", PartString)
	if len(PartString) > 1 && strings.HasPrefix(PartString[1], "BATTLE_START") {
		fmt.Println("Battle message detected")
		handleBattleMessage(PartString[1])
	}
}

func handleBattleMessage(message string) {
	// Parse the battle message
	PartString := strings.Split(message, " ")
	if len(PartString) < 3 {
		fmt.Println("Invalid battle message format")
		return
	}

	player1ConnAdd := PartString[1]
	player2ConnAdd := PartString[2]

	// Read the clients data from the store file
	var clients bluePrint.UserData
	err := bluePrint.ReadBattleState("storeFile/clients.json", &clients)
	if err != nil {
		fmt.Printf("Error reading clients data: %v\n", err)
		return
	}

	// Find the players based on their connection addresses
	var player1, player2 bluePrint.User
	for _, client := range clients.User {
		if client.ConnAdd == player1ConnAdd {
			player1 = client
		} else if client.ConnAdd == player2ConnAdd {
			player2 = client
		}
	}

	// Ensure both players are found
	if player1.ConnAdd == "" || player2.ConnAdd == "" {
		fmt.Println("One or both players not found")
		return
	}

	var player1Client = battle.Client(player1)
	var player2Client = battle.Client(player2)

	// Initialize the players
	player1Data := battle.InitializePlayer(player1Client)
	player2Data := battle.InitializePlayer(player2Client)

	// Create a Clients struct for the battle
	battleClients := battle.Clients{
		User: []bluePrint.User{player1, player2},
	}

	// Conduct the battle
	battle.Battle(battleClients)

	// Update the battle state in the store file
	battleState := map[string]interface{}{
		"winner": player1Data.ConnAdd,
		"loser":  player2Data.ConnAdd,
		"player1": map[string]interface{}{
			"connAdd": player1Data.ConnAdd,
			"pokemons": player1Data.Pokemons,
		},
		"player2": map[string]interface{}{
			"connAdd": player2Data.ConnAdd,
			"pokemons": player2Data.Pokemons,
		},
	}

	err = bluePrint.UpdateBattleState("storeFile/battle.json", battleState)
	if err != nil {
		fmt.Printf("Error updating battle state: %v\n", err)
	}
}


