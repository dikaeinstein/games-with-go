package vector

import "math"

// Vector represents a 3D vector
type Vector struct {
	X, Y, Z float32
}

// Length returns the length/magnitude of the vector
func (v Vector) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Add adds the two given vectors
func Add(v1, v2 Vector) Vector {
	return Vector{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

// Multiply is the dot product of the two given vectors
func Multiply(v Vector, factor float32) Vector {
	return Vector{
		X: v.X * factor,
		Y: v.Y * factor,
		Z: v.Z * factor,
	}
}

// Distance is the distance between the two vectors
func Distance(v1, v2 Vector) float32 {
	dx := v1.X - v2.X
	dy := v1.Y - v2.Y
	dz := v1.Z - v1.Z

	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

// DistanceSquared is the sum of the square of the given vectors
func DistanceSquared(v1, v2 Vector) float32 {
	dx := v1.X - v2.X
	dy := v1.Y - v2.Y
	dz := v1.Z - v1.Z

	return dx*dx + dy*dy + dz*dz
}

// Normalize the given vector i.e convert it into a unit vector
func Normalize(v Vector) Vector {
	len := v.Length()
	return Vector{
		X: v.X / len,
		Y: v.Y / len,
		Z: v.Z / len,
	}
}
