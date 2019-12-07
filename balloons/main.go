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

func (t *texture) draw(p pos, w, h int, pixels []byte) {
	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			screenX := x + int(p.x)
			screenY := y + int(p.y)
			if screenX >= 0 && screenX < w && screenY >= 0 && screenY < h {
				texIndex := y*t.pitch + x*4
				screenIndex := (screenY*w + screenX) * 4

				pixels[screenIndex] = t.pixels[texIndex]
				pixels[screenIndex+1] = t.pixels[texIndex+1]
				pixels[screenIndex+2] = t.pixels[texIndex+2]
				pixels[screenIndex+3] = t.pixels[texIndex+3]
			}
		}
	}
}

func (t *texture) drawAlpha(p pos, w, h int, pixels []byte) {
	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			screenX := x + int(p.x)
			screenY := y + int(p.y)
			if screenX >= 0 && screenX < w && screenY >= 0 && screenY < h {
				texIndex := y*t.pitch + x*4
				screenIndex := (screenY*w + screenX) * 4

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

	window, err := sdl.CreateWindow("Ballons", sdl.WINDOWPOS_UNDEFINED,
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
	imgs := loadImages("images", "balloon_")
	balloonTextures := loadBalloons(imgs, winWidth, winHeight)

	for i, t := range balloonTextures {
		v := float32(i) * 40
		t.drawAlpha(pos{v, v}, winWidth, winHeight, pixels)
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

func loadBalloon(img image.Image, w, h int) *texture {
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

func loadBalloons(imgs []image.Image, w, h int) []*texture {
	textures := make([]*texture, len(imgs))
	for i, img := range imgs {
		textures[i] = loadBalloon(img, w, h)
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
