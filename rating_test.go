package goglicko

import "testing"

func TestScaleRescale(t *testing.T) {
	def := NewDefaultRating()
	p2 := def.toGlicko2().fromGlicko2()
	if !def.MostlyEquals(p2, 0.0001) {
		t.Errorf("Test Failed. def %v != p2 %v", def, p2)
	}
}

func TestStringyf(t *testing.T) {
	def := NewDefaultRating()
	if def.String() != "{Rating[1500.000] Deviation[350.000] Volatility[0.060]}" {
		t.Errorf("Error. String form was %v", def.String())
	}
}
