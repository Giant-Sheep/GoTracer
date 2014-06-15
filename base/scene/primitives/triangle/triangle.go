package triangle

import (
	"gotracer/base/scene/primitives/vertex"
)

type Triangle struct {
	Normals  [3]*vertex.Normal
	Vertices [3]*vertex.Vertex
}
