package sceneparser

import (
	"encoding/json"
	"gotracer/base/scene"
	"io/ioutil"
	"log"
)

var (
	decoder *json.Decoder
)

func parseMetadata(world *scene.Scene, metadata map[string]interface{}) {
	world.Info.FormatVersion = float32(metadata["formatVersion"].(float64))
	world.Info.Objects = int32(metadata["objects"].(float64))
	world.Info.Geometries = int32(metadata["geometries"].(float64))
	world.Info.Materials = int32(metadata["materials"].(float64))
	world.Info.Textures = int32(metadata["textures"].(float64))
}

func parseScene(world *scene.Scene, jsondata map[string]interface{}) (successful bool) {
	log.Println("Parsing Json Object")
	parseMetadata(world, jsondata["metadata"].(map[string]interface{}))
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
