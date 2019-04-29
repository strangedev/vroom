package algebra

type Rectangle struct {
	Ul, Ur, Dl, Dr Vector2
}

func (r Rectangle) Covers(v Vector2) bool {
	return v[0] >= r.Ul[0] && v[0] <= r.Ur[0] && v[1] >= r.Dl[1] && v[1] <= r.Ul[1]
}

func (r Rectangle) Clips(other Rectangle) bool {
	return r.Covers(other.Ul) || r.Covers(other.Ur) || r.Covers(other.Dl) || r.Covers(other.Dr)
}

func (r Rectangle) Size() (width, height float64) {
	width, height = r.Ur[0]-r.Ul[0], r.Ul[1]-r.Dr[1]
	return
}

func (r Rectangle) Translate(v Vector2) Rectangle {
	return Rectangle{
		r.Ul.Add(v), r.Ur.Add(v), r.Dl.Add(v), r.Dr.Add(v),
	}
}

func (r *Rectangle) TranslateInPlace(v Vector2) {
	r.Ul.AddInPlace(v)
	r.Ur.AddInPlace(v)
	r.Dl.AddInPlace(v)
	r.Dl.AddInPlace(v)
}
