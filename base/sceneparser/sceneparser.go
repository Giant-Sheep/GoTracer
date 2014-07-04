package sceneparser

import (
	"encoding/json"
	"gotracer/base/quaternion"
	"gotracer/base/scene"
	"gotracer/base/scene/objects"
	"gotracer/base/scene/primitives/triangle"
	"gotracer/base/scene/primitives/vertex"
	"gotracer/base/util"
	"gotracer/base/vector"
	"io/ioutil"
	"log"
	"math"
	"strings"
)

var (
	decoder *json.Decoder
)

func parseMetadata(world *scene.Scene, metadata map[string]interface{}) {
	log.Println("Parsing Json Metadata")
	world.Info.FormatVersion = float32(metadata["formatVersion"].(float64))
	world.Info.Objects = int32(metadata["objects"].(float64))
	world.Info.Geometries = int32(metadata["geometries"].(float64))
	world.Info.Materials = int32(metadata["materials"].(float64))
	world.Info.Textures = int32(metadata["textures"].(float64))
}

func parsePhysicalObjects(world *scene.Scene, jsondata, items map[string]interface{}) (retList []objects.Object) {
	retList = make([]objects.Object, 0, 50)
	for k, v := range items {
		object := objects.NewObject()
		object.ID = k
		if v.(map[string]interface{})["geometry"] != nil {
			object.Geometry = parseGeometry(world, jsondata, v.(map[string]interface{})["geometry"].(string))
		}

		if v.(map[string]interface{})["material"] != nil {
			object.Material = parseMaterial(world, jsondata, v.(map[string]interface{})["material"].(string))
		}

		object.Position.X = float32(v.(map[string]interface{})["position"].([]interface{})[0].(float64))
		object.Position.Y = float32(v.(map[string]interface{})["position"].([]interface{})[1].(float64))
		object.Position.Z = float32(v.(map[string]interface{})["position"].([]interface{})[2].(float64))
		object.Position.W = 1
		object.Rotation = quaternion.NewQuaternion(vector.Vec3f{X: 0, Y: 0, Z: 0}, float32(1.0))
		object.Rotation.Mul(quaternion.NewQuaternion(vector.Vec3f{X: 0, Y: 1, Z: 0}, float32(v.(map[string]interface{})["rotation"].([]interface{})[0].(float64))))
		object.Rotation.Mul(quaternion.NewQuaternion(vector.Vec3f{X: 1, Y: 0, Z: 0}, float32(v.(map[string]interface{})["rotation"].([]interface{})[1].(float64))))
		object.Rotation.Mul(quaternion.NewQuaternion(vector.Vec3f{X: 0, Y: 0, Z: 1}, float32(v.(map[string]interface{})["rotation"].([]interface{})[2].(float64))))
		object.Scale.X = float32(v.(map[string]interface{})["scale"].([]interface{})[0].(float64))
		object.Scale.Y = float32(v.(map[string]interface{})["scale"].([]interface{})[1].(float64))
		object.Scale.Z = float32(v.(map[string]interface{})["scale"].([]interface{})[2].(float64))
		object.Visible = v.(map[string]interface{})["visible"].(bool)
		if val, ok := v.(map[string]interface{})["children"]; ok {
			object.Children = parsePhysicalObjects(world, jsondata, val.(map[string]interface{}))
		}
	}
	return retList
}

func parseObjects(world *scene.Scene, jsondata map[string]interface{}) {
	items := jsondata["objects"].(map[string]interface{})
	log.Println("Parsing Json Objects")
	for k, v := range items {
		if strings.Contains(k, "camera") {
			val := v.(map[string]interface{})
			world.Camera.Type = val["type"].(string)
			world.Camera.FOV = float32(val["fov"].(float64) * (2 * math.Pi / 360))
			world.Camera.Near = float32(val["near"].(float64))
			world.Camera.Far = float32(val["far"].(float64))
			// world.Camera.Position.X = float32(val["position"].([]interface{})[0].(float64))
			// world.Camera.Position.Y = float32(val["position"].([]interface{})[1].(float64))
			// world.Camera.Position.Z = float32(val["position"].([]interface{})[2].(float64))
			// Let's override the camera position and put it in a sensible position
			world.Camera.Position.X = 0
			world.Camera.Position.Y = 0
			world.Camera.Position.Z = float32((-1) * (800 / math.Tan((float64(world.Camera.FOV) / 2))))
			world.Camera.Position.W = 1
			world.Camera.Rotation = quaternion.Quaternion{X: 0, Y: 0, Z: 0, W: 1}
			delete(items, k)
		} else if strings.Contains(k, "light") {
			var light objects.Light
			val := v.(map[string]interface{})
			light.Type = val["type"].(string)
			light.Colour = util.Uint32toVec3ui8(uint32(val["color"].(float64)))
			light.Intensity = float32(val["intensity"].(float64))
			light.Direction.X = float32(val["direction"].([]interface{})[0].(float64))
			light.Direction.Y = float32(val["direction"].([]interface{})[1].(float64))
			light.Direction.Z = float32(val["direction"].([]interface{})[2].(float64))
			light.Target = val["target"].(string)
			world.Lights = append(world.Lights, light)
			delete(items, k)
		}
	}
	world.Objects = append(world.Objects, parsePhysicalObjects(world, jsondata, jsondata["objects"].(map[string]interface{}))...)
}

func parseMaterial(world *scene.Scene, jsondata map[string]interface{}, material_id string) (material *objects.Material) {
	log.Println("Parsing Json Material")
	materialdata := jsondata["materials"].(map[string]interface{})[material_id].(map[string]interface{})
	material = new(objects.Material)

	material.ID = material_id
	material.Type = materialdata["type"].(string)

	if materialdata["parameters"].(map[string]interface{})["materials"] == nil {
		material.Color = util.Uint32toVec3ui8(uint32(materialdata["parameters"].(map[string]interface{})["color"].(float64)))
		material.Ambient = util.Uint32toVec3ui8(uint32(materialdata["parameters"].(map[string]interface{})["ambient"].(float64)))
		material.Emissive = util.Uint32toVec3ui8(uint32(materialdata["parameters"].(map[string]interface{})["emissive"].(float64)))
		material.Transparent = materialdata["parameters"].(map[string]interface{})["transparent"].(bool)
		material.Reflectivity = float32(materialdata["parameters"].(map[string]interface{})["reflectivity"].(float64))
		material.Opacity = float32(materialdata["parameters"].(map[string]interface{})["opacity"].(float64))
		material.Wireframe = materialdata["parameters"].(map[string]interface{})["wireframe"].(bool)
		material.WireframeLinewidth = int32(materialdata["parameters"].(map[string]interface{})["wireframeLinewidth"].(float64))
	} else {
		materials := materialdata["parameters"].(map[string]interface{})["materials"].([]interface{})
		for _, v := range materials {
			id := v.(string)
			if world.HasMaterial(id) {
				material.Materials = append(material.Materials, world.GetMaterial(id))
			} else {
				material.Materials = append(material.Materials, parseMaterial(world, jsondata, id))
			}
		}
	}

	return
}

func parseGeometry(world *scene.Scene, jsondata map[string]interface{}, geometry_id string) (geometry *objects.Geometry) {
	log.Println("Parsing Json Geometry")
	embeds := jsondata["embeds"].(map[string]interface{})
	geometrydata := jsondata["geometries"].(map[string]interface{})[geometry_id].(map[string]interface{})

	geometry = new(objects.Geometry)
	geometry.ID = geometry_id
	if geometrydata["type"].(string) == "embedded" {
		embed := embeds[geometrydata["id"].(string)].(map[string]interface{})

		// Parse metadata
		geometry.Metadata.Vertices = int32(embed["metadata"].(map[string]interface{})["vertices"].(float64))
		geometry.Metadata.Normals = int32(embed["metadata"].(map[string]interface{})["normals"].(float64))
		geometry.Metadata.Colors = int32(embed["metadata"].(map[string]interface{})["colors"].(float64))
		geometry.Metadata.Faces = int32(embed["metadata"].(map[string]interface{})["faces"].(float64))

		// Parse bounding box
		geometry.BB.Min.X = float32(embed["boundingBox"].(map[string]interface{})["min"].([]interface{})[0].(float64))
		geometry.BB.Min.Y = float32(embed["boundingBox"].(map[string]interface{})["min"].([]interface{})[1].(float64))
		geometry.BB.Min.Z = float32(embed["boundingBox"].(map[string]interface{})["min"].([]interface{})[2].(float64))
		geometry.BB.Max.X = float32(embed["boundingBox"].(map[string]interface{})["max"].([]interface{})[0].(float64))
		geometry.BB.Max.Y = float32(embed["boundingBox"].(map[string]interface{})["max"].([]interface{})[1].(float64))
		geometry.BB.Max.Z = float32(embed["boundingBox"].(map[string]interface{})["max"].([]interface{})[2].(float64))

		geometry.Scale = float32(embed["scale"].(float64))
		vertices := embed["vertices"].([]interface{})
		for i := 0; i < len(vertices); {
			v := vertex.NewVertex(float32(vertices[i].(float64)), float32(vertices[i+1].(float64)), float32(vertices[i+2].(float64)))
			i += 3
			world.Vertices = append(world.Vertices, *v)
			geometry.Vertices = append(geometry.Vertices, &world.Vertices[len(world.Vertices)-1])
		}

		normals := embed["normals"].([]interface{})
		for i := 0; i < len(normals); {
			v := vertex.NewNormal(float32(normals[i].(float64)), float32(normals[i+1].(float64)), float32(normals[i+2].(float64)))
			i += 3
			world.Normals = append(world.Normals, *v)
			geometry.Normals = append(geometry.Normals, &world.Normals[len(world.Normals)-1])
		}

		faces := embed["faces"].([]interface{})
		for i := 0; i < len(faces); {
			mask := uint32(faces[i].(float64))
			i++

			// is not triangle?
			if isBitSet(mask, 0) {
				log.Fatalln("Quads are not supported")
			}
			tri := triangle.Triangle{Vertices: [3]*vertex.Vertex{&world.Vertices[int32(faces[i].(float64))], &world.Vertices[int32(faces[i+1].(float64))], &world.Vertices[int32(faces[i+2].(float64))]}}
			i += 3
			// has face material
			if isBitSet(mask, 1) {
				log.Println("Ignoring face material, material defined in the geometry json")
				i++
			}

			// has face uv
			if isBitSet(mask, 2) {
				log.Fatalln("Face UVs are not supported")
				i++
			}

			// has face vertex uv
			if isBitSet(mask, 3) {
				log.Fatalln("Face vertex UVs are not supported")
				i += 3
			}

			// has face normals
			if isBitSet(mask, 4) {
				log.Fatalln("Face normals are not supported")
				i++
			}

			// has face vertex normals
			if isBitSet(mask, 5) {
				tri.Normals = [3]*vertex.Normal{&world.Normals[int32(faces[i].(float64))], &world.Normals[int32(faces[i+1].(float64))], &world.Normals[int32(faces[i+2].(float64))]}
				i += 3
			}

			// has face color
			if isBitSet(mask, 6) {
				log.Fatalln("Face colors are not supported")
				i++
			}

			if isBitSet(mask, 7) {
				log.Fatalln("Face vertex colors are not supported")
				i += 3
			}

			world.Triangles = append(world.Triangles, tri)
			log.Println(tri)
		}
	} else {
		log.Panicln("I don't know what to do with non-embedded geometry types!")
	}

	return
}

func parseScene(world *scene.Scene, jsondata map[string]interface{}) (successful bool) {
	log.Println("Parsing Json Scene")
	parseMetadata(world, jsondata["metadata"].(map[string]interface{}))
	parseObjects(world, jsondata)
	return
}

func ParseJsonFile(world *scene.Scene, filepath string) (successful bool) {
	var jsondata map[string]interface{}
	rawdata, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(rawdata, &jsondata)
	parseScene(world, jsondata)

	successful = true
	return
}

func isBitSet(value, position uint32) bool {
	shifted := value & (1 << position)

	if shifted == 0 {
		return false
	} else {
		return true
	}
}
