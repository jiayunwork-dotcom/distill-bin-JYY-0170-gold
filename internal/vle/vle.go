package vle

import (
	"errors"
	"math"
)

var (
	ErrAlphaNonSeparating = errors.New("vle: relative volatility at most 1, no separation")
	ErrCompositionOutOfRange = errors.New("vle: composition must be in (0,1)")
	ErrInvalidFeed = errors.New("vle: feed composition must lie between bottoms and distillate")
)

type VLE struct {
	Alpha  float64
	Points []CurvePoint
}

type CurvePoint struct {
	X float64
	Y float64
}

func NewConstantAlpha(alpha float64) (VLE, error) {
	if alpha <= 1 {
		return VLE{}, ErrAlphaNonSeparating
	}
	return VLE{Alpha: alpha}, nil
}

func NewCurve(points []CurvePoint) (VLE, error) {
	for i, p := range points {
		isEndpoint := i == 0 && p.X == 0 && p.Y == 0
		isTopEndpoint := i == len(points)-1 && p.X == 1 && p.Y == 1
		if isEndpoint || isTopEndpoint {
			continue
		}
		if p.X <= 0 || p.X >= 1 || p.Y <= 0 || p.Y >= 1 {
			return VLE{}, ErrCompositionOutOfRange
		}
	}
	return VLE{Points: points}, nil
}

func (v VLE) Y(x float64) float64 {
	if len(v.Points) > 0 {
		return interpolate(v.Points, x)
	}
	return v.Alpha * x / (1 + (v.Alpha-1)*x)
}

func (v VLE) X(y float64) float64 {
	if len(v.Points) > 0 {
		return interpolateX(v.Points, y)
	}
	return y / (v.Alpha - (v.Alpha-1)*y)
}

func (v VLE) ValidComposition(c float64) bool {
	return c > 0 && c < 1
}

func interpolate(points []CurvePoint, x float64) float64 {
	if x <= points[0].X {
		return points[0].Y
	}
	for i := 1; i < len(points); i++ {
		if x <= points[i].X {
			p0 := points[i-1]
			p1 := points[i]
			if p1.X == p0.X {
				return p1.Y
			}
			t := (x - p0.X) / (p1.X - p0.X)
			return p0.Y + t*(p1.Y-p0.Y)
		}
	}
	return points[len(points)-1].Y
}

func interpolateX(points []CurvePoint, y float64) float64 {
	if y <= points[0].Y {
		return points[0].X
	}
	for i := 1; i < len(points); i++ {
		if y <= points[i].Y {
			p0 := points[i-1]
			p1 := points[i]
			if p1.Y == p0.Y {
				return p1.X
			}
			t := (y - p0.Y) / (p1.Y - p0.Y)
			return p0.X + t*(p1.X-p0.X)
		}
	}
	return points[len(points)-1].X
}

func RelativeVolatilityFromPoints(points []CurvePoint) float64 {
	sum := 0.0
	n := 0
	for _, p := range points {
		if p.X <= 0 || p.X >= 1 {
			continue
		}
		alpha := p.Y * (1 - p.X) / (p.X * (1 - p.Y))
		if alpha > 1 && !math.IsInf(alpha, 0) && !math.IsNaN(alpha) {
			sum += alpha
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}
