package vertex

import (
	"gotracer/base/vector"
)

type Vertex struct {
	vector.Vec4f
}

type Normal struct {
	vector.Vec4f
}

type UV struct {
	vector.Vec3f
}

func NewVertex(x, y, z float32) (vertex *Vertex) {
	vertex = new(Vertex)

	vertex.X = x
	vertex.Y = y
	vertex.Z = z
	vertex.W = 1.0

	return
}

func NewNormal(x, y, z float32) (normal *Normal) {
	normal = new(Normal)

	normal.X = x
	normal.Y = y
	normal.Z = z
	normal.W = 1.0

	return
}
