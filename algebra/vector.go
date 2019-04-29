package algebra

import (
	"github.com/faiface/pixel"
	"math"
)

type Vector2 [2]float64

func (v Vector2) Length() float64 {
	return math.Sqrt(
		math.Pow(v[0], 2) +
			math.Pow(v[1], 2),
	)
}

func (v1 Vector2) Dot(v2 Vector2) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1]
}

func (v Vector2) Normalize() (r Vector2) {
	l := v.Length()
	if l == 0 {
		return v
	}
	return v.Scale(1 / l)
}

func (v *Vector2) NormalizeInPlace() {
	l := v.Length()
	if l == 0 {
		return
	}
	v.ScaleInPlace(1 / l)
}

func (v Vector2) Add(w Vector2) (r Vector2) {
	return Vector2{
		v[0] + w[0],
		v[1] + w[1],
	}
}

func (v Vector2) Sub(w Vector2) (r Vector2) {
	return Vector2{
		v[0] - w[0],
		v[1] - w[1],
	}
}

func (v Vector2) Scale(s float64) (r Vector2) {
	r[0] = s * v[0]
	r[1] = s * v[1]
	return
}

func (v *Vector2) AddInPlace(w Vector2) {
	v[0] += w[0]
	v[1] += w[1]
}

func (v *Vector2) SubInPlace(w Vector2) {
	v[0] -= w[0]
	v[1] -= w[1]
}

func (v *Vector2) ScaleInPlace(s float64) {
	v[0] = s * v[0]
	v[1] = s * v[1]
}

func (v Vector2) Rotate(rad float64) (r Vector2) {
	m := Rotation2(RadToDeg(rad))
	return m.MulV(v)
}

func (v *Vector2) RotateInPlace(rad float64) {
	m := Rotation2(RadToDeg(rad))
	m.MulVInPlace(v)
}

func (v Vector2) AngleTo(w Vector2) (rad float64) {
	m := Matrix2FromColumns(v, w)
	rad = math.Atan2(m.Det(), v.Dot(w))
	return
}

func (v Vector2) ToPixelVec() pixel.Vec {
	return pixel.V(v[0], v[1])
}
