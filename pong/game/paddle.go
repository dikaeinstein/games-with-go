package game

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// Paddle represents the paddle in the pong game
type Paddle struct {
	Pos
	w, h  float32
	speed float32
	score Score
	color Color
}

// NewPaddle creates an instance of a Paddle
func NewPaddle(pos Pos, w, h, speed float32, score Score, color Color) *Paddle {
	return &Paddle{
		pos,
		w, h,
		speed,
		score,
		color,
	}
}

// Draw renders the paddle in the pixels buffer
func (p *Paddle) Draw(pixels []byte) {
	startX := int(p.X - p.w/2)
	startY := int(p.Y - p.h/2)

	var x, y int
	for y = 0; y < int(p.h); y++ {
		for x = 0; x < int(p.w); x++ {
			setPixel(startX+x, startY+y, p.color, pixels)
		}
	}

	scoreX := Lerp(p.X, GetCenter().X, 0.4)
	p.score.Draw(Pos{X: scoreX, Y: 70}, p.color, 5, pixels)
}

// Update updates the position of the paddle based on input(up/down) from the keyboard
func (p *Paddle) Update(keyboardStates []uint8, controllerAxis int16, elapsedTime float32) {
	if keyboardStates[sdl.SCANCODE_UP] != 0 {
		p.Y -= p.speed * elapsedTime
	}
	if keyboardStates[sdl.SCANCODE_DOWN] != 0 {
		p.Y += p.speed * elapsedTime
	}

	if math.Abs(float64(controllerAxis)) > 1500 {
		p.Y += p.speed * float32(controllerAxis) / 32767.0 * elapsedTime
	}
}

// AIUpdate updates the position of the paddle based on the AI move
func (p *Paddle) AIUpdate(ball *Ball, elapsedTime float32) {
	p.Y = ball.Y
}

// GetScore returns the current score of this player/paddle
func (p *Paddle) GetScore() Score {
	return p.score
}

// ResetScore resets the paddle/player score to zero
func (p *Paddle) ResetScore() {
	p.score = 0
}
