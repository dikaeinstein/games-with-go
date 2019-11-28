package noise

import (
	"github.com/dikaeinstein/games-with-go/simplex"
)

// Noise is a slice of 2D simplex noise values
type Noise []float32

type Color struct {
	R, G, B byte
}

// MakeNoise creates a slice 2D simplex Noise
func MakeNoise(w, h int, frequency, lacunarity, gain float32, octaves int) Noise {
	noise := make([]float32, w*h)

	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// noise[i] = simplex.SNoise2(float32(x)/100.0, float32(y)/100.0)
			noise[i] = turbulence(
				float32(x), float32(y),
				frequency, lacunarity,
				gain, octaves,
			)

			if noise[i] < min {
				min = noise[i]
			} else if noise[i] > max {
				max = noise[i]
			}
			i++
		}
	}

	return noise
}

var min = float32(9999.0)
var max = float32(-9999.0)

// RescaleAndDraw scales the noise values by the given scale
// and draws them to the pixels buffer
func (noise Noise) RescaleAndDraw(scale float32, gradient []Color, pixels []byte) {
	offset := min * scale

	for i := range noise {
		noise[i] = noise[i]*scale - offset
		c := gradient[clamp(0, 255, int(noise[i]))]
		p := i * 4
		pixels[p] = c.R
		pixels[p+1] = c.G
		pixels[p+2] = c.B
	}
}

// CalcScale calculates and returns the scale value using the give value
// scale = value / (max-min)
func CalcScale(value float32) float32 {
	return value / (max - min)
}

func fBm2(x, y, frequency, lacunarity, gain float32, octaves int) float32 {
	var sum float32
	amplitude := float32(1.0)

	for i := 0; i < octaves; i++ {
		sum += simplex.SNoise2(x*frequency, y*frequency) * amplitude
		frequency *= lacunarity
		amplitude *= gain
	}

	return sum
}

func turbulence(x, y, frequency, lacunarity, gain float32, octaves int) float32 {
	var sum float32
	amplitude := float32(1.0)

	for i := 0; i < octaves; i++ {
		n := simplex.SNoise2(x*frequency, y*frequency) * amplitude
		if n < 0 {
			n = -1.0 * n
		}
		sum += n
		frequency *= lacunarity
		amplitude *= gain
	}

	return sum
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
