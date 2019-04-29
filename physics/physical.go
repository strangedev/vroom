package physics

import "github.com/strangedev/vroom/algebra"

type Physical interface {
	Tick(dt float64)
}

type Locatable interface {
	GetPosition() algebra.Vector2
}

type Positionable interface {
	Locatable
	MoveTo(v algebra.Vector2)
}

type ClippingShape interface {
	GetBoundingBox() algebra.Rectangle
}

type SolidBody interface {
	ApplyImpulse(p algebra.Vector2)
}
