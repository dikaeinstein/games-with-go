package noise

import (
	"math"

	"github.com/dikaeinstein/games-with-go/simplex"
)

// Noise is a slice of 2D simplex noise values
type Noise []float32

// Type indicates which noise MakeNoise will generate
type Type uint

const (
	// FBM is the fractal brownian motion noise
	FBM Type = iota
	// TURBULENCE is the turbulence fractal noise
	TURBULENCE
)

// MakeNoise creates a slice 2D simplex Noise
func MakeNoise(noiseType Type, w, h int, frequency, lacunarity, gain float32, octaves int) Noise {
	noise := make([]float32, w*h)

	// numCPUs := runtime.NumCPU()
	// batchSize := len(noise) / numCPUs
	// var wg sync.WaitGroup
	// wg.Add(numCPUs)

	// for i := 0; i < numCPUs; i++ {
	// 	go func(i int) {
	// 		defer wg.Done()
	// 		start := i * batchSize
	// 		end := start + batchSize - 1
	// 		for j := start; j < end; j++ {
	// 			x := j % w
	// 			y := (j - x) / w
	// 			noise[j] = Turbulence(
	// 				float32(x), float32(y),
	// 				frequency, lacunarity,
	// 				gain, octaves,
	// 			)

	// 			if noise[j] < min {
	// 				min = noise[j]
	// 			} else if noise[j] > max {
	// 				max = noise[j]
	// 			}
	// 		}
	// 	}(i)
	// }

	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// noise[i] = simplex.SNoise2(float32(x)/100.0, float32(y)/100.0)
			if noiseType == TURBULENCE {
				noise[i] = Turbulence(
					float32(x), float32(y),
					frequency, lacunarity,
					gain, octaves,
				)
			} else {
				noise[i] = Fbm2(
					float32(x), float32(y),
					frequency, lacunarity,
					gain, octaves,
				)
			}

			if noise[i] < min {
				min = noise[i]
			} else if noise[i] > max {
				max = noise[i]
			}
			i++
		}
	}

	// wg.Wait()
	return noise
}

var min = float32(math.MaxFloat32)
var max = float32(-math.MaxFloat32)

// Rescale scales the noise values by the given scale in place
func (noise Noise) Rescale(scale float32) {
	offset := min * scale

	for i := range noise {
		noise[i] = noise[i]*scale - offset
	}
}

// CalcScale calculates and returns the scale value using the give value.
// `scale = value / (max-min)`
func CalcScale(value float32) float32 {
	return value / (max - min)
}

// Fbm2 generates fractal brownian motion noise
func Fbm2(x, y, frequency, lacunarity, gain float32, octaves int) float32 {
	var sum float32
	amplitude := float32(1.0)

	for i := 0; i < octaves; i++ {
		sum += simplex.SNoise2(x*frequency, y*frequency) * amplitude
		frequency *= lacunarity
		amplitude *= gain
	}

	return sum
}

// Turbulence generates turbulence fractal noise
func Turbulence(x, y, frequency, lacunarity, gain float32, octaves int) float32 {
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
