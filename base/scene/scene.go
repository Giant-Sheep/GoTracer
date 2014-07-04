package scene

import (
	"gotracer/base/scene/objects"
	"gotracer/base/scene/primitives/triangle"
	"gotracer/base/scene/primitives/vertex"
	"sync"
)

type Metadata struct {
	FormatVersion float32
	Objects       int32
	Geometries    int32
	Materials     int32
	Textures      int32
}

type Scene struct {
	sync.Mutex
	Triangles  []triangle.Triangle
	Vertices   []vertex.Vertex
	Normals    []vertex.Normal
	UVs        []vertex.UV
	Lights     []objects.Light
	Materials  []objects.Material
	Objects    []objects.Object
	Geometries []objects.Geometry
	Info       Metadata
	Camera     objects.Camera
}

func NewScene() (scene *Scene) {
	scene = new(Scene)
	scene.Lights = make([]objects.Light, 0, 10)
	scene.Triangles = make([]triangle.Triangle, 0, 100000)
	scene.Materials = make([]objects.Material, 0, 1000)
	scene.Objects = make([]objects.Object, 0, 50)
	scene.Geometries = make([]objects.Geometry, 0, 1000)
	return
}

func (world *Scene) HasGeometry(id string) bool {
	for _, geometry := range world.Geometries {
		if geometry.ID == id {
			return true
			break
		}
	}

	return false
}

func (world *Scene) GetGeometry(id string) *objects.Geometry {
	for _, geometry := range world.Geometries {
		if geometry.ID == id {
			return &geometry
			break
		}
	}

	return nil
}

func (world *Scene) HasMaterial(id string) bool {
	for _, material := range world.Materials {
		if material.ID == id {
			return true
			break
		}
	}

	return false
}

func (world *Scene) GetMaterial(id string) *objects.Material {
	for _, material := range world.Materials {
		if material.ID == id {
			return &material
			break
		}
	}

	return nil
}
