package UI

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/rand"
	"time"

	"image/color"
	"os"
)

const (
	MaxMapH  = 100
	MapMapW  = 100
	CellSize = 16
)

var (
	PlayerX, PlayerY int
	PlayerRect       *canvas.Rectangle
	EnemyRect        []Enemy
	ScrollContainer  *container.Scroll
	CoordsLabel      *widget.Label
	CurrentMap       [][]rune
	CurrentMapWidth  int
	CurrentMapHeight int
	W                fyne.Window
	Grid             *fyne.Container
	Rectangles       []fyne.CanvasObject
)

type Enemy struct {
	Rect     *canvas.Rectangle // Rectangle representing the enemy
	Position fyne.Position     // Position of the enemy (x, y)
}

// createGrid initializes the grid layout based on the map dimensions
func CreateGrid(width, height int) {
	Rectangles = make([]fyne.CanvasObject, 0, width*height)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			var rect *canvas.Rectangle
			if i < len(CurrentMap) && j < len(CurrentMap[0]) {
				if CurrentMap[i][j] == '1' {
					rect = canvas.NewRectangle(color.Black) // Boundary tile (black)
				} else {
					rect = canvas.NewRectangle(color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Peace tile (green)
				}
			} else {
				rect = canvas.NewRectangle(color.Black) // If outside map bounds, show boundary (black)
			}
			rect.SetMinSize(fyne.NewSize(CellSize, CellSize))
			Rectangles = append(Rectangles, rect)
		}
	}

	// Add the player object
	PlayerRect = canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	PlayerRect.SetMinSize(fyne.NewSize(CellSize, CellSize))

	// Create a grid layout
	Grid = container.New(layout.NewGridLayout(width), Rectangles...)
	Grid.Objects[PlayerY*width+PlayerY] = PlayerRect

	// Add enemies to the grid
	AddEnemies()

	// Create a scroll container for the grid
	ScrollContainer = container.NewScroll(Grid)
	ScrollContainer.SetMinSize(fyne.NewSize(CellSize*25, CellSize*17)) // Viewport size

	// Set initial scroll position to center on the player
	ScrollTo(PlayerY, PlayerX)
}

// readMap reads the map from a file and returns a 2D slice of runes
func ReadMap(filename string) ([][]rune, error) {
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
		if len(line) > MapMapW {
			line = line[:MapMapW]
		}
		mapData = append(mapData, []rune(line))

		if len(mapData) >= MaxMapH {
			break
		}
		count++
	}
	fmt.Println(count)
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return mapData, nil
}

// movePlayer moves the player object based on the key press
// movePlayer moves the player object based on the key press
func MovePlayer(key *fyne.KeyEvent, a fyne.App) {
	newX, newY := PlayerX, PlayerY

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
	for _, enemy := range EnemyRect {
		if newX == int(enemy.Position.X) && newY == int(enemy.Position.Y) {
			// Collision detected
			// Handle collision logic here (e.g., player takes damage)
			HandleCollision(W)
			return
		}
	}

	// Check if the new position is within bounds and not a boundary
	if newX >= 0 && newX < CurrentMapWidth && newY >= 0 && newY < CurrentMapHeight &&
		newX < len(CurrentMap[0]) && newY < len(CurrentMap) && CurrentMap[newY][newX] == '0' {
		// Move the player
		Grid.Objects[PlayerX*CurrentMapWidth+PlayerX] = canvas.NewRectangle(color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Restore the previous tile
		PlayerX, PlayerY = newX, newY
		Grid.Objects[PlayerY*CurrentMapWidth+PlayerX] = PlayerRect // Move the player to the new position
		Grid.Refresh()

		// Update the scroll position to keep the player in view
		ScrollTo(PlayerY, PlayerX)

		// Update the coordinates label
		CoordsLabel.SetText(fmt.Sprintf("Coordinates: (%d, %d)", PlayerX, PlayerY))
	}
}
func HandleCollision(win fyne.Window) {
	dialog.ShowConfirm("Collision Detected", "Do you want to open a new window?", func(confirmed bool) {
		if confirmed {
			// Open a new window or perform any action here if the user confirms
			fmt.Println("Opening a new window...")
			OpenNewWindow(win)
		} else {
			// Handle the case when the user clicks "No"
			fmt.Println("Collision ignored.")
		}
	}, W)
}

func OpenNewWindow(parent fyne.Window) {
	app := app.New()
	newWindow := app.NewWindow("New Window")
	newWindow.Resize(fyne.NewSize(300, 200))
	newWindow.CenterOnScreen()

	closeButton := widget.NewButton("Close", func() {
		newWindow.Close()
	})

	newWindow.SetContent(container.NewVBox(
		widget.NewLabel("This is a new window!"),
		closeButton,
	))
	newWindow.Show()
}

func ScrollTo(x, y int) {
	scrollX := float32(x)*CellSize - ScrollContainer.Size().Width/2 - CellSize  //449
	scrollY := float32(y)*CellSize - ScrollContainer.Size().Height/2 - CellSize //272

	ScrollContainer.Offset.X = scrollX
	ScrollContainer.Offset.Y = scrollY + scrollY/4
	ScrollContainer.Refresh()
}

func CreateEnemy(x, y int, color color.Color) Enemy {
	rect := canvas.NewRectangle(color)
	rect.SetMinSize(fyne.NewSize(CellSize, CellSize)) // Set the size of the rectangle
	return Enemy{
		Rect:     rect,
		Position: fyne.NewPos(float32(x), float32(y)),
	}
}

func AddEnemies() {
	// Generate enemy positions randomly
	rand.Seed(uint64(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		x := rand.Intn(50) + 1
		y := rand.Intn(50) + 1
		enemy := CreateEnemy(x, y, color.RGBA{R: 255, G: 165, B: 0, A: 255}) // Orange color
		EnemyRect = append(EnemyRect, enemy)
		Grid.Objects[y*CurrentMapWidth+x] = enemy.Rect
	}
	fmt.Println(EnemyRect)
}
