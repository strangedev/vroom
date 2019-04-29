package algebra

import "math"

type LineSegment2 struct {
	Start *Vector2
	End   *Vector2
}

func (l1 LineSegment2) IntersectionWith(l2 LineSegment2) (v Vector2, e error) {
	epsilon := math.Nextafter(1, 2) - 1
	ab := l1.End.Sub(*l1.Start)
	cd := l2.End.Sub(*l2.Start)
	ca := l2.Start.Sub(*l1.Start)

	if ab.Length() < epsilon {
		return v, &undefinedError{"First line is a point."}
	} else if cd.Length() < epsilon {
		return v, &undefinedError{"Second line is a point."}
	}

	abn := ab.Normalize()
	cdn := cd.Normalize()
	cdn.ScaleInPlace(-1)

	coeffs := Matrix2FromColumns(abn, cdn)
	st := coeffs.Solve(ca)

	abn.ScaleInPlace(st[0])
	v = l1.Start.Add(abn)

	if math.IsNaN(v[0]) || math.IsNaN(v[1]) {
		return v, &undefinedError{"Lines are parallel"}
	}

	// check if v is on ab
	av := v.Sub(*l1.Start)
	abDotAv := ab.Dot(av)
	abDotAb := ab.Dot(ab)
	if abDotAv < 0 || abDotAv > abDotAb {
		return v, &undefinedError{"No intersection within the first segment."}
	}

	// check if v is on cd
	cv := v.Sub(*l2.Start)
	cdDotCv := cd.Dot(cv)
	cdDotCd := cd.Dot(cd)
	if cdDotCv < 0 || cdDotCv > cdDotCd {
		return v, &undefinedError{"No intersection within the second segment."}
	}

	return v, nil
}
