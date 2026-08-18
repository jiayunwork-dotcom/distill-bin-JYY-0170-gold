package stage

import (
	"math"

	"distill-bin/internal/vle"
)

type OperatingPoint struct {
	R          float64
	RMin       float64
	NMin       float64
	N          float64
	Q          float64
	Z          float64
	Feasible   bool
}

func OperatingPointAt(v vle.VLE, r, xD, xB, q, zF, rMin, nMin float64) OperatingPoint {
	n, err := TotalStagesAt(v, r, xD, xB, q, zF)
	if err != nil {
		return OperatingPoint{R: r, RMin: rMin, NMin: nMin, Q: q, Z: zF, Feasible: false}
	}
	return OperatingPoint{R: r, RMin: rMin, NMin: nMin, N: float64(n), Q: q, Z: zF, Feasible: true}
}

func RefluxMultiplier(r, rMin float64) float64 {
	if rMin <= 0 {
		return 0
	}
	return r / rMin
}

func AbsoluteGain(n1, n2 float64) float64 {
	return n1 - n2
}

func PercentGain(n1, n2 float64) float64 {
	if n1 <= 0 {
		return 0
	}
	return (n1 - n2) / n1 * 100
}

func CrossOverReflux(v vle.VLE, xD, xB, q, zF, rMin float64) float64 {
	// binary search for the reflux where stages stop dropping fast
	lo := rMin * 1.01
	hi := rMin * 10
	prev := math.Inf(1)
	for i := 0; i < 60; i++ {
		mid := (lo + hi) / 2
		n, err := TotalStagesAt(v, mid, xD, xB, q, zF)
		if err != nil {
			return mid
		}
		cur := float64(n)
		if cur >= prev {
			return mid
		}
		prev = cur
		lo = mid
	}
	return (lo + hi) / 2
}

func StableReflux(study RefluxStudy, tolerance int) float64 {
	for i := 1; i < len(study.Refluxes); i++ {
		if study.Stages[i] > 0 && study.Stages[i-1] > 0 {
			diff := study.Stages[i-1] - study.Stages[i]
			if diff <= tolerance {
				return study.Refluxes[i]
			}
		}
	}
	if len(study.Refluxes) == 0 {
		return 0
	}
	return study.Refluxes[len(study.Refluxes)-1]
}
