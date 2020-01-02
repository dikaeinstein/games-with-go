package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dikaeinstein/games-with-go/balloons2/balloon"
	"github.com/dikaeinstein/games-with-go/noise"
	"github.com/dikaeinstein/games-with-go/vector"
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

	var audioSpec sdl.AudioSpec
	deviceID, err := sdl.OpenAudioDevice("", false, &audioSpec, nil, 0)
	if err != nil {
		panic(err)
	}
	explosionBytes, _ := sdl.LoadWAV("explode.wav")
	defer sdl.FreeWAV(explosionBytes)
	audioState := &balloon.AudioState{deviceID, audioSpec, explosionBytes}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	cloudNoise := noise.MakeNoise(noise.FBM, winWidth, winHeight, .009, 0.5, 3, 3)
	cloudGradient := getGradient(rgba{0, 0, 255}, rgba{255, 255, 255})
	scale := noise.CalcScale(255.0)
	cloudNoise.Rescale(scale)
	cloudPixels := make([]byte, winWidth*winHeight*4)
	drawNoise(cloudNoise, cloudGradient, cloudPixels)
	cloudTexture := pixelsToTexture(renderer, cloudPixels, winWidth, winHeight)

	imgs := loadImages("images", "balloon_")
	balloons := loadBalloons(renderer, imgs, 25)

	currentMouseState := balloon.GetMouseState()
	previousMouseState := currentMouseState
	var elapsedTime float32
	running := true
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

		renderer.Copy(cloudTexture, nil, nil)

		balloon.UpdateBalloons(balloons, elapsedTime, currentMouseState,
			previousMouseState, audioState, winWidth, winHeight, winDepth)
		sort.Sort(balloon.Slice(balloons))
		for _, b := range balloons {
			b.Draw(renderer)
		}

		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Milliseconds())
		// // fmt.Println("ms per frame", elapsedTime)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Milliseconds())
		}

		previousMouseState = currentMouseState
	}
}

func loadImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	return img
}

func loadImages(dirname, pattern string) []image.Image {
	fileInfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	var imgs []image.Image
	for _, fileInfo := range fileInfos {
		if strings.Contains(fileInfo.Name(), pattern) {
			filename := filepath.Join("images", fileInfo.Name())
			imgs = append(imgs, loadImage(filename))
		}
	}

	return imgs
}

func imgToTexture(renderer *sdl.Renderer, img image.Image, w, h int) *sdl.Texture {
	pixels := make([]byte, w*h*4)

	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[i] = byte(r / 256)
			pixels[i+1] = byte(g / 256)
			pixels[i+2] = byte(b / 256)
			pixels[i+3] = byte(a / 256)
			i += 4
		}
	}
	return pixelsToTexture(renderer, pixels, w, h)
}

func loadBalloons(renderer *sdl.Renderer, imgs []image.Image, numOfBalloons int) []*balloon.Balloon {
	numOfImages := len(imgs)
	balloonTextures := make([]*sdl.Texture, numOfImages)

	for i, img := range imgs {
		w := img.Bounds().Max.X
		h := img.Bounds().Max.Y
		tex := imgToTexture(renderer, img, w, h)
		err := tex.SetBlendMode(sdl.BLENDMODE_BLEND)
		if err != nil {
			panic(err)
		}

		balloonTextures[i] = tex
	}

	balloons := make([]*balloon.Balloon, numOfBalloons)
	for i := range balloons {
		// create equal no. of balloon textures
		t := balloonTextures[i%numOfImages]
		pos := vector.Vector{
			X: rand.Float32() * float32(winWidth),
			Y: rand.Float32() * float32(winHeight),
			Z: rand.Float32() * float32(winDepth),
		}
		velocity := vector.Vector{
			X: rand.Float32()*0.5 - 0.25,
			Y: rand.Float32()*0.5 - 0.25,
			Z: rand.Float32()*0.5 - .25/2,
		}

		explosionImg := loadImage("images/explosion.png")
		b := explosionImg.Bounds()
		explosionTex := imgToTexture(renderer, explosionImg, b.Max.X, b.Max.Y)
		err := explosionTex.SetBlendMode(sdl.BLENDMODE_BLEND)
		if err != nil {
			panic(err)
		}
		explosionAnimation := &balloon.Animation{false, false, 30, time.Now(), explosionTex}
		balloons[i] = balloon.NewBalloon(t, pos, velocity, explosionAnimation)
	}
	return balloons
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

func lerp(b1, b2 byte, pct float32) byte {
	return uint8(float32(b1) + pct*(float32(b2)-float32(b1)))
}

func colorLerp(c1, c2 rgba, pct float32) rgba {
	return rgba{
		lerp(c1.r, c2.r, pct),
		lerp(c1.g, c2.g, pct),
		lerp(c1.b, c2.b, pct),
	}
}

func getGradient(c1, c2 rgba) []rgba {
	result := make([]rgba, 256)

	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}

	return result
}

func getDualGradient(c1, c2, c3, c4 rgba) []rgba {
	result := make([]rgba, 256)

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

// drawNoise draws noise to the pixels buffer
func drawNoise(noise []float32, gradient []rgba, pixels []byte) {
	for i := range noise {
		c := gradient[clamp(0, 255, int(noise[i]))]
		p := i * 4
		pixels[p] = c.r
		pixels[p+1] = c.g
		pixels[p+2] = c.b
	}
}
