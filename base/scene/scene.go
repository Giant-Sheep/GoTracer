package scene

import (
	"gotracer/base/quaternion"
	"gotracer/base/scene/primitives/triangle"
	"gotracer/base/scene/primitives/vertex"
	"gotracer/base/vector"
)

type Metadata struct {
	FormatVersion float32
	Objects       int32
	Geometries    int32
	Materials     int32
	Textures      int32
}

// TODO: Containt only fields for directional light for the time being
type Light struct {
	Type      string
	Colour    vector.Vec3ui8
	Intensity float32
	Direction vector.Vec3f
	Target    string
}

type Camera struct {
	Type     string
	FOV      float32
	Near     float32
	Far      float32
	Position vector.Vec3f
}

type Geometry struct {
	ID       string
	Metadata struct {
		Vertices int32
		Normals  int32
		Colors   int32
		Faces    int32
		//UVs	[] TODO: No textures yet, not sure of the format of this
	}
	BB struct {
		Min vector.Vec3f
		Max vector.Vec3f
	}
	Scale     float32
	Materials *Material
	Vertices  []*vertex.Vertex
	Normals   []*vertex.Normal
	//Colors []*vertex.Colors TODO: Not supported yet
	//UVs []*vertex.UV TODO: Not supported yet
}

type Material struct {
	ID                 string
	Type               string
	Color              vector.Vec3ui8
	Ambient            vector.Vec3ui8
	Emissive           vector.Vec3ui8
	Transparent        bool
	Reflectivity       float32
	Opacity            float32
	Wireframe          bool
	WireframeLinewidth int32
}

type Object struct {
	ID       string
	Geometry *Geometry
	Material *Material
	Position vector.Vec4f
	Rotation quaternion.Quaternion
	Scale    vector.Vec3f
	Visible  bool
	Children []Object
}

type Scene struct {
	Triangles  []triangle.Triangle
	Vertices   []vertex.Vertex
	Normals    []vertex.Normal
	UVs        []vertex.UV
	Lights     []Light
	Materials  []Material
	Objects    []Object
	Geometries []Geometry
	Info       Metadata
	Camera     Camera
}

func NewScene() (scene *Scene) {
	scene = new(Scene)
	scene.Lights = make([]Light, 0, 10)
	scene.Triangles = make([]triangle.Triangle, 0, 100000)
	scene.Materials = make([]Material, 0, 1000)
	scene.Objects = make([]Object, 0, 50)
	scene.Geometries = make([]Geometry, 0, 1000)
	return
}

func NewObject() (object *Object) {
	object = new(Object)
	object.Rotation = quaternion.Quaternion{0, 0, 0, 1}
	object.Position = vector.Vec4f{0, 0, 0, 1}
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

func (world *Scene) GetGeometry(id string) *Geometry {
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

func (world *Scene) GetMaterial(id string) *Material {
	for _, material := range world.Materials {
		if material.ID == id {
			return &material
			break
		}
	}

	return nil
}
