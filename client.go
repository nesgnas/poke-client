package main

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"image/color"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	mapWidth  = 1000
	mapHeight = 1000
	gridSize  = 30
	cellSize  = 16
)

var (
	playerX, playerY int
	playerRect       *canvas.Rectangle
	scrollContainer  *container.Scroll
)

func main() {
	// Read the map from file
	mapData, err := readMap("map.txt")
	if err != nil {
		fmt.Println("Error reading map file:", err)
		return
	}

	a := app.New()
	w := a.NewWindow("Poke Client")

	// Create a grid of rectangles based on the map data
	var rectangles []fyne.CanvasObject
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			var rect *canvas.Rectangle
			if i < mapHeight && j < mapWidth {
				if mapData[i][j] == '1' {
					rect = canvas.NewRectangle(color.Black) // Boundary tile (black)
				} else {
					rect = canvas.NewRectangle(color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Peace tile (green)
				}
			} else {
				rect = canvas.NewRectangle(color.Black) // If outside map bounds, show boundary (black)
			}
			rect.SetMinSize(fyne.NewSize(cellSize, cellSize))
			rectangles = append(rectangles, rect)
		}
	}

	// Add the player object
	playerX, playerY = 12, 8 // Initial position
	playerRect = canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	playerRect.SetMinSize(fyne.NewSize(cellSize, cellSize))

	// Create a grid layout
	grid := container.New(layout.NewGridLayout(gridSize), rectangles...)
	grid.Objects[playerY*gridSize+playerX] = playerRect

	// Create a scroll container for the grid
	scrollContainer = container.NewScroll(grid)
	scrollContainer.SetMinSize(fyne.NewSize(cellSize*25, cellSize*17)) // Viewport size

	// Set initial scroll position to center on the player
	scrollTo(playerX, playerY)

	// Set up key event handling
	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		movePlayer(key, mapData, grid)
	})

	// Add the scroll container to the window
	w.SetContent(container.NewVBox(
		widget.NewLabel("Use WASD to move the red object"),
		scrollContainer,
	))
	w.Resize(fyne.NewSize(cellSize*25, cellSize*17+50)) // +50 to accommodate the label
	w.ShowAndRun()
}

// readMap reads the map from a file and returns a 2D slice of runes
func readMap(filename string) ([][]rune, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var mapData [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		mapData = append(mapData, []rune(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return mapData, nil
}

// movePlayer moves the player object based on the key press
func movePlayer(key *fyne.KeyEvent, mapData [][]rune, grid *fyne.Container) {
	newX, newY := playerX, playerY

	switch key.Name {
	case fyne.KeyW:
		newY--
	case fyne.KeyS:
		newY++
	case fyne.KeyA:
		newX--
	case fyne.KeyD:
		newX++
	}

	if newX >= 0 && newX < gridSize && newY >= 0 && newY < gridSize &&
		newX < mapWidth && newY < mapHeight && mapData[newY][newX] == '0' {
		// Move the player if the new position is within bounds and not a boundary
		grid.Objects[playerY*gridSize+playerX] = canvas.NewRectangle(color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Restore the previous tile
		playerX, playerY = newX, newY
		grid.Objects[playerY*gridSize+playerX] = playerRect // Move the player to the new position
		grid.Refresh()

		// Update the scroll position to keep the player in view
		scrollTo(playerX, playerY)
	}
}

// scrollTo scrolls the container to keep the player in view
func scrollTo(x, y int) {
	scrollX := float32(x)*cellSize - scrollContainer.Size().Width/2 + cellSize/2
	scrollY := float32(y)*cellSize - scrollContainer.Size().Height/2 + cellSize/2

	// Ensure scroll offsets are within valid range
	maxScrollX := float32(mapWidth)*cellSize - scrollContainer.Size().Width
	maxScrollY := float32(mapHeight)*cellSize - scrollContainer.Size().Height

	if scrollX < 0 {
		scrollX = 0
	} else if scrollX > maxScrollX {
		scrollX = maxScrollX
	}
	if scrollY < 0 {
		scrollY = 0
	} else if scrollY > maxScrollY {
		scrollY = maxScrollY
	}

	scrollContainer.Offset.X = scrollX
	scrollContainer.Offset.Y = scrollY
	scrollContainer.Refresh()
}
