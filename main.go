package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/strangedev/vroom/algebra"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
	"fmt"
)

const (
	width  = 1024
	height = 768
)

func RandFloat(f float64) float64 {
	return (0.2 * f * rand.Float64())
}

func run() {
	world := NewWorld(width, height)

	for i := 0; i < 100; i++ {
		c := NewCarWithGfx()
		c.Steering.Acceleration = RandFloat(0.6)
		c.Steering.SteerRadians = (rand.Float64() - 0.5) * RandFloat(0.0001)
		c.Position = &algebra.Vector2{RandFloat(width), RandFloat(height)}
		fmt.Printf("Placing @ %v\n", c.Position)
		world.AddObject(&c)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "Vroom!",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Skyblue)

	last := time.Now()
	for !win.Closed() {
		win.Clear(colornames.Skyblue)
		dt := time.Since(last).Seconds()
		last = time.Now()

		world.Tick(dt)
		world.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
