package algebra

import (
	"fmt"
	"math"
)

const Epsilon = 1e-9

/*
	matrix of elements a_ji:
	|a11 a12|
	|a21 a22|
*/
type Matrix2 [2][2]float64

func Matrix2FromColumns(c1, c2 Vector2) (m Matrix2) {
	m[0][0] = c1[0]
	m[1][0] = c1[1]
	m[0][1] = c2[0]
	m[1][1] = c2[1]
	return
}

func Identity2() Matrix2 {
	return Matrix2{[2]float64{1, 0}, [2]float64{0, 1}}
}

func Rotation2(angle float64) Matrix2 {
	return Matrix2{
		[2]float64{math.Cos(angle), -math.Sin(angle)},
		[2]float64{math.Sin(angle), math.Cos(angle)},
	}
}

func (m Matrix2) Print() {
	fmt.Printf("|%e\t%e|\n", m[0][0], m[0][1])
	fmt.Printf("|%e\t%e|\n", m[1][0], m[1][1])
}

func (m Matrix2) Pivot() (p, a Matrix2) {
	p = Identity2()
	if math.Abs(m[0][0]) < math.Abs(m[1][0]) {
		p[0], p[1] = p[1], p[0]
	}
	a = p.Mul(m)
	return
}

func (m Matrix2) Mul(k Matrix2) (r Matrix2) {
	return Matrix2{
		[2]float64{
			m[0][0]*k[0][0] + m[0][1]*k[1][0],
			m[0][0]*k[0][1] + m[0][1]*k[1][1],
		},
		[2]float64{
			m[1][0]*k[0][0] + m[1][1]*k[1][0],
			m[1][0]*k[0][1] + m[1][1]*k[1][1],
		},
	}
}

func (m Matrix2) MulV(v Vector2) (w Vector2) {
	return Vector2{
		m[0][0]*v[0] + m[0][1]*v[1],
		m[1][0]*v[0] + m[1][1]*v[1],
	}
}

func (m Matrix2) MulVInPlace(v *Vector2) {
	v[0], v[1] = m[0][0]*v[0]+m[0][1]*v[1], m[1][0]*v[0]+m[1][1]*v[1]
}

func (m Matrix2) LUDecomposition() (l, u, p Matrix2) {
	var a Matrix2
	p, a = m.Pivot()
	u = Matrix2{a[0]} // upper triangular
	l = Identity2()   // lower triangular

	l[1][0] = a[1][0] / u[0][0]
	u[1][1] = a[1][1] - u[0][1]*l[1][0]
	return
}

func (m Matrix2) Solve(v Vector2) (x Vector2) {
	l, u, p := m.LUDecomposition()
	pv := p.MulV(v)

	// forwards substitution
	y := Vector2{pv[0], pv[1] - l[1][0]*pv[0]}
	// backwards substitution
	x[1] = y[1] / u[1][1]
	x[0] = (y[0] - u[0][1]*x[1]) / u[0][0]
	return
}

func (m Matrix2) Det() float64 {
	return m[0][0]*m[1][1] - m[0][1]*m[1][0]
}
