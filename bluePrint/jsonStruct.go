package bluePrint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Data UserData

type User struct {
	ConnAdd     string      `json:"connAdd"`
	ListPokemon [][]Pokemon `json:"listPokemon"` // Change type to [][]Pokemon
	MaxValue    string      `json:"maxValue"`
	PositionX   int         `json:"positionX"`
	PositionY   int         `json:"positionY"`
	SpaceLeft   string      `json:"spaceLeft"`
	UID         string      `json:"uID"`
}

type Pokemon struct {
	UID string  `json:"uid"`
	ID  int     `json:"id"`
	Exp int     `json:"exp"`
	EV  float64 `json:"ev"`
	Lv  int     `json:"lv"`
}

type UserData struct {
	User []User `json:"user"`
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

	//for _, user := range Data.User {
	//	fmt.Println("User:", user.UID)
	//	fmt.Println("Connection Address:", user.ConnAdd)
	//	fmt.Println("Position:", user.PositionX, user.PositionY)
	//	fmt.Println("------------------------------")
	//}
}

func UpdateBattleState(fileName string, state interface{}) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, 0644)
}

func ReadBattleState(fileName string, state interface{}) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, state)
}
