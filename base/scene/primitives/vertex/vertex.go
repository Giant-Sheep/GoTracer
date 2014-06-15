package vertex

import (
	"gotracer/base/vector"
)

type Vertex struct {
	Pos vector.Vec4f
}

type Normal struct {
	Dir vector.Vec4f
}

type UV struct {
	Pos vector.Vec3f
}

func NewVertex(x, y, z float32) (vertex *Vertex) {
	vertex = new(Vertex)

	vertex.Pos.X = x
	vertex.Pos.Y = y
	vertex.Pos.Z = z
	vertex.Pos.W = 1.0

	return
}

func NewNormal(x, y, z float32) (normal *Normal) {
	normal = new(Normal)

	normal.Dir.X = x
	normal.Dir.Y = y
	normal.Dir.Z = z
	normal.Dir.W = 1.0

	return
}
