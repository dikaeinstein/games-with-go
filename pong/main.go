package main

import (
	"fmt"
	"time"

	"github.com/dikaeinstein/games-with-go/pong/game"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth = 800
const winHeight = 600

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Pong Game", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println("Could not create window:", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Could not create renderer:", err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING, winWidth, winHeight)
	if err != nil {
		fmt.Println("Could not create texture:", err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)
	game.Init(winWidth, winHeight)
	player1 := game.NewPaddle(game.Pos{X: 100, Y: 100}, 10, 100, 400,
		game.Score(0), game.Color{R: 255, G: 255, B: 255})
	aiPlayer := game.NewPaddle(game.Pos{X: winWidth - 100, Y: 100}, 10, 100, 400,
		game.Score(0), game.Color{R: 255, G: 255, B: 255})
	ball := game.NewBall(game.GetCenter(), 10, 300, 300, game.Color{R: 255, G: 255, B: 255})

	keyboardState := sdl.GetKeyboardState()
	var elapsedTime float32
	game.InitState()
	controllers := game.SetupControllers()
	var controllerAxis1 int16

	running := true
	for running {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}

		for _, controller := range controllers {
			if controller != nil {
				controllerAxis1 = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
			}
		}

		if game.GetState() == game.StatePlay {
			player1.Update(keyboardState, controllerAxis1, elapsedTime)
			ball.Update(player1, aiPlayer, elapsedTime)
			aiPlayer.AIUpdate(ball, elapsedTime)
		} else if game.GetState() == game.StateStart {
			if keyboardState[sdl.SCANCODE_SPACE] != 0 {
				if player1.GetScore() == 3 || aiPlayer.GetScore() == 3 {
					player1.ResetScore()
					aiPlayer.ResetScore()
				}
				game.SetState(game.StatePlay)
			}
		}

		game.ClearPixels(pixels)
		player1.Draw(pixels)
		aiPlayer.Draw(pixels)
		ball.Draw(pixels)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < .005 {
			sdl.Delay(5 - uint32(elapsedTime*1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}
}
