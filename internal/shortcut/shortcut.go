package shortcut

import (
	"errors"
	"math"

	"distill-bin/internal/vle"
)

var (
	ErrNonSeparating = errors.New("shortcut: alpha at most 1, no separation")
	ErrInvalidComposition = errors.New("shortcut: compositions must satisfy xD > zF > xB")
	ErrNoRealRoot = errors.New("shortcut: Underwood root search failed")
)

type Shortcut struct {
	Alpha float64
	Z     float64
	XD    float64
	XB    float64
	Q     float64
}

func New(alpha, zF, xD, xB, q float64) (Shortcut, error) {
	if alpha <= 1 {
		return Shortcut{}, ErrNonSeparating
	}
	if !(xD > zF && zF > xB) {
		return Shortcut{}, ErrInvalidComposition
	}
	if xD <= 0 || xD >= 1 || xB <= 0 || xB >= 1 || zF <= 0 || zF >= 1 {
		return Shortcut{}, ErrInvalidComposition
	}
	return Shortcut{Alpha: alpha, Z: zF, XD: xD, XB: xB, Q: q}, nil
}

func (s Shortcut) FenskeMinimumStages() float64 {
	if s.Alpha <= 1 {
		return math.Inf(1)
	}
	num := math.Log((s.XD / (1 - s.XD)) * ((1 - s.XB) / s.XB))
	return num / math.Log(s.Alpha)
}

func (s Shortcut) UnderwoodTheta() (float64, error) {
	lo := 1.0 + 1e-9
	hi := s.Alpha - 1e-9
	f := func(theta float64) float64 {
		return s.Alpha*s.Z/(s.Alpha-theta) + (1-s.Z)/(1-theta) - (1 - s.Q)
	}
	flo := f(lo)
	fhi := f(hi)
	if flo*fhi > 0 {
		return 0, ErrNoRealRoot
	}
	for i := 0; i < 200; i++ {
		mid := (lo + hi) / 2
		fv := f(mid)
		if math.Abs(fv) < 1e-10 {
			return mid, nil
		}
		if flo*fv < 0 {
			hi = mid
			fhi = fv
		} else {
			lo = mid
			flo = fv
		}
	}
	return (lo + hi) / 2, nil
}

func (s Shortcut) UnderwoodMinReflux() (float64, error) {
	theta, err := s.UnderwoodTheta()
	if err != nil {
		return 0, err
	}
	num := s.Alpha * s.XD / (s.Alpha - theta)
	if num < 0 {
		return 0, ErrNoRealRoot
	}
	return num - 1, nil
}

func (s Shortcut) GillilandStages(r, rMin, nMin float64) float64 {
	if r <= rMin {
		return math.Inf(1)
	}
	x := (r - rMin) / (r + 1)
	if x < 0 {
		x = 0
	}
	// Eduljee correlation: (N - Nmin)/(N + 1) = 0.75*(1 - x^0.5668)
	y := 0.75 * (1 - math.Pow(x, 0.5668))
	if y >= 1 {
		return math.Inf(1)
	}
	return (nMin + y) / (1 - y)
}

func (s Shortcut) RefluxRatioFromGilliland(nStages, nMin, rMin float64) float64 {
	if nStages <= nMin {
		return rMin
	}
	y := (nStages - nMin) / (nStages + 1)
	if y >= 0.75 {
		return rMin
	}
	x := math.Pow(1-y/0.75, 1/0.5668)
	return (x + rMin) / (1 - x)
}

func EffectiveAlpha(v vle.VLE) float64 {
	return vle.RelativeVolatilityFromPoints(v.Points)
}

func (s Shortcut) StagesAtReflux(r float64) float64 {
	if r < 0 {
		return math.Inf(1)
	}
	nMin := s.FenskeMinimumStages()
	if r >= 1e9 {
		return nMin
	}
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		return math.Inf(1)
	}
	if r <= rMin {
		return math.Inf(1)
	}
	return s.GillilandStages(r, rMin, nMin)
}

func (s Shortcut) FenskeMinimumReflux() float64 {
	// approximation from Fenske when R is infinite
	return 0
}

func (s Shortcut) RMinOrZero() float64 {
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		return 0
	}
	return rMin
}
