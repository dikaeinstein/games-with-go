package main

import (
	"fmt"

	"github.com/dikaeinstein/games-with-go/simplexnoise/noise"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth = 800
const winHeight = 600

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

	// for y := 0; y < winHeight; y++ {
	// 	for x := 0; x < winWidth; x++ {
	// 		setPixel(x, y, color{255, 0, 0}, pixels)
	// 	}
	// }

	octaves := 3
	var frequency float32 = 0.01
	var gain float32 = 0.2
	var lacunarity float32 = 3.0
	n := noise.MakeNoise(winWidth, winHeight, frequency, lacunarity, gain, octaves)
	scale := noise.CalcScale(255.0)
	gradient := getDualGradient(
		noise.Color{R: 0, G: 0, B: 175}, noise.Color{R: 80, G: 160, B: 244},
		noise.Color{R: 12, G: 192, B: 75}, noise.Color{R: 255, G: 255, B: 255},
	)
	n.RescaleAndDraw(scale, gradient, pixels)

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
			n := noise.MakeNoise(winWidth, winHeight, frequency, lacunarity, gain, octaves)
			scale := noise.CalcScale(255.0)
			n.RescaleAndDraw(scale, gradient, pixels)
		}

		if keyboardState[sdl.SCANCODE_F] != 0 {
			frequency = frequency + float32(0.001)*mult
			n := noise.MakeNoise(winWidth, winHeight, frequency, lacunarity, gain, octaves)
			scale := noise.CalcScale(255.0)
			n.RescaleAndDraw(scale, gradient, pixels)
		}

		if keyboardState[sdl.SCANCODE_G] != 0 {
			gain = gain + float32(0.1)*mult
			n := noise.MakeNoise(winWidth, winHeight, frequency, lacunarity, gain, octaves)
			scale := noise.CalcScale(255.0)
			n.RescaleAndDraw(scale, gradient, pixels)
		}

		if keyboardState[sdl.SCANCODE_L] != 0 {
			lacunarity = lacunarity + float32(0.1)*mult
			n := noise.MakeNoise(winWidth, winHeight, frequency, lacunarity, gain, octaves)
			scale := noise.CalcScale(255.0)
			n.RescaleAndDraw(scale, gradient, pixels)
		}

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}
}

func setPixel(x, y int, c noise.Color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
	}
}

func lerp(b1, b2 byte, pct float32) byte {
	return uint8(float32(b1) + pct*(float32(b2)-float32(b1)))
}

func colorLerp(c1, c2 noise.Color, pct float32) noise.Color {
	return noise.Color{
		R: lerp(c1.R, c2.R, pct),
		G: lerp(c1.G, c2.G, pct),
		B: lerp(c1.B, c2.B, pct),
	}
}

func getGradient(c1, c2 noise.Color) []noise.Color {
	result := make([]noise.Color, 256)

	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}

	return result
}

func getDualGradient(c1, c2, c3, c4 noise.Color) []noise.Color {
	result := make([]noise.Color, 256)

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
