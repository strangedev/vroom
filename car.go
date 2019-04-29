package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/strangedev/vroom/algebra"
	"github.com/strangedev/vroom/gfx"
	"github.com/faiface/pixel/imdraw"
	"github.com/rs/xid"
)

type SteeringIntent struct {
	SteerRadians float64
	Acceleration float64
}

type Car struct {
	Position    *algebra.Vector2
	Orientation *algebra.Vector2
	Velocity    float64
	Steering    *SteeringIntent
	tickCount   uint64
	id    		string
	Sprite      *pixel.Sprite
}

func NewCar() Car {
	id := xid.New()
	return Car{
		Position:    &algebra.Vector2{},
		Orientation: &algebra.Vector2{1, 0},
		Velocity:    0,
		Steering:    &SteeringIntent{},
		id:			 id.String(),
	}
}

func NewCarWithGfx() (c Car) {
	c = NewCar()
	pic, err := gfx.LoadPicture("assets/sprites/turtle-small.png")
	if err != nil {
		panic(err)
	}
	c.Sprite = pixel.NewSprite(pic, pic.Bounds())
	return
}

func (c *Car) GetTransformationMatrix() (m pixel.Matrix) {
	carAngle := -c.Orientation.AngleTo(algebra.Vector2{1, 0})
	carPosition := c.Position.ToPixelVec()
	m = pixel.IM
	m = m.Scaled(pixel.ZV, 0.5)
	m = m.Rotated(pixel.ZV, carAngle)
	m = m.Moved(carPosition)
	return
}

func (c *Car) Draw(win *pixelgl.Window) {
	bounds := c.GetBoundingBox()
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(0, 1, 0)
	imd.Push(
		bounds.Ul.ToPixelVec(),
		bounds.Ur.ToPixelVec(),
	)
	imd.Line(2)
	imd.Push(
		bounds.Ul.ToPixelVec(),
		bounds.Dl.ToPixelVec(),
	)
	imd.Line(2)
	imd.Push(
		bounds.Ur.ToPixelVec(),
		bounds.Dr.ToPixelVec(),
	)
	imd.Line(2)
	imd.Push(
		bounds.Dl.ToPixelVec(),
		bounds.Dr.ToPixelVec(),
	)
	imd.Line(2)
	imd.Draw(win)

	c.Sprite.Draw(win, c.GetTransformationMatrix())
}

func (c *Car) Tick(dt float64) {
	c.Velocity += c.Steering.Acceleration * dt
	c.Position.AddInPlace(c.Orientation.Scale(c.Velocity))
	c.Orientation.RotateInPlace(c.Steering.SteerRadians)
	c.tickCount++
}

func (c Car) Print() {
	fmt.Printf("Car: %p (%v ticks)\n", &c, c.tickCount)
	fmt.Printf("  Pos:\t%v\n", *c.Position)
	fmt.Printf("  Ori:\t%v\n", *c.Orientation)
	fmt.Printf("  Vel:\t%v\n", c.Velocity)
	fmt.Println("  Steer:")
	fmt.Printf("    Acc:\t%v\n", c.Steering.Acceleration)
	fmt.Printf("    Rad:\t%v\n", c.Steering.SteerRadians)
	fmt.Println("---")
}

func (c *Car) GetBoundingBox() algebra.Rectangle {
	return algebra.Rectangle{
		Ul: algebra.Vector2{c.Position[0] - 5, c.Position[1] + 5},
		Ur: algebra.Vector2{c.Position[0] + 5, c.Position[1] + 5},
		Dl: algebra.Vector2{c.Position[0] - 5, c.Position[1] - 5},
		Dr: algebra.Vector2{c.Position[0] + 5, c.Position[1] - 5},
	}
}

func (c *Car) GetPosition() algebra.Vector2 {
	return *(c.Position)
}

func (c *Car) MoveTo(v algebra.Vector2) {
	c.Position = &v
}

func (c *Car) GetId() string {
	return c.id
}