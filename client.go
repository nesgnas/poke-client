package main

import (
	"bufio"
	"client/connectionWorld"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"golang.org/x/exp/rand"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	initialMapWidth  = 30
	initialMapHeight = 30

	maxMapH  = 100
	mapMapW  = 100
	cellSize = 16
)

var (
	playerX, playerY int
	playerRect       *canvas.Rectangle
	enemyRect        []Enemy
	scrollContainer  *container.Scroll
	coordsLabel      *widget.Label
	currentMap       [][]rune
	currentMapWidth  int
	currentMapHeight int
	w                fyne.Window
	grid             *fyne.Container
	rectangles       []fyne.CanvasObject
	OnBattle         bool
)

type Enemy struct {
	Rect     *canvas.Rectangle // Rectangle representing the enemy
	Position fyne.Position     // Position of the enemy (x, y)
}

func main() {
	OnBattle = false
	// Read the initial map from file
	mapData, err := readMap("map.txt")
	if err != nil {
		fmt.Println("Error reading map file:", err)
		return
	}

	// Read opponent positions from JSON file

	a := app.New()
	w = a.NewWindow("Poke Client")

	// Initialize map dimensions and data
	currentMap = mapData
	currentMapWidth = mapMapW
	currentMapHeight = maxMapH

	// Set initial player position
	playerX, playerY = 1, 1

	// Create the initial grid
	createGrid(currentMapWidth, currentMapHeight)

	// Create a label for displaying coordinates
	coordsLabel = widget.NewLabel(fmt.Sprintf("Coordinates: (%d, %d)", playerX, playerY))

	// Set up key event handling
	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyReturn {

		} else {
			movePlayer(key, a)
		}
	})

	// Add the scroll container and label to the window
	w.SetContent(container.NewVBox(
		widget.NewLabel("Use WASD to move the red object. Press Enter to expand the map."),
		scrollContainer,
		coordsLabel,
	))
	w.Resize(fyne.NewSize(cellSize*25, cellSize*17+50)) // +50 to accommodate the label
	w.Show()

	// Run the main Goroutine
	go connectionWorld.ConnecWorld()

	a.Run()
}

// createGrid initializes the grid layout based on the map dimensions
func createGrid(width, height int) {
	rectangles = make([]fyne.CanvasObject, 0, width*height)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			var rect *canvas.Rectangle
			if i < len(currentMap) && j < len(currentMap[0]) {
				if currentMap[i][j] == '1' {
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
	playerRect = canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	playerRect.SetMinSize(fyne.NewSize(cellSize, cellSize))

	// Create a grid layout
	grid = container.New(layout.NewGridLayout(width), rectangles...)
	grid.Objects[playerY*width+playerX] = playerRect

	// Add enemies to the grid
	//addEnemies()

	// Create a scroll container for the grid
	scrollContainer = container.NewScroll(grid)
	scrollContainer.SetMinSize(fyne.NewSize(cellSize*25, cellSize*17)) // Viewport size

	// Set initial scroll position to center on the player
	scrollTo(playerX, playerY)
}

// readMap reads the map from a file and returns a 2D slice of runes
func readMap(filename string) ([][]rune, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	count := 0

	var mapData [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > mapMapW {
			line = line[:mapMapW]
		}

		mapData = append(mapData, []rune(line))

		if len(mapData) >= maxMapH {
			break
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return mapData, nil
}

// movePlayer moves the player object based on the key press
// movePlayer moves the player object based on the key press
func movePlayer(key *fyne.KeyEvent, a fyne.App) {
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

	// Check collision with enemies
	for _, enemy := range enemyRect {
		if newX == int(enemy.Position.X) && newY == int(enemy.Position.Y) {
			// Collision detected
			// Handle collision logic here (e.g., player takes damage)
			handleCollision(w, a)
			return
		}
	}

	// Check if the new position is within bounds and not a boundary
	if newX >= 0 && newX < currentMapWidth && newY >= 0 && newY < currentMapHeight &&
		newX < len(currentMap[0]) && newY < len(currentMap) && currentMap[newY][newX] == '0' {
		// Move the player
		grid.Objects[playerY*currentMapWidth+playerX] = canvas.NewRectangle(color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Restore the previous tile
		playerX, playerY = newX, newY
		grid.Objects[playerY*currentMapWidth+playerX] = playerRect // Move the player to the new position

		addEnemiesFromOpponent("enemy.json")

		grid.Refresh()

		// Update the scroll position to keep the player in view
		scrollTo(playerX, playerY)

		// Update the coordinates label
		coordsLabel.SetText(fmt.Sprintf("Coordinates: (%d, %d)", playerX, playerY))
	}
}
func handleCollision(win fyne.Window, a fyne.App) {
	dialog.ShowConfirm("Collision Detected", "Do you want to open a BATTLE?", func(confirmed bool) {
		if confirmed {
			OnBattle = true
			// Open a new window or perform any action here if the user confirms
			fmt.Println("Opening a new window...")
			openNewWindow(win, a)
		} else {
			// Handle the case when the user clicks "No"
			fmt.Println("Collision ignored.")
		}
	}, w)
}

func openNewWindow(parent fyne.Window, a fyne.App) {
	w1 := a.NewWindow("BATTLE")
	closeButton := widget.NewButton("Close", func() {

	})
	w1.SetContent(container.NewVBox(
		widget.NewLabel("Battle with another"),
		closeButton,
	))

	w1.Show()
}

func scrollTo(x, y int) {
	scrollX := float32(x)*cellSize - scrollContainer.Size().Width/2 - cellSize  //449
	scrollY := float32(y)*cellSize - scrollContainer.Size().Height/2 - cellSize //272

	scrollContainer.Offset.X = scrollX
	scrollContainer.Offset.Y = scrollY + scrollY/4
	scrollContainer.Refresh()
}

func createEnemy(x, y int, color color.Color) Enemy {
	rect := canvas.NewRectangle(color)
	rect.SetMinSize(fyne.NewSize(cellSize, cellSize)) // Set the size of the rectangle
	return Enemy{
		Rect:     rect,
		Position: fyne.NewPos(float32(x), float32(y)),
	}
}

func addEnemies() {
	// Generate enemy positions randomly
	rand.Seed(uint64(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		x := rand.Intn(50)
		y := rand.Intn(50)
		enemy := createEnemy(x, y, color.RGBA{R: 255, G: 165, B: 0, A: 255}) // Orange color
		enemyRect = append(enemyRect, enemy)
		grid.Objects[y*currentMapWidth+x] = enemy.Rect
	}
	fmt.Println(enemyRect)
}

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Opponent struct {
	ID        int          `json:"id"`
	Positions []Coordinate `json:"position"`
}

func readOpponentPositions(filename string) ([]Opponent, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var opponentData struct {
		Opponents []Opponent `json:"enemy"`
	}

	err = json.NewDecoder(file).Decode(&opponentData)
	if err != nil {
		return nil, err
	}

	return opponentData.Opponents, nil
}

func addEnemiesFromOpponent(filename string) error {
	// Clear existing enemies
	enemyRect = nil

	// Read opponent positions from JSON file
	opponentRects, err := readOpponentPositions(filename)
	if err != nil {
		return err
	}

	// Add enemies to the grid based on opponent positions
	for _, opponent := range opponentRects {
		for _, pos := range opponent.Positions {
			// Create an enemy and add it to the grid
			enemy := createEnemy(pos.X, pos.Y, color.RGBA{R: 255, G: 165, B: 0, A: 255}) // Orange color
			enemyRect = append(enemyRect, enemy)
			grid.Objects[pos.Y*currentMapWidth+pos.X] = enemy.Rect
		}
	}

	return nil
}
