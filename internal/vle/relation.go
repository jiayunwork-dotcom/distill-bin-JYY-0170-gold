package vle

import (
	"math"
)

type DewBubble struct {
	X bubblePoint
	Y dewPoint
}

type bubblePoint struct{ V float64 }
type dewPoint struct{ V float64 }

func BubblePressure(alpha, x float64) float64 {
	if alpha <= 0 || x <= 0 || x >= 1 {
		return 0
	}
	y := alpha * x / (1 + (alpha-1)*x)
	return y
}

func DewComposition(alpha, y float64) float64 {
	if alpha <= 1 || y <= 0 || y >= 1 {
		return 0
	}
	return y / (alpha - (alpha-1)*y)
}

func SeparationFactor(xD, xB float64) float64 {
	if xB <= 0 || xD >= 1 {
		return 0
	}
	return (xD / (1 - xD)) / (xB / (1 - xB))
}

func MinimumBoilup(nMin, rMin float64) float64 {
	return rMin + nMin
}

func VaporRate(D, R float64) float64 {
	return D * (R + 1)
}

func LiquidRate(D, R float64) float64 {
	return D * R
}

func BoilupRate(B, V, L float64) float64 {
	return V - L
}

func CheckRelativeVolatility(alpha float64) bool {
	return alpha > 1
}

func LiquidPhaseFraction(q float64) float64 {
	return q
}

func VaporPhaseFraction(q float64) float64 {
	return 1 - q
}

func EnrichmentRatio(alpha, x float64) float64 {
	if x <= 0 || x >= 1 {
		return 0
	}
	y := alpha * x / (1 + (alpha-1)*x)
	return y / x
}

func LogMeanSeparation(alpha, n float64) float64 {
	if n <= 0 {
		return 0
	}
	return math.Pow(alpha, n)
}

func StrippingFactor(alpha, x float64) float64 {
	return alpha * x / (1 + (alpha-1)*x)
}
