package main

import (
	"fmt"

	"github.com/dikaeinstein/games-with-go/noise"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth = 800
const winHeight = 600

type color struct {
	r, g, b byte
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED,
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

	octaves := 3
	var frequency float32 = 0.01
	var gain float32 = 0.2
	var lacunarity float32 = 3.0
	n := noise.MakeNoise(noise.TURBULENCE, winWidth, winHeight, frequency, lacunarity, gain, octaves)
	scale := noise.CalcScale(255.0)
	n.Rescale(scale)
	gradient := getDualGradient(
		color{0, 0, 175}, color{80, 160, 244},
		color{12, 192, 75}, color{255, 255, 255},
	)
	drawNoise(n, gradient, pixels)

	keyboardState := sdl.GetKeyboardState()
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}

		var mult float32 = 1
		if keyboardState[sdl.SCANCODE_LSHIFT] != 0 || keyboardState[sdl.SCANCODE_RSHIFT] != 0 {
			mult = -1
		}

		if keyboardState[sdl.SCANCODE_O] != 0 {
			octaves = octaves + 1*int(mult)
			n := noise.MakeNoise(noise.TURBULENCE, winWidth, winHeight, frequency, lacunarity, gain, octaves)
			scale := noise.CalcScale(255.0)
			n.Rescale(scale)
			drawNoise(n, gradient, pixels)
		}

		if keyboardState[sdl.SCANCODE_F] != 0 {
			frequency = frequency + float32(0.001)*mult
			n := noise.MakeNoise(noise.TURBULENCE, winWidth, winHeight, frequency, lacunarity, gain, octaves)
			scale := noise.CalcScale(255.0)
			n.Rescale(scale)
			drawNoise(n, gradient, pixels)
		}

		if keyboardState[sdl.SCANCODE_G] != 0 {
			gain = gain + float32(0.1)*mult
			n := noise.MakeNoise(noise.TURBULENCE, winWidth, winHeight, frequency, lacunarity, gain, octaves)
			scale := noise.CalcScale(255.0)
			n.Rescale(scale)
			drawNoise(n, gradient, pixels)
		}

		if keyboardState[sdl.SCANCODE_L] != 0 {
			lacunarity = lacunarity + float32(0.1)*mult
			n := noise.MakeNoise(noise.TURBULENCE, winWidth, winHeight, frequency, lacunarity, gain, octaves)
			scale := noise.CalcScale(255.0)
			n.Rescale(scale)
			drawNoise(n, gradient, pixels)
		}

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}
}

func lerp(b1, b2 byte, pct float32) byte {
	return uint8(float32(b1) + pct*(float32(b2)-float32(b1)))
}

func colorLerp(c1, c2 color, pct float32) color {
	return color{
		lerp(c1.r, c2.r, pct),
		lerp(c1.g, c2.g, pct),
		lerp(c1.b, c2.b, pct),
	}
}

func getGradient(c1, c2 color) []color {
	result := make([]color, 256)

	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}

	return result
}

func getDualGradient(c1, c2, c3, c4 color) []color {
	result := make([]color, 256)

	for i := range result {
		pct := float32(i) / float32(255)
		if pct < 0.5 {
			result[i] = colorLerp(c1, c2, pct*float32(2))
		} else {
			result[i] = colorLerp(c3, c4, pct*float32(1.5)-float32(0.5))
		}
	}

	return result
}

// clamp ensure v is within this interval or boundary
func clamp(min, max, v int) int {
	if v < min {
		min = v
	} else if v > max {
		max = v
	}

	return v
}

// RescaleAndDraw scales the noise values by the given scale
// and draws them to the pixels buffer
func drawNoise(noise []float32, gradient []color, pixels []byte) {
	for i := range noise {
		c := gradient[clamp(0, 255, int(noise[i]))]
		p := i * 4
		pixels[p] = c.r
		pixels[p+1] = c.g
		pixels[p+2] = c.b
	}
}
