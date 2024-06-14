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
	fmt.Println("Trigger")
	PartString = strings.SplitN(mess, " ", 3)
	fmt.Println(PartString[0])
	fmt.Println(PartString[1])
	fmt.Println(PartString[2])
	if len(PartString) > 1 && strings.HasPrefix(PartString[2], "BATTLE_START") {
		handleBattleMessage(PartString[2])
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
	clients ,err := battle.LoadClients("clients.json")
	if err != nil {
		fmt.Printf("Error reading clients data: %v\n", err)
		return
	}

	// Find the players based on their connection addresses
	var player1, player2 battle.Client
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

	// Initialize the players
	player1Data := battle.InitializePlayer(player1)
	player2Data := battle.InitializePlayer(player2)

	// Create a Clients struct for the battle
	battleClients := battle.Clients{
		User: []battle.Client{player1, player2},
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

	err = bluePrint.UpdateBattleState("battle.json", battleState)
	if err != nil {
		fmt.Printf("Error updating battle state: %v\n", err)
	}
}


