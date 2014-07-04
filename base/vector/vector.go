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

type Vec3ui8 struct {
	X, Y, Z uint8
}

type Vec4i struct {
	X, Y, Z, W int32
}

type Vec4ui struct {
	X, Y, Z, W uint32
}

type Vec4ui8 struct {
	X, Y, Z, W uint8
}

func (v *Vec2i) Min() (min int) {
	min = int(math.Min(float64(v.X), float64(v.Y)))

	return
}

func (v *Vec2i) Max() (max int) {
	max = int(math.Max(float64(v.X), float64(v.Y)))

	return
}

func (a *Vec4f) Cross(b Vec4f) (c Vec4f) {
	c.X = a.Y*b.Z - a.Z*b.Y
	c.Y = a.Z*b.X - a.X*b.Z
	c.Z = a.X*b.Y - a.Y*b.X
	c.W = 1.0

	return
}

func (a *Vec4f) Dot(b Vec4f) (dot float32) {
	dot = a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W

	return
}

func (a *Vec4f) Subtract(b *Vec4f) (c Vec4f) {
	c.X = a.X - b.X
	c.Y = a.Y - b.Y
	c.Z = a.Z - b.Z
	c.W = a.W

	return
}

func (a *Vec4f) Scale(b float32) (c Vec4f) {
	c.X = a.X * b
	c.Y = a.Y * b
	c.Z = a.Z * b
	c.W = a.W

	return
}

func (a *Vec4f) Add(b *Vec4f) (c Vec4f) {
	c.X = a.X + b.X
	c.Y = a.Y + b.Y
	c.Z = a.Z + b.Z
	c.W = a.W

	return
}

func (vec *Vec4f) Normalize() {
	length := float32(math.Sqrt(math.Pow(float64(vec.X), 2) + math.Pow(float64(vec.Y), 2) + math.Pow(float64(vec.Z), 2)))
	vec.X /= length
	vec.Y /= length
	vec.Z /= length
	return
}
