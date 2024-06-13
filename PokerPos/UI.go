package PokerPos

import (
	"client/bluePrint"
	"client/connectionWorld"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/rand"
)

var (
	FilePos   = "storeFile/pokerPos.json"
	IndexFile int
)

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PokerPos struct {
	ID        string       `json:"id"`
	Positions []Coordinate `json:"position"`
}

type EnemyPos struct {
	ID        string       `json:"id"`
	Positions []Coordinate `json:"position"`
}

// Global variable to hold the poker player data
var PokerPlayer PokerPos

func RandomPos(filename string) int {
	rand.Seed(uint64(time.Now().UnixNano()))

	// Generate random coordinates
	IndexFile := getNextIndex(filename)
	PokerPlayer.Positions = append(PokerPlayer.Positions, Coordinate{
		X: rand.Intn(49) + 1,
		Y: rand.Intn(49) + 1,
	})
	PokerPlayer.ID = connectionWorld.LocalAddress

	fmt.Println(IndexFile)

	err := UpdatePokerPos(FilePos, IndexFile, PokerPlayer.Positions[0].X, PokerPlayer.Positions[0].Y, PokerPlayer.ID)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return IndexFile
}

// ReadPokerPos reads the initial position from a JSON file
func ReadPokerPos(filename string) ([]PokerPos, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var data struct {
		PokerPoses []PokerPos `json:"Poker"`
	}

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	file.Close()
	return data.PokerPoses, nil
}

// getNextIndex calculates the next index based on the current number of entries in the JSON file
func getNextIndex(filename string) int {
	pokerPoses, err := ReadPokerPos(filename)
	if err != nil {
		return 0
	}
	return len(pokerPoses)
}

// UpdatePokerPos updates the x and y coordinates in pokerPos.json
func UpdatePokerPos(filename string, index int, newX int, newY int, id string) error {

	pokerPoses, err := ReadPokerPos(filename)

	if err != nil {
		return err
	}

	// Ensure the pokerPoses slice has enough space for the new index

	if index >= len(pokerPoses) {
		defaultPosition := Coordinate{X: -1, Y: -1}
		pokerPoses = append(pokerPoses, PokerPos{
			ID:        id,
			Positions: []Coordinate{defaultPosition},
		})

	}

	PokerPlayer = pokerPoses[index]

	PokerPlayer.Positions[0].X = newX

	PokerPlayer.Positions[0].Y = newY
	PokerPlayer.ID = id

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(struct {
		Poker []PokerPos `json:"Poker"`
	}{Poker: pokerPoses})
	if err != nil {
		return err
	}

	fmt.Printf("Updated "+FilePos+" with x=%d, y=%d\n", newX, newY)
	return nil
}
func DeleteInvalidPokerPos(clientFile, pokerFile string) error {
	// Read and parse clients.json
	clientFileData, err := os.ReadFile(clientFile)
	if err != nil {
		return fmt.Errorf("failed to read client file: %w", err)
	}

	var userDatas bluePrint.UserData
	if err := json.Unmarshal(clientFileData, &userDatas); err != nil {
		return fmt.Errorf("failed to unmarshal client data: %w", err)
	}

	// Extract valid connection addresses
	validConnAdd := make(map[string]struct{})
	for _, client := range userDatas.User {
		validConnAdd[client.ConnAdd] = struct{}{}
	}

	// Read and parse pokerPos.json
	pokerFileData, err := os.ReadFile(pokerFile)
	if err != nil {
		return fmt.Errorf("failed to read poker file: %w", err)
	}

	var pokerData struct {
		Poker []PokerPos `json:"Poker"`
	}
	if err := json.Unmarshal(pokerFileData, &pokerData); err != nil {
		return fmt.Errorf("failed to unmarshal poker data: %w", err)
	}

	// Filter out invalid PokerPos entries and collect existing IDs
	var updatedPoker []PokerPos
	for _, poker := range pokerData.Poker {
		if _, valid := validConnAdd[poker.ID]; valid {
			updatedPoker = append(updatedPoker, poker)
			delete(validConnAdd, poker.ID) // Remove from validConnAdd as it is already present
		}
	}

	// Add missing valid IDs to pokerPos.json
	for connAdd := range validConnAdd {
		newPokerPos := PokerPos{ID: connAdd}
		updatedPoker = append(updatedPoker, newPokerPos)
	}

	// Write the updated pokerPos.json
	pokerData.Poker = updatedPoker
	updatedPokerFileData, err := json.MarshalIndent(pokerData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated poker data: %w", err)
	}

	if err := os.WriteFile(pokerFile, updatedPokerFileData, 0644); err != nil {
		return fmt.Errorf("failed to write updated poker file: %w", err)
	}

	fmt.Printf("Updated %s by removing invalid entries and adding missing valid ones\n", pokerFile)
	return nil
}

func UpdateEnemyFile(pokerFile, enemyFile string, localAddress string) error {
	// Read and parse pokerPos.json
	pokerFileData, err := os.ReadFile(pokerFile)
	if err != nil {
		return err
	}

	var pokerData struct {
		Poker []PokerPos `json:"Poker"`
	}
	err = json.Unmarshal(pokerFileData, &pokerData)
	if err != nil {
		return err
	}

	// Filter PokerPos entries to move to EnemyPos
	var enemyData struct {
		Enemies []EnemyPos `json:"enemy"`
	}
	for _, poker := range pokerData.Poker {
		enemyData.Enemies = append(enemyData.Enemies, EnemyPos{
			ID:        poker.ID,
			Positions: poker.Positions,
		})
	}

	// Write the updated enemy.json
	updatedEnemyFileData, err := json.MarshalIndent(enemyData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(enemyFile, updatedEnemyFileData, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Updated %s based on local address %s\n", enemyFile, localAddress)
	return nil
}
