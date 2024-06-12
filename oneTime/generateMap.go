package oneTime

import (
	"fmt"
	"os"
)

const (
// gridSize = 50
)

func mains() {
	file, err := os.Create("map.txt")
	if err != nil {
		fmt.Println("Error creating map file:", err)
		return
	}
	defer file.Close()

	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			if i == 0 || i == gridSize-1 || j == 0 || j == gridSize-1 {
				_, err = file.WriteString("1")
			} else {
				_, err = file.WriteString("0")
			}
			if err != nil {
				fmt.Println("Error writing to map file:", err)
				return
			}
		}
		_, err = file.WriteString("\n")
		if err != nil {
			fmt.Println("Error writing to map file:", err)
			return
		}
	}

	fmt.Println("Map file 'map.txt' generated successfully.")
}
