package bluePrint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var Data UserData
var (
	PokeDex PokedexType
	TypeEffectivenessMap map[string]TypeEffectiveness
)

type User struct {
	ConnAdd     string   `json:"connAdd"`
	// ListPokemon []string `json:"listPokemon"`
	ListPokemon []struct {
		UID string `json:"uid"`
		ID  int    `json:"id"`
		Exp int    `json:"exp"`
		EV  int    `json:"ev"`
		Lv  int    `json:"lv"`
	} `json:"listPokemon"`
	MaxValue    string   `json:"maxValue"`
	PositionX   int      `json:"positionX"`
	PositionY   int      `json:"positionY"`
	SpaceLeft   string   `json:"spaceLeft"`
	UID         string   `json:"uID"`
}

type UserData struct {
	User []User `json:"user"`
}

type Pokemon struct {
	ID       int
	Name     struct {
		English string
	}
	Type []string
	Base struct {
		HP       int
		Attack   int
		Defense  int
		SpAttack int `json:"Sp. Attack"`
		SpDefense int `json:"Sp. Defense"`
		Speed		int
	}
	Exp  int
	Ev	int
	Level int
	Alive bool
	Evolution Evolution `json:"evolution"`
}

type Evolution struct {
	Next [][]string `json:"next"`
	Prev []string `json:"prev"`
}

type TypeEffectiveness struct {
	English     string   `json:"english"`
	Effective   []string `json:"effective"`
	InEffective []string `json:"ineffective"`
	NoEffect    []string `json:"no_effect"`
}

type Player struct {
	ConnAdd string
	Pokemons []Pokemon
	Active   int
}

type PokedexType map[int]Pokemon

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
	//for _, user := range Data.Users {
	//	fmt.Println("User:", user.UID)
	//	fmt.Println("Connection Address:", user.ConnAdd)
	//	fmt.Println("Position:", user.PositionX, user.PositionY)
	//	fmt.Println("------------------------------")
	//}
}

func LoadTypeEffectiveness(filename string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
			return err
	}

	var types []TypeEffectiveness
	if err := json.Unmarshal(bytes, &types); err != nil {
			return err
	}

	TypeEffectivenessMap = make(map[string]TypeEffectiveness)
	for _, t := range types {
		TypeEffectivenessMap[t.English] = t
	}
	return nil
}

func LoadPokedex(filename string) (PokedexType, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var pokemons []Pokemon // Assuming Pokemon is the struct that matches the JSON structure
	if err := json.Unmarshal(bytes, &pokemons); err != nil {
			return nil, err
	}

	pokedex := make(PokedexType) // Assuming Pokedex is a map[int]Pokemon
	for _, pokemon := range pokemons {
			pokedex[pokemon.ID] = pokemon
	}

	return pokedex, nil
}