package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Color is the RGB color of a pixel
type Color struct {
	R, G, B byte
}

// Pos represents an object position in a 2D space
type Pos struct {
	X, Y float32
}

// Object is an interface to game objects like paddle, ball, score etc
type Object interface {
	Draw()
}

// State is an unsigned integer that represents the various game states
type State uint

const (
	// Unknown is the uninitialized game state
	Unknown State = iota
	// StateStart game state
	StateStart
	// StatePlay game state
	StatePlay
)

var state State

// SetState updates the current game state with the newState
func SetState(newState State) {
	state = newState
}

// GetState returns the current game state. It returns the Unknown state
// if the game state has not been initialized
func GetState() State {
	return state
}

// InitState is a helper to initialize the game state to the start state
func InitState() {
	state = StateStart
}

var winWidth int
var winHeight int

// Init initializes the game library with the given window width and height
func Init(windowWidth, windowHeight int) {
	winWidth = windowWidth
	winHeight = windowHeight
}

func setPixel(x, y int, c Color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
	}
}

// ClearPixels resets/clears the pixel buffer
func ClearPixels(pixels []byte) {
	for index := range pixels {
		pixels[index] = 0
	}
}

// GetCenter returns the center position of the window
func GetCenter() Pos {
	return Pos{float32(winWidth) / 2, float32(winHeight) / 2}
}

// Lerp is the linear interpolation between point a and b
func Lerp(a, b, percent float32) float32 {
	return a + percent*(b-a)
}

// SetupControllers setups game controllers/joysticks and returns
// the list of game controllers attached
func SetupControllers() []*sdl.GameController {
	var gameControllers []*sdl.GameController

	for i := 0; i < sdl.NumJoysticks(); i++ {
		gameController := sdl.GameControllerOpen(i)
		gameControllers = append(gameControllers, gameController)
	}

	return gameControllers
}
