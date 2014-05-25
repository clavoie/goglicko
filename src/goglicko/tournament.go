package goglicko

import (
	"fmt"
	"math"
)

// A Tournament represents a player and a series of matches played by
// that player.
type Tournament struct {
	// The player for whom we're calculating the ratings (as Glicko2 rating)
	player *Rating

	// The ratings of the players played, as Glicko2 ratings. len(opponents) must
	// equal len(results).
	opponents []*Rating

	// Results of the matches. Should be one of Win, Loss, or Draw
	results []Result

	// Various Caches to make calculations faster.
	geeCache map[float64]float64

	etaCache map[int]float64

	convCache map[float64]float64
}

func NewTournament(p *Rating, o []*Rating, r []Result) (*Tournament, error) {
	// Step 1, Initialization, should have already taken place.
	if len(o) != len(r) {
		return nil, fmt.Errorf("Len(Opponents) must == Len(Results). Was %v %v",
			len(o), len(r))
	}

	// Step 2: Convert to Glicko2 ratings
	p2 := p.ToGlicko2()

	var o2 = make([]*Rating, len(o), len(o))
	for i := 0; i < len(o); i++ {
		o2[i] = o[i].ToGlicko2()
	}

	return &Tournament{
		p2,
		o2,
		r,
		make(map[float64]float64),
		make(map[int]float64),
		make(map[float64]float64),
	}, nil
}

// Calculate a new Glicko2 Rating based on the tournament results.
//
// Note: I expect the process to be completely inscrutable, but it should be
// easy(ish) to understand if you have the paper available.
func (t *Tournament) CalcRating() *Rating {
	estvar := t.estVariance()          // v
	estimp := t.estImprovement(estvar) // delta
	newvol := t.newVolatility(estvar, estimp)
	fmt.Printf("%v\n", newvol)
	return nil
}

// Step 3. Calculate the Estimated Variance based only on game outcomes.
func (t *Tournament) estVariance() float64 {
	sum := 0.0
	for i := 0; i < len(t.results); i++ {
		o := t.opponents[i]
		sum += t.gCached(o.Deviation) * t.gCached(o.Deviation) * t.eCached(i) * (1 - t.eCached(i))
	}
	sum = 1 / sum
	return sum
}

// Step 4. Calculate the estimated improvement (Delta), based only on game
// outcomes.
func (t *Tournament) estImprovement(estVar float64) float64 {
	sum := 0.0
	for i := 0; i < len(t.results); i++ {
		o := t.opponents[i]
		sum += t.gCached(o.Deviation) * (float64(t.results[i]) - t.eCached(i))
	}
	return sum * estVar
}

// Step 5. Calculate the new volatility
func (t *Tournament) newVolatility(estvar, estimp float64) float64 {
	epsilon := 0.000001
	a := math.Log(t.player.Volatility)
	deltaSq := sq(estimp)
	phiSq := sq(t.player.Deviation)
	tauSq := sq(DefaultTau)

	f := func(x float64) float64 {
		return conv(x, a, estvar, deltaSq, phiSq, tauSq)
	}

	A := a
	B := 0.0
	if deltaSq > (phiSq + estvar) {
		B = math.Log(deltaSq - phiSq - estvar)
	} else {
		val := -1.0
		k := 1
		for ; val < 0; k++ {
			val = f(a - float64(k)*DefaultTau)
		}
		B = a - float64(k)*DefaultTau
	}
	// Now: A < ln(sigma'^2) < B

	fA := f(A)
	fB := f(B)
	fC := 0.0
	for math.Abs(B-A) > epsilon {
		C := A + (A-B)*fA/(fB-fA)
		fC = f(C)
		if fC*fB < 0 {
			A = B
			fA = fB
		} else {
			fA = fA / 2
		}
		B = C
		fB = fC
	}

	newVol := math.Exp(A / 2)
	return newVol
}

//////////////////////
// Helper Functions //
//////////////////////

func (t *Tournament) eCached(i int) float64 {
	if val, ok := t.etaCache[i]; ok {
		return val
	}
	var o = t.opponents[i]
	val := ee(t.player.Rating, o.Rating, o.Deviation)
	t.etaCache[i] = val
	return val
}

func (t *Tournament) gCached(phi float64) float64 {
	if val, ok := t.geeCache[phi]; ok {
		return val
	}
	val := gee(phi)
	t.geeCache[phi] = val
	return val
}
