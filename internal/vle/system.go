package vle

import "math"

type System struct {
	Alpha float64
	PB    float64
	PA    float64
}

func NewSystem(alpha, pb float64) System {
	return System{Alpha: alpha, PB: pb, PA: alpha * pb}
}

func (s System) PartialPressure(x float64) (float64, float64) {
	// Raoult's law for a binary ideal system
	pA := x * s.PA
	pB := (1 - x) * s.PB
	return pA, pB
}

func (s System) TotalPressure(x float64) float64 {
	pA, pB := s.PartialPressure(x)
	return pA + pB
}

func (s System) BubbleY(x float64) float64 {
	return s.Alpha * x / (1 + (s.Alpha-1)*x)
}

func (s System) DewX(y float64) float64 {
	return y / (s.Alpha - (s.Alpha-1)*y)
}

func (s System) TotalPressureAtBubble(x float64) float64 {
	return s.TotalPressure(x)
}

func (s System) CheckThermodynamicConsistency(x, y float64) bool {
	if x <= 0 || x >= 1 || y <= 0 || y >= 1 {
		return false
	}
	calcY := s.BubbleY(x)
	return math.Abs(calcY-y) < 1e-6
}

func (s System) BubbleTemperature(psatA, psatB float64) float64 {
	// constant alpha implies constant ratio; temperature implied by ratio
	if psatB <= 0 {
		return 0
	}
	return psatA / psatB
}

func (s System) KValue(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return s.Alpha / (1 + (s.Alpha-1)*x)
}

func (s System) RelativeVolatilityAt(x float64) float64 {
	if x <= 0 || x >= 1 {
		return 0
	}
	y := s.BubbleY(x)
	return y * (1 - x) / (x * (1 - y))
}

func (s System) VolatilityRatio(x1, x2 float64) float64 {
	a1 := s.RelativeVolatilityAt(x1)
	a2 := s.RelativeVolatilityAt(x2)
	if a2 <= 0 {
		return 0
	}
	return a1 / a2
}

func (s System) Constant() bool {
	return math.Abs(s.Alpha-s.RelativeVolatilityAt(0.5)) < 1e-9
}
