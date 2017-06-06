package goglicko

import (
	"fmt"
)

// Represents a player's rating and the confidence in a player's rating.
type Rating struct {
	rating     float64 // Player's rating. Usually starts off at 1500.
	deviation  float64 // Confidence/uncertainty in a player's rating
	volatility float64 // Measures erratic performances
	system     *System // the values from which the rating was created
}

// Creates a default Rating using:
// 	Rating     = DefaultRat
// 	Deviation  = DefaultDev
// 	Volatility = DefaultVol
func NewDefaultRating() *Rating {
	return &Rating{DefaultRat, DefaultDev, DefaultVol, NewDefaultSystem()}
}

// Creates a new custom Rating.
func NewRating(r, rd, s float64, sys *System) *Rating {
	return &Rating{r, rd, s, sys}
}

// Creates a new rating, converted from Glicko1 scaling to Glicko2 scaling.
// This assumes the starting rating value is 1500.
func (r *Rating) toGlicko2() *Rating {
	return NewRating(
		(r.rating-r.system.baseRating)/glicko2Scale,
		(r.deviation)/glicko2Scale,
		r.volatility, r.system)
}

// Creates a new rating, converted from Glicko2 scaling to Glicko1 scaling.
// This assumes the starting rating value is 1500.
func (r *Rating) fromGlicko2() *Rating {
	return NewRating(
		r.rating*glicko2Scale+r.system.baseRating,
		r.deviation*glicko2Scale,
		r.volatility, r.system)
}

func (r *Rating) String() string {
	return fmt.Sprintf("{Rating[%.3f] Deviation[%.3f] Volatility[%.3f]}",
		r.rating, r.deviation, r.volatility)
}

// Create a duplicate rating with the same values.
func (r *Rating) Copy() *Rating {
	rCopy := *r

	return &rCopy
}

// GetSystem returns the system that created this Rating
func (r *Rating) GetSystem() *System {
	return r.system
}

// GetValues returns the rating, deviation, and volatility
func (r *Rating) GetValues() (float64, float64, float64) {
	return r.rating, r.deviation, r.volatility
}
