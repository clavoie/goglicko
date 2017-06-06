package goglicko

// System represents the Glicko defaults used to create the rating
type System struct {
	baseRating     float64
	baseDeviation  float64
	baseVolatility float64
	tau            float64 // constrains system volatility, should be between 0.3-1.2
}

// NewDefaultSystem creates a new System using DefaultRat, DefaultDev, DefaultVol, and DefaultTau
func NewDefaultSystem() *System {
	return NewSystem(DefaultRat, DefaultDev, DefaultVol, DefaultTau)
}

// NewSystem creates a custom System
func NewSystem(baseRating, baseDeviation, baseVolitility, tau float64) *System {
	return &System{baseRating, baseDeviation, baseVolitility, tau}
}

// GetValues returns the base rating, deviation, volatility, and tau
func (s *System) GetValues() (float64, float64, float64, float64) {
	return s.baseRating, s.baseDeviation, s.baseVolatility, s.tau
}
