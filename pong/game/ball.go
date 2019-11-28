package game

// Ball represents the ball in the pong game
type Ball struct {
	Pos
	radius               float32
	xVelocity, yVelocity float32
	color                Color
}

// NewBall creates an instance of a Ball
func NewBall(pos Pos, radius, xVelocity, yVelocity float32, color Color) *Ball {
	return &Ball{
		pos,
		radius,
		xVelocity, yVelocity,
		color,
	}
}

// Draw renders the ball in the pixels buffer
func (b *Ball) Draw(pixels []byte) {
	bRadius := int(b.radius)
	for y := -bRadius; y < bRadius; y++ {
		for x := -bRadius; x < bRadius; x++ {
			if x*x+y*y < bRadius*bRadius {
				setPixel(int(b.X)+x, int(b.Y)+y, b.color, pixels)
			}
		}
	}
}

// Update updates the position of the ball based on collision with paddles
func (b *Ball) Update(leftPaddle, rightPaddle *Paddle, elapsedTime float32) {
	b.X += b.xVelocity * elapsedTime
	b.Y += b.yVelocity * elapsedTime

	if b.Y-b.radius < 0 || b.Y+b.radius > float32(winHeight) {
		// invert the y component of the ball velocity
		b.yVelocity = -b.yVelocity
	}

	if b.X < 0 {
		rightPaddle.score++
		b.Pos = GetCenter()
		state = StateStart
	} else if int(b.X) > winWidth {
		leftPaddle.score++
		b.Pos = GetCenter()
		state = StateStart
	}

	if b.X-b.radius < leftPaddle.X+leftPaddle.w/2 {
		if b.Y > leftPaddle.Y-leftPaddle.h/2 && b.Y < leftPaddle.Y+leftPaddle.h/2 {
			b.xVelocity = -b.xVelocity
			// minimum translation vector after collision
			b.X = leftPaddle.X + leftPaddle.w/2 + b.radius
		}
	}

	if b.X+b.radius > rightPaddle.X-rightPaddle.w/2 {
		if b.Y > rightPaddle.Y-rightPaddle.h/2 && b.Y < rightPaddle.Y+rightPaddle.h/2 {
			b.xVelocity = -b.xVelocity
			// minimum translation vector
			b.X = rightPaddle.X - rightPaddle.w/2 - b.radius
		}
	}
}
