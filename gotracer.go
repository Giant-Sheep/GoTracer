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
	"gotracer/base/scene"
	"gotracer/base/sceneparser"
	"image"
	"log"
	"os"
	"sync"
	"time"
)

var (
	defaultScenePath = "/Users/mepiispa/go/src/gotracer/scenes/polygon.json"
	cvs              *canvas.Canvas
	sceneMutex       sync.Mutex
	ticker           *time.Ticker
	world            *scene.Scene
)

func init() {
	fmt.Println("Initialising program")
	world = new(scene.Scene)
}

func draw() {
	sceneMutex.Lock()
	cvs.Flush()
	sceneMutex.Unlock()
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

	bg := canvas.NewBackground(screen.(*image.RGBA), image.Black, nil)
	cvs = canvas.NewCanvas(nil, bg.Rect())
	bg.SetItem(cvs)

	// Set even handler for ui events
	ec := win.EventChan()
	go readKeyboardInput(ec)

	var last time.Time
	var now time.Time
	last = time.Now()

	// Parse scene
	sceneparser.ParseJsonFile(world, defaultScenePath)

	// Main loop
	for {
		now = <-ticker.C
		draw()

		dt := now.Nanosecond() - last.Nanosecond()
		last = now
		dt = dt
	}
}
