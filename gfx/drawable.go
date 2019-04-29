package gfx

import "github.com/faiface/pixel/pixelgl"

type Drawable interface {
	Draw(win *pixelgl.Window)
}
