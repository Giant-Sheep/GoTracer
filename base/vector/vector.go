package vector

import (
	"math"
)

type Vec3f struct {
	X, Y, Z float32
}

type Vec4f struct {
	X, Y, Z, W float32
}

type Vec2i struct {
	X, Y int32
}

type Vec3i struct {
	X, Y, Z int32
}

func (v *Vec2i) Min() (min int) {
	min = int(math.Min(float64(v.X), float64(v.Y)))

	return
}

func (v *Vec2i) Max() (max int) {
	max = int(math.Max(float64(v.X), float64(v.Y)))

	return
}