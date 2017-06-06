package goglicko

import "testing"

// Ensure that some other Rating is equal to this rating, given some epsilon. In
// other words, find the error between this rating's values and the other
// rating's values and make sure it's less than epsilon in absolute value.
func (r *Rating) MostlyEquals(o *Rating, epsilon float64) bool {
	return floatsMostlyEqual(r.rating, o.rating, epsilon) &&
		floatsMostlyEqual(r.deviation, o.deviation, epsilon) &&
		floatsMostlyEqual(r.volatility, o.volatility, epsilon)
}

func TestGlicko(t *testing.T) {
	// Much of this data comes from the paper:
	// http://en.wikipedia.org/wiki/Glicko_rating_system
	var pl *Rating
	var opps []*Rating
	var results []Result

	reset := func() {
		pl = NewRating(1500, 200, DefaultVol, NewDefaultSystem())
		opps = []*Rating{
			NewRating(1400, 30, DefaultVol, NewDefaultSystem()),
			NewRating(1550, 100, DefaultVol, NewDefaultSystem()),
			NewRating(1700, 300, DefaultVol, NewDefaultSystem()),
		}
		results = []Result{1, 0, 0}
	}

	t.Run("TestEquivTransfOpps", func(t *testing.T) {
		reset()

		for i := range opps {
			o := opps[i]
			o2 := opps[i].toGlicko2().fromGlicko2()
			if !o.MostlyEquals(o2, 0.0001) {
				t.Errorf("o %v != o2 %v", o, o2)
			}
		}
	})

	t.Run("TestToGlicko2", func(t *testing.T) {
		reset()

		p2 := pl.toGlicko2()
		exp := NewRating(0, 1.1513, DefaultVol, NewDefaultSystem())
		if !p2.MostlyEquals(exp, 0.0001) {
			t.Errorf("p2 %v != expected %v", p2, exp)
		}
	})

	t.Run("TestOppToGlicko2", func(t *testing.T) {
		reset()

		exp := []*Rating{
			NewRating(-0.5756, 0.1727, DefaultVol, NewDefaultSystem()),
			NewRating(0.2878, 0.5756, DefaultVol, NewDefaultSystem()),
			NewRating(1.1513, 1.7269, DefaultVol, NewDefaultSystem()),
		}
		for i := range exp {
			g2 := opps[i].toGlicko2()
			if !g2.MostlyEquals(exp[i], 0.0001) {
				t.Errorf("For i=%v: Glicko2 scaled opp %v != expected %v\n", i, g2, exp[i])
			}
		}
	})

	t.Run("TestEeGeeValues", func(t *testing.T) {
		reset()

		expGee := []float64{0.9955, 0.9531, 0.7242}
		expEe := []float64{0.639, 0.432, 0.303}
		p2 := pl.toGlicko2()
		for i := range opps {
			o := opps[i].toGlicko2()
			geeVal := gee(o.deviation)
			if !floatsMostlyEqual(geeVal, expGee[i], 0.0001) {
				t.Errorf("Floats not mostly equal. g=%v exp_g=%v", geeVal, expGee[i])
			}
			eeVal := ee(p2.rating, o.rating, o.deviation)
			if !floatsMostlyEqual(eeVal, expEe[i], 0.001) {
				t.Errorf("Floats not mostly equal. ee=%v exp_ee=%v", eeVal, expEe[i])
			}
		}
	})

	t.Run("TestAlgorithm", func(t *testing.T) {
		reset()

		p2 := pl.toGlicko2()
		gees := make([]float64, len(opps))
		ees := make([]float64, len(opps))
		for i := range opps {
			o := opps[i].toGlicko2()
			gees[i] = gee(o.deviation)
			ees[i] = ee(p2.rating, o.rating, o.deviation)
		}
		estVar := estVariance(gees, ees)
		exp := 1.7785
		if !floatsMostlyEqual(estVar, exp, 0.001) {
			t.Errorf("estvar %v != exp %v", estVar, exp)
		}

		// Test Delta
		estImpPart := estImprovePartial(gees, ees, results)
		estImp := estVar * estImpPart
		expEstImp := -0.4834
		if !floatsMostlyEqual(estImp, expEstImp, 0.001) {
			t.Errorf("delta %v != exp %v", estImp, expEstImp)
		}

		// Test calculating the new volatility
		newVol := p2.newVolatility(estVar, estImp)
		expNewVol := 0.05999
		if !floatsMostlyEqual(newVol, expNewVol, 0.0001) {
			t.Errorf("newVol %v != expNewVol %v", newVol, expNewVol)
		}

		newDev := newDeviation(p2.deviation, newVol, estVar)
		expNewDev := 0.8722
		if !floatsMostlyEqual(newDev, expNewDev, 0.0001) {
			t.Errorf("newDev %v != expNewDev %v", newDev, expNewDev)
		}

		newRating := newRatingVal(p2.rating, newDev, estImpPart)
		expNewRating := -0.2069
		if !floatsMostlyEqual(newRating, expNewRating, 0.0001) {
			t.Errorf("newRating %v != expNewRating %v", newRating, expNewRating)
		}

		newPlayer := NewRating(newRating, newDev, newVol, NewDefaultSystem()).fromGlicko2()
		expNewRatingV1 := 1464.06
		if !floatsMostlyEqual(newPlayer.rating, expNewRatingV1, 0.01) {
			t.Errorf("newPlayer.Rating %v != expNewRatingV1 %v",
				newPlayer.rating, expNewRatingV1)
		}
		expNewDevV1 := 151.52
		if !floatsMostlyEqual(newPlayer.deviation, expNewDevV1, 0.01) {
			t.Errorf("newPlayer.Deviation %v != expNewDevV1 %v",
				newPlayer.deviation, expNewDevV1)
		}
	})

	t.Run("TestUpdate", func(t *testing.T) {
		reset()

		err := pl.Update(opps, results)
		if err != nil {
			t.Errorf("Error while calculating results: %v", err)
			return
		}

		expNewVol := 0.05999
		if !floatsMostlyEqual(pl.volatility, expNewVol, 0.0001) {
			t.Errorf("pl.Volatility %v != expNewVol %v", pl.volatility, expNewVol)
		}
		expNewRatingV1 := 1464.06
		if !floatsMostlyEqual(pl.rating, expNewRatingV1, 0.01) {
			t.Errorf("pl.Rating %v != expNewRatingV1 %v",
				pl.rating, expNewRatingV1)
		}
		expNewDevV1 := 151.52
		if !floatsMostlyEqual(pl.deviation, expNewDevV1, 0.01) {
			t.Errorf("pl.Deviation %v != expNewDevV1 %v",
				pl.deviation, expNewDevV1)
		}
	})
}
