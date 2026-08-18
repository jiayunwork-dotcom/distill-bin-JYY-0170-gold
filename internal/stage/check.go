package stage

import (
	"errors"
	"math"

	"distill-bin/internal/shortcut"
)

var ErrRefluxBelowMin = errors.New("stage: reflux below minimum")

type Check struct {
	R     float64
	RMin  float64
	OK    bool
	Reason string
}

func CheckReflux(r, rMin float64) Check {
	if r < rMin {
		return Check{R: r, RMin: rMin, OK: false, Reason: "reflux below minimum"}
	}
	return Check{R: r, RMin: rMin, OK: true}
}

func CheckCompositionOrder(zF, xD, xB float64) error {
	if !(xD > zF && zF > xB) {
		return errors.New("stage: distillate must exceed feed which must exceed bottoms")
	}
	return nil
}

func SafeStages(s shortcut.Shortcut, r float64) float64 {
	if r < s.RMinOrZero() {
		return math.Inf(1)
	}
	nMin := s.FenskeMinimumStages()
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		return math.Inf(1)
	}
	return s.GillilandStages(r, rMin, nMin)
}

func RoundedStages(n float64) int {
	if math.IsInf(n, 1) || math.IsNaN(n) {
		return 0
	}
	return int(math.Round(n))
}

func FeasibleR(r, rMin float64) bool {
	return r >= rMin
}

func MarginRatio(r, rMin float64) float64 {
	if rMin <= 0 {
		return 0
	}
	return r / rMin
}
