/*
 * Author: Matias Piispanen (@gmail.com)
 *
 * The purpose of this project is to learn about ray tracing techniques and the
 * Go language simultaneously. Eventually I'll try to add OpenCL acceleration to
 * the ray tracer.
 */

package main

import (
	"code.google.com/p/rog-go/canvas"
	"code.google.com/p/x-go-binding/ui"
	"code.google.com/p/x-go-binding/ui/x11"
	"fmt"
	"gotracer/base/renderer"
	"gotracer/base/scene"
	"gotracer/base/sceneparser"
	"image"
	"image/color"
	"log"
	"os"
	"time"
)

var (
	defaultScenePath = "/Users/mepiispa/go/src/gotracer/scenes/polygon.json"
	cvs              *canvas.Canvas
	ticker           *time.Ticker
	world            *scene.Scene
	width, height    uint32
	output           *image.RGBA
)

func init() {
	fmt.Println("Initialising program")
	world = scene.NewScene()
}

func clear() {
	bbox := cvs.Bbox()
	output = image.NewRGBA(bbox)
	for i := 0; i < bbox.Dx(); i++ {
		for j := 0; j < bbox.Dy(); j++ {
			output.Set(i, j, color.RGBA{0, 0, 0, 0xff})
		}
	}
	cvs.AddItem(canvas.NewImage(output, true, image.Point{X: 0, Y: 0}))
}

func paint() {
	world.Lock()
	cvs.AddItem(canvas.NewImage(output, true, image.Point{X: 0, Y: 0}))
	cvs.Flush()
	world.Unlock()
}

func readKeyboardInput(ec <-chan interface{}) {
	for {
		select {
		case e, ok := <-ec:
			if !ok {
				fmt.Println("Something went wrong with event handling")
				return
			}

			switch e := e.(type) {
			case ui.KeyEvent:
				switch e.Key {
				case 'q':
					fmt.Println("Quitting")
					ticker.Stop()
					os.Exit(0)
				}
			}
		}
	}
}

func main() {
	// Ticker that keeps the world ticking
	ticker = time.NewTicker(time.Nanosecond * 40000)
	// Initialize screen
	win, err := x11.NewWindow()
	if win == nil {
		log.Fatalf("no window: %v", err)
	}
	screen := win.Screen()

	bg := canvas.NewBackground(screen.(*image.RGBA), image.White, nil)
	cvs = canvas.NewCanvas(nil, bg.Rect())
	bg.SetItem(cvs)
	bbox := cvs.Bbox()
	height = uint32(bbox.Dy())
	width = uint32(bbox.Dx())

	// Set even handler for ui events
	ec := win.EventChan()
	go readKeyboardInput(ec)

	var last time.Time
	var now time.Time
	last = time.Now()

	// Parse scene
	sceneparser.ParseJsonFile(world, defaultScenePath)

	// Clear output buffer
	clear()

	renderer.Render(world, width, height, output)
	paint()
	win.FlushImage()

	// Main loop
	for {
		now = <-ticker.C
		paint()
		win.FlushImage()

		dt := now.Nanosecond() - last.Nanosecond()
		last = now
		dt = dt
	}
}
