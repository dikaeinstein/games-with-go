package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type rgba struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type texture struct {
	pixels      []byte
	w, h, pitch int
}

func (t *texture) draw(p pos, pixels []byte) {
	for y := 0; y < t.h; y++ {
		screenY := y + int(p.y)
		for x := 0; x < t.w; x++ {
			screenX := x + int(p.x)
			if screenX >= 0 && screenX < t.w && screenY >= 0 && screenY < t.h {
				texIndex := y*t.pitch + x*4
				screenIndex := (screenY*t.w + screenX) * 4

				pixels[screenIndex] = t.pixels[texIndex]
				pixels[screenIndex+1] = t.pixels[texIndex+1]
				pixels[screenIndex+2] = t.pixels[texIndex+2]
				pixels[screenIndex+3] = t.pixels[texIndex+3]
			}
		}
	}
}

func (t *texture) drawScaled(p pos, scaleX, scaleY float32, pixels []byte) {
	newWidth := int(scaleX * float32(t.w))
	newHeight := int(scaleY * float32(t.h))

	for y := 0; y < newHeight; y++ {
		fy := float32(y) / float32(newHeight) * float32(t.h-1)
		iy := int(fy)
		screenY := int(fy*scaleY) + int(p.y)
		screenIndex := screenY*t.w*4 + int(p.x)*4
		for x := 0; x < newWidth; x++ {
			fx := float32(x) / float32(newWidth) * float32(t.w-1)
			screenX := int(fx*scaleX) + int(p.x)
			if screenX >= 0 && screenX < t.w && screenY >= 0 && screenY < t.h {
				texIndex := iy*t.pitch + int(fx)*4
				// screenIndex := (screenY*t.w + screenX) * 4

				pixels[screenIndex] = t.pixels[texIndex]
				pixels[screenIndex+1] = t.pixels[texIndex+1]
				pixels[screenIndex+2] = t.pixels[texIndex+2]
				pixels[screenIndex+3] = t.pixels[texIndex+3]
			}
		}
	}
}

func (t *texture) drawAlpha(p pos, pixels []byte) {
	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			screenX := x + int(p.x)
			screenY := y + int(p.y)
			if screenX >= 0 && screenX < t.w && screenY >= 0 && screenY < t.h {
				texIndex := y*t.pitch + x*4
				screenIndex := (screenY*t.w + screenX) * 4

				srcR := int(t.pixels[texIndex])
				srcG := int(t.pixels[texIndex+1])
				srcB := int(t.pixels[texIndex+2])
				srcA := int(t.pixels[texIndex+3])

				dstR := int(pixels[screenIndex])
				dstG := int(pixels[screenIndex+1])
				dstB := int(pixels[screenIndex+2])

				outR := (srcR*255 + dstR*(255-srcA)) / 255
				outG := (srcG*255 + dstG*(255-srcA)) / 255
				outB := (srcB*255 + dstB*(255-srcA)) / 255

				pixels[screenIndex] = byte(outR)
				pixels[screenIndex+1] = byte(outG)
				pixels[screenIndex+2] = byte(outB)
				// pixels[screenIndex+3] =
			}
		}
	}
}

const winWidth = 800
const winHeight = 600

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Balloons", sdl.WINDOWPOS_UNDEFINED,
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

	// cloudNoise := noise.MakeNoise(noise.FBM, winWidth, winHeight, .009, 0.5, 3, 3)
	// cloudGradient := getGradient(rgba{0, 0, 255}, rgba{255, 255, 255})
	// scale := noise.CalcScale(255.0)
	// cloudNoise.Rescale(scale)
	// cloudPixels := make([]byte, winWidth*winHeight*4)
	// drawNoise(cloudNoise, cloudGradient, cloudPixels)
	// cloudTexture := texture{cloudPixels, winWidth, winHeight, winWidth * 4}

	pixels := make([]byte, winWidth*winHeight*4)
	imgs := loadImages("images", "balloon_")
	balloonTextures := loadBalloons(imgs, winWidth, winHeight)

	// cloudTexture.draw(pos{0, 0}, pixels)
	for i, t := range balloonTextures {
		v := float32(i) * 40
		t.drawScaled(pos{v, v}, 1, 1, pixels)
	}

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

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(16)
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

func loadBalloon(img image.Image) *texture {
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

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

	return &texture{pixels, w, h, w * 4}
}

func loadBalloons(imgs []image.Image) []*texture {
	textures := make([]*texture, len(imgs))
	for i, img := range imgs {
		textures[i] = loadBalloon(img)
	}
	return textures
}

func setPixel(x, y int, c rgba, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
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
