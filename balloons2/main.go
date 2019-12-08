package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dikaeinstein/games-with-go/noise"
	"github.com/veandco/go-sdl2/sdl"
)

type rgba struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type balloon struct {
	tex *sdl.Texture
	pos
	scale float32
	w, h  int
}

func (b *balloon) draw(renderer *sdl.Renderer) {
	newWidth := int32(float32(b.w) * b.scale)
	newHeight := int32(float32(b.h) * b.scale)

	x := int32(b.x - float32(b.w)/2)
	y := int32(b.y - float32(b.h)/2)

	rect := &sdl.Rect{X: x, Y: y, W: newWidth, H: newHeight}
	renderer.Copy(b.tex, nil, rect)
}

const winWidth = 800
const winHeight = 600

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Balloons 2", sdl.WINDOWPOS_UNDEFINED,
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

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	cloudNoise := noise.MakeNoise(noise.FBM, winWidth, winHeight, .009, 0.5, 3, 3)
	cloudGradient := getGradient(rgba{0, 0, 255}, rgba{255, 255, 255})
	scale := noise.CalcScale(255.0)
	cloudNoise.Rescale(scale)
	cloudPixels := make([]byte, winWidth*winHeight*4)
	drawNoise(cloudNoise, cloudGradient, cloudPixels)
	cloudTexture := pixelsToTexture(renderer, cloudPixels, winWidth, winHeight)

	imgs := loadImages("images", "balloon_")
	balloons := loadBalloons(renderer, imgs)

	dir := 1
	var elapsedTime float32
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

		renderer.Copy(cloudTexture, nil, nil)

		for _, b := range balloons {
			b.draw(renderer)
		}

		balloons[1].x += float32(1 * dir)
		if balloons[1].x > 400 || balloons[1].x < 0 {
			dir = dir * -1
		}

		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Milliseconds())
		fmt.Println("ms per frame", elapsedTime)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Milliseconds())
		}
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

func loadBalloonTexture(renderer *sdl.Renderer, img image.Image, w, h int) *sdl.Texture {
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

func loadBalloons(renderer *sdl.Renderer, imgs []image.Image) []balloon {
	balloons := make([]balloon, len(imgs))

	for i, img := range imgs {
		w := img.Bounds().Max.X
		h := img.Bounds().Max.Y
		tex := loadBalloonTexture(renderer, img, w, h)
		err := tex.SetBlendMode(sdl.BLENDMODE_BLEND)
		if err != nil {
			panic(err)
		}

		// v := float32(i * 120)
		balloons[i] = balloon{tex,
			pos{float32(i * 120), float32(i * 120)}, float32(1+i) / 2, w, h}
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
