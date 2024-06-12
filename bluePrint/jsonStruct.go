package bluePrint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Data UserData

type User struct {
	ConnAdd     string   `json:"connAdd"`
	ListPokemon []string `json:"listPokemon"`
	MaxValue    string   `json:"maxValue"`
	PositionX   int      `json:"positionX"`
	PositionY   int      `json:"positionY"`
	SpaceLeft   string   `json:"spaceLeft"`
	UID         string   `json:"uID"`
}

type UserData struct {
	Users []User `json:"user"`
}

func ReadUser(fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Parse JSON data

	if err := json.Unmarshal(file, &Data); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Extract connection addresses
	for _, user := range Data.Users {
		fmt.Println("User:", user.UID)
		fmt.Println("Connection Address:", user.ConnAdd)
		fmt.Println("Position:", user.PositionX, user.PositionY)
		fmt.Println("------------------------------")
	}
}
