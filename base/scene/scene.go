package scene

import (
	"gotracer/base/scene/primitives/triangle"
	"gotracer/base/scene/primitives/vertex"
)

type Metadata struct {
	FormatVersion float32
	Objects       int32
	Geometries    int32
	Materials     int32
	Textures      int32
}

type Scene struct {
	Triangles []triangle.Triangle
	Vertices  []vertex.Vertex
	Info      Metadata
}
