package main

import (
	"github.com/faiface/pixel/pixelgl"
	//"github.com/faiface/pixel"
	"github.com/strangedev/vroom/gfx"
	"github.com/strangedev/vroom/physics"
	"github.com/faiface/pixel/imdraw"
)

type Object interface {
	physics.Physical
	gfx.Drawable
	GetId() string
}

type World struct {
	objects       []Object
	width, height float64
	quadTree physics.QuadTree
	imd *imdraw.IMDraw
}

func NewWorld(width, height float64) World {
	return World{
		make([]Object, 0, 100),
		width, height,
		physics.NewQuadTree(width, height),
		imdraw.New(nil),
	}
}

func (w *World) AddObject(o Object) {
	w.objects = append(w.objects, o)
}

func (w *World) Tick(dt float64) {
	w.quadTree = physics.NewQuadTree(width, height)
	for _, o := range w.objects {
		o.Tick(dt)
		if l, ok := o.(physics.Positionable); ok {
			pos := l.GetPosition()
			if pos[0] < 0 {
				pos[0] = w.width
			}
			if pos[0] > w.width {
				pos[0] = 0
			}
			if pos[1] < 0 {
				pos[1] = w.height
			}
			if pos[1] > w.height {
				pos[1] = 0
			}
			l.MoveTo(pos)
		}
		if c, ok := o.(physics.ClippingShape); ok {
			w.quadTree.InsertAt(&c, c.GetBoundingBox())
		}
	}
	//w.imd.Clear()
	//w.imd.Color = pixel.RGB(1, 0, 0)
	/*
	for _, o := range w.objects {
		if c, ok := o.(physics.ClippingShape); ok {
			for s := range w.quadTree.ClippingCandidatesAt(c.GetBoundingBox()) {
				if i, ok := s.Pnt.(Object); ok {
					if i.GetId() == o.GetId() {
						continue
					}
				}
				w.imd.Push(
					s.Bounds.Ul.ToPixelVec(),
					s.Bounds.Ur.ToPixelVec(),
				)
				w.imd.Line(5)
				w.imd.Push(
					s.Bounds.Ul.ToPixelVec(),
					s.Bounds.Dl.ToPixelVec(),
				)
				w.imd.Line(5)
				w.imd.Push(
					s.Bounds.Ur.ToPixelVec(),
					s.Bounds.Dr.ToPixelVec(),
				)
				w.imd.Line(5)
				w.imd.Push(
					s.Bounds.Dl.ToPixelVec(),
					s.Bounds.Dr.ToPixelVec(),
				)
				w.imd.Line(5)
			}
		}
		//break
	}*/
}

func (w *World) Draw(win *pixelgl.Window) {
	w.quadTree.Draw(win)
	w.imd.Draw(win)
	for _, o := range w.objects {
		o.Draw(win)
	}
}
