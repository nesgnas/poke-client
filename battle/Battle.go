package battle

import (
	"client/bluePrint"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
)

type Player bluePrint.Player
type PokedexType bluePrint.PokedexType
type Evolution bluePrint.Evolution
type TypeEffectiveness bluePrint.TypeEffectiveness
type Pokemon bluePrint.Pokemon
type Client bluePrint.User
type Clients bluePrint.UserData

var pokedex = bluePrint.PokeDex
var typeEffectivenessMap = bluePrint.TypeEffectivenessMap

func (p *Player) SwitchPokemon(index int) {
	if index < len(p.Pokemons) && p.Pokemons[index].Alive {
		p.Active = index
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
			if s == item {
					return true
			}
	}
	return false
}

func (p *Player) Attack(target *Player, attackerID string, defenderID string) {
	attacker := p.Pokemons[p.Active]
	defender := target.Pokemons[target.Active]

	// Ensure the attacker is alive
	if !attacker.Alive {
		return
	}

	// Determine the most effective attack type
	attackType := 0 // Default to normal attack
	for _, atkType := range attacker.Type {
		if typeData, ok := typeEffectivenessMap[atkType]; ok {
			for _, defType := range defender.Type {
				if contains(typeData.Effective, defType) {
					attackType = 1 // Use special attack if effective
					break
				}
				if contains(typeData.NoEffect, defType) {
					attackType = 0 // Use normal attack if no effect
					break
				}
			}
		}
	}

	var damage int
	if attackType == 0 { // Normal attack
		damage = attacker.Base.Attack - defender.Base.Defense
	} else { // Special attack
		damage = attacker.Base.SpAttack*2 - defender.Base.SpDefense
	}

	if damage < 0 {
		damage = 0
	}

	defender.Base.HP -= damage
	defender.Base.HP -= damage
	fmt.Printf("%s (Player %s) attacks %s (Player %s) for %d damage.\n", attacker.Name.English, attackerID, defender.Name.English, defenderID, damage)

	if defender.Base.HP <= 0 {
		defender.Alive = false
		fmt.Printf("%s has fainted!\n", defender.Name.English)
		for i, poke := range target.Pokemons {
			if poke.Alive {
				target.SwitchPokemon(i)
				break
			}
		}
	}
}

func (p *Player) LevelUp(exp int, level int) (int, int) {
	expToLevelUp := 100
	for i := 1; i < level; i++ {
		expToLevelUp *= 2
	}
	// fmt.Printf("Exp to level up to level %d: %d\n", level+1, expToLevelUp)

	if exp >= expToLevelUp {
		return p.LevelUp(exp-expToLevelUp, level+1)
	} else {
		return level, exp
	}
}

func (p *Player) CheckEvolution(pokemon *Pokemon) {
	// Assuming pokedex is a global variable or passed to this function
	if evolution, exists := pokedex[pokemon.ID]; exists {
		for _, next := range evolution.Evolution.Next {
			nextID := next[0]
			condition := next[1]

			// Use regex to check if the condition is level-based or item-based
			levelRe := regexp.MustCompile(`\d+`)
			itemRe := regexp.MustCompile(`use (.+)`)

			if levelRe.MatchString(condition) {
				// Level-based evolution
				levelStr := levelRe.FindString(condition)
				requiredLevel, err := strconv.Atoi(levelStr)
				if err != nil {
					fmt.Printf("Error converting required level: %v\n", err)
					continue
				}
				if pokemon.Level >= requiredLevel {
					nextPokemonID, err := strconv.Atoi(nextID)
					if err != nil {
						fmt.Printf("Error converting next Pokemon ID: %v\n", err)
						continue
					}
					if nextPokemon, exists := pokedex[nextPokemonID]; exists {
						currentName := pokemon.Name.English
						// Perform evolution
						pokemon.ID = nextPokemon.ID
						pokemon.Name = nextPokemon.Name
						pokemon.Base = nextPokemon.Base
						fmt.Printf("%s has evolved into %s!\n", currentName,nextPokemon.Name.English)
					}
				}
			} else if matches := itemRe.FindStringSubmatch(condition); matches != nil {
				// Item-based evolution
				item := matches[1]
				fmt.Printf("%s needs to use %s to evolve.\n", pokemon.Name.English, item)
				// Here you can add logic to handle item-based evolution if needed
			} else {
				fmt.Printf("Unknown evolution condition: %s\n", condition)
			}
		}
	}
}

func InitializePlayer(client Client) Player {
	var connAdd string = client.ConnAdd
	var pokemons []bluePrint.Pokemon
	for _, p := range client.ListPokemon {
		blueprintPokemon := pokedex[p.ID]
		pokemon := bluePrint.Pokemon{
			ID:     blueprintPokemon.ID,
			Name:   blueprintPokemon.Name,
			Type:   blueprintPokemon.Type,
			Base:   blueprintPokemon.Base,
			Level:  p.Lv,
			Exp:    p.Exp,
			Ev:     p.EV,
			Alive:  true,
		}
		pokemons = append(pokemons, pokemon)
	}
	return Player{ConnAdd: connAdd, Pokemons: pokemons, Active: 0}
}

func Battle(clients Clients) {
	if len(clients.User) < 1 {
		fmt.Println("Not enough clients to start a battle.")
		return
	}

	// Initialize the single player
	player := InitializePlayer(Client(clients.User[0]))

	// DEBUG: Print the initial state of the players
	fmt.Printf("Player %s: %v\n", player.ConnAdd, player)

	player1Wins := 0
	player2Wins := 0

	// Best-of-three matches, each round using a different Pokémon from the first 3
	for round := 0; round < 3; round++ {
		fmt.Printf("\nRound %d:\n", round+1)
		winner := conductMatch(&player1, &player2, round)
		if winner == 1 {
			player1Wins++
		} else if winner == 2 {
			player2Wins++
		}

		fmt.Printf("Player 1 wins: %d, Player 2 wins: %d\n", player1Wins, player2Wins)

		// Check if either player has won two matches
		if player1Wins == 2 || player2Wins == 2 {
			break
		}
	}

	// Break line for better readability
	fmt.Println("--------------------")

	// Determine the overall winner and calculate experience
	var winner, loser *Player
	if player1Wins > player2Wins {
		fmt.Printf("Player 1 wins the battle!\n")
		winner = &player1
		loser = &player2
	} else {
		fmt.Printf("Player 2 wins the battle!\n")
		winner = &player2
		loser = &player1
	}

	// Calculate total experience from the loser's Pokémon
	fmt.Printf("\nGaining experience...\n")
	totalExp := 100
	for _, pokemon := range loser.Pokemons {
		totalExp += pokemon.Level * 100
		for i :=1; i < pokemon.Level; i++ {
			totalExp = int(float64(totalExp) * 1.5)
		}
	}
	fmt.Printf("Total experience gained: %d\n", totalExp)

	// Distribute experience to the winner's Pokémon
	expPerPokemon := totalExp / len(winner.Pokemons)

	for i := range winner.Pokemons {
		winner.Pokemons[i].Exp += expPerPokemon
		fmt.Printf("%s gained %d experience.\n", winner.Pokemons[i].Name.English, expPerPokemon)
		currentLevel := winner.Pokemons[i].Level
		newLevel, remainingExp := winner.LevelUp(winner.Pokemons[i].Exp, winner.Pokemons[i].Level)
		winner.Pokemons[i].Level = newLevel
		winner.Pokemons[i].Exp = remainingExp
		
		fmt.Printf("%s level up from %d to %d\n", winner.Pokemons[i].Name.English, currentLevel, newLevel)

		// Check for evolution
		fmt.Printf("\nChecking for evolution...\n")
		winner.CheckEvolution(&winner.Pokemons[i])

		// Break line for better readability
		fmt.Println("--------------------")
	}



	// Notify about evolutions and end the battle
	for _, pokemon := range winner.Pokemons {
		fmt.Printf("%s is now level %d.\n", pokemon.Name.English, pokemon.Level)
	}
}

func conductMatch(player1, player2 *Player, round int) int {
	// Set the active Pokémon for this round
	player1.Active = round
	player2.Active = round

	// Determine the order of attacks based on speed
	if player1.Pokemons[player1.Active].Base.Speed > player2.Pokemons[player2.Active].Base.Speed {
		// Player 1 attacks first
		player1.Attack(player2, "1", "2")
		if !player2.Pokemons[player2.Active].Alive {
			return 1 // Player 1 wins
		}
		player2.Attack(player1, "2", "1")
		if !player1.Pokemons[player1.Active].Alive {
			return 2 // Player 2 wins
		}
	} else {
		// Player 2 attacks first
		player2.Attack(player1, "2", "1")
		if !player1.Pokemons[player1.Active].Alive {
			return 2 // Player 2 wins
		}
		player1.Attack(player2, "1", "2")
		if !player2.Pokemons[player2.Active].Alive {
			return 1 // Player 1 wins
		}
	}

	// If no Pokémon has fainted, decide randomly (for simplicity)
	return rand.Intn(2) + 1
}

func main() {
	// Load Pokedex
	var err error
	pokedex, err = LoadPokedex("pokedex.json")
	if err != nil {
		fmt.Println("Error loading Pokedex:", err)
		return
	}

	// Load Type Effectiveness
	if err := LoadTypeEffectiveness("types.json"); err != nil {
		fmt.Println("Error loading type effectiveness:", err)
		return
	}

	// Load Clients
	clients, err := LoadClients("clients.json")
	if err != nil {
		fmt.Println("Error loading Clients:", err)
		return
	}

	Battle(clients)
}