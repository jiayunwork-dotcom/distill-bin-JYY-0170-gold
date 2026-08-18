package shortcut

import (
	"math"
)

type Design struct {
	Shortcut Shortcut
	RMin     float64
	NMin     float64
	N        float64
	R        float64
	Feasible bool
}

func DesignAtReflux(s Shortcut, r float64) (Design, error) {
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		return Design{}, err
	}
	nMin := s.FenskeMinimumStages()
	if r < rMin {
		return Design{
			Shortcut: s,
			RMin:     rMin,
			NMin:     nMin,
			R:        r,
			Feasible: false,
		}, nil
	}
	n := s.GillilandStages(r, rMin, nMin)
	return Design{
		Shortcut: s,
		RMin:     rMin,
		NMin:     nMin,
		N:        n,
		R:        r,
		Feasible: !math.IsInf(n, 0) && !math.IsNaN(n),
	}, nil
}

func OptimumRefluxFactor(rMin float64) float64 {
	if rMin < 0 {
		return 0
	}
	return 1.2 * rMin
}

func TotalStagesAt(s Shortcut, r float64) (int, error) {
	d, err := DesignAtReflux(s, r)
	if err != nil {
		return 0, err
	}
	if !d.Feasible {
		return 0, nil
	}
	return int(math.Round(d.N)), nil
}

func FeedStageFraction(s Shortcut, q float64) float64 {
	if q <= 0 {
		return 0.5
	}
	if q >= 1 {
		return 0.4
	}
	return q * 0.5
}

func StagesBelowFeed(total int, q float64) int {
	if total <= 0 {
		return 0
	}
	frac := FeedStageFraction(Shortcut{Q: q}, q)
	below := int(math.Round(float64(total) * frac))
	if below > total-1 {
		below = total - 1
	}
	return below
}
