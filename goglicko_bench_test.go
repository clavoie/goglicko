package goglicko

import "testing"

func BenchmarkSimpleExample(b *testing.B) {
	p := NewDefaultRating()
	o := []*Rating{
		NewRating(1400, 30, DefaultVol, NewDefaultSystem()),
		NewRating(1550, 100, DefaultVol, NewDefaultSystem()),
		NewRating(1700, 300, DefaultVol, NewDefaultSystem()),
	}
	res := []Result{1, 0, 0}

	for i := 0; i < b.N; i++ {
		p.Update(o, res)
	}
}
