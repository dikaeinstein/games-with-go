package main

import (
	"fmt"
	"time"

	"github.com/dikaeinstein/games-with-go/balloons2/balloon"
	"github.com/dikaeinstein/games-with-go/evolvingpictures/apt"
	"github.com/veandco/go-sdl2/sdl"
)

type rgba struct {
	r, g, b byte
}

const winWidth = 800
const winHeight = 600
const winDepth = 100

func main() {
	sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Exploding Balloons", sdl.WINDOWPOS_UNDEFINED,
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

	// var audioSpec sdl.AudioSpec
	// deviceID, err := sdl.OpenAudioDevice("", false, &audioSpec, nil, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// explosionBytes, _ := sdl.LoadWAV("explode.wav")
	// defer sdl.FreeWAV(explosionBytes)
	// audioState := &balloon.AudioState{deviceID, audioSpec, explosionBytes}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	// imgs := loadImages("images", "balloon_")
	// balloons := loadBalloons(renderer, imgs, 25)

	currentMouseState := balloon.GetMouseState()
	// previousMouseState := currentMouseState
	var elapsedTime float32
	running := true

	x := apt.OpX{}
	y := apt.OpY{}
	sine := &apt.OpSin{}
	plus := &apt.OpPlus{}

	sine.Child = x
	plus.LeftChild = sine
	plus.RightChild = y

	tex := aptToTexture(plus, renderer, 800, 600)

	for running {
		frameStart := time.Now()

		currentMouseState = balloon.GetMouseState()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.TouchFingerEvent:
				if e.Type == sdl.FINGERDOWN {
					touchX := int(e.X * float32(winWidth))
					touchY := int(e.Y * float32(winHeight))
					currentMouseState.X = touchX
					currentMouseState.Y = touchY
					currentMouseState.LeftButton = true
				}
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}

		// if !currentMouseState.leftButton && previousMouseState.leftButton {
		// 	fmt.Println("left click")
		// }

		renderer.Copy(tex, nil, nil)

		// balloon.UpdateBalloons(balloons, elapsedTime, currentMouseState,
		// 	previousMouseState, audioState, winWidth, winHeight, winDepth)
		// sort.Sort(balloon.Slice(balloons))
		// for _, b := range balloons {
		// 	b.Draw(renderer)
		// }

		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Milliseconds())
		// // fmt.Println("ms per frame", elapsedTime)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Milliseconds())
		}

		// previousMouseState = currentMouseState
	}
}

func aptToTexture(node apt.Node, renderer *sdl.Renderer, w, h int) *sdl.Texture {
	pixels := make([]byte, w*h*4)
	scale := float32(255 / 2)
	offset := float32(-1.0 * scale)

	i := 0
	for yi := 0; yi < h; yi++ {
		y := float32(yi)/float32(h)*2 - 1
		for xi := 0; xi < w; xi++ {
			x := float32(xi)/float32(w)*2 - 1
			c := node.Eval(x, y)
			pixels[i] = byte(c*scale - offset)
			pixels[i+1] = byte(c*scale - offset)
			pixels[i+2] = byte(c*scale - offset)
			i += 4
		}
	}

	return pixelsToTexture(renderer, pixels, w, h)
}

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
	if err != nil {
		panic(err)
	}

	tex.Update(nil, pixels, w*4)
	return tex
}
