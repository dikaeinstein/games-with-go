package game

var nums = [][]byte{
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		0, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
}

// Score represents a player score
type Score int

// Scores returns a slice of byte slice(scores represented as bits)
func Scores() [][]byte {
	return nums
}

// Draw renders the score in the pixels buffer
func (s Score) Draw(pos Pos, color Color, size int, pixels []byte) {
	startX := int(pos.X) - size*3/2
	startY := int(pos.Y) - size*5/2

	scores := Scores()
	for i, v := range scores[s] {
		if v == 1 {
			for y := startY; y < startY+int(size); y++ {
				for x := startX; x < startX+int(size); x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += int(size)
		if (i+1)%3 == 0 {
			// move down startY
			startY += int(size)
			// move back startX
			startX -= (int(size) * 3)
		}
	}
}
