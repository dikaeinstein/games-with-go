package balloon

import (
	"math"
	"time"

	"github.com/dikaeinstein/games-with-go/vector"
	"github.com/veandco/go-sdl2/sdl"
)

type Animation struct {
	Running   bool
	Complete  bool
	Interval  float32
	StartTime time.Time
	Texture   *sdl.Texture
}

type AudioState struct {
	DeviceID       sdl.AudioDeviceID
	AudioSPec      sdl.AudioSpec
	ExplosionBytes []byte
}

type MouseState struct {
	LeftButton, RightButton bool
	X, Y                    int
}

func GetMouseState() MouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	LeftButton := mouseButtonState & sdl.ButtonLMask()
	RightButton := mouseButtonState & sdl.ButtonRMask()
	return MouseState{
		!(LeftButton == 0),
		!(RightButton == 0),
		int(mouseX),
		int(mouseY),
	}
}

type Balloon struct {
	tex                *sdl.Texture
	pos, velocity      vector.Vector
	w, h               int
	explosionAnimation *Animation
}

func NewBalloon(t *sdl.Texture, pos, velocity vector.Vector, explosionAnimation *Animation) *Balloon {
	_, _, w, h, err := t.Query()
	if err != nil {
		panic(err)
	}

	return &Balloon{t, pos, velocity, int(w), int(h), explosionAnimation}
}

func (b *Balloon) Draw(renderer *sdl.Renderer) {
	scale := b.Scale()
	newWidth := int32(float32(b.w) * scale)
	newHeight := int32(float32(b.h) * scale)

	x := int32(b.pos.X - float32(b.w)/2)
	y := int32(b.pos.Y - float32(b.h)/2)

	rect := &sdl.Rect{X: x, Y: y, W: newWidth, H: newHeight}
	renderer.Copy(b.tex, nil, rect)

	if b.explosionAnimation.Running {
		numOfAnimations := 16
		animationElapsedTime := float32(time.Since(b.explosionAnimation.StartTime).Milliseconds())
		animationIndex := numOfAnimations - 1 - int(animationElapsedTime/b.explosionAnimation.Interval)
		animationX := animationIndex % 4
		animationY := 64 * ((animationIndex - animationX) / 4)
		animationX *= 64

		animationRect := &sdl.Rect{X: int32(animationX), Y: int32(animationY), W: 64, H: 64}
		rect.X -= rect.W / 2
		rect.Y -= rect.H / 2
		rect.W *= 2
		rect.H *= 2
		renderer.Copy(b.explosionAnimation.Texture, animationRect, rect)
	}
}

func UpdateBalloons(balloons []*Balloon, elapsedTime float32,
	current, previous MouseState, audioState *AudioState, w, h, d int) []*Balloon {
	numOfAnimations := 16
	balloonClicked := false
	balloonExploaded := false

	for i := len(balloons) - 1; i > 0; i-- {
		b := balloons[i]
		animationElapsedTime := float32(time.Since(b.explosionAnimation.StartTime).Milliseconds())
		animationIndex := numOfAnimations - 1 - int(animationElapsedTime/b.explosionAnimation.Interval)

		if b.explosionAnimation.Running {
			if animationIndex < 0 {
				b.explosionAnimation.Running = false
				b.explosionAnimation.Complete = true
				balloonExploaded = true
			}
		}

		if !balloonClicked && !previous.LeftButton && current.LeftButton {
			x, y, r := b.Circle()
			dx := float32(current.X) - x
			dy := float32(current.Y) - y
			dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
			if dist < r {
				balloonClicked = true
				sdl.ClearQueuedAudio(audioState.DeviceID)
				sdl.QueueAudio(audioState.DeviceID, audioState.ExplosionBytes)
				sdl.PauseAudioDevice(audioState.DeviceID, false)
				b.explosionAnimation.Running = true
				b.explosionAnimation.StartTime = time.Now()
			}
		}

		// compute the new position for the ballon based on its current postion,
		// velocity and the elapsedTime for the previous frame
		p := vector.Add(b.pos, vector.Multiply(b.velocity, elapsedTime))
		if p.X < 0 || p.X > float32(w) {
			b.velocity.X = -b.velocity.X
		}
		if p.Y < 0 || p.Y > float32(h) {
			b.velocity.Y = -b.velocity.Y
		}
		if p.Z < 0 || p.Z > float32(d) {
			b.velocity.Z = -b.velocity.Z
		}

		b.pos = vector.Add(b.pos, vector.Multiply(b.velocity, elapsedTime))
	}

	if balloonExploaded {
		filteredBalloons := balloons[:0]
		// Filtering without allocating
		for _, b := range balloons {
			if !b.explosionAnimation.Complete {
				filteredBalloons = append(filteredBalloons, b)
			}
		}
		balloons = filteredBalloons
	}

	return balloons
}

func (b *Balloon) Scale() float32 {
	return (b.pos.Z/200 + 1) / 2
}

func (b *Balloon) Circle() (x, y, r float32) {
	scale := b.Scale()
	x = b.pos.X
	y = b.pos.Y - 30*scale
	r = (float32(b.w) / 2) * scale

	return x, y, r
}

type Slice []*Balloon

func (bs Slice) Len() int {
	return len(bs)
}

func (bs Slice) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

func (bs Slice) Less(i, j int) bool {
	diff := bs[i].pos.Z - bs[j].pos.Z
	return diff < -0.5
}
