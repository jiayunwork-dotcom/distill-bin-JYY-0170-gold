package stage

import (
	"math"

	"distill-bin/internal/vle"
)

type RefluxStudy struct {
	Refluxes []float64
	Stages   []int
}

func RefluxSweep(v vle.VLE, xD, xB, q, zF float64, from, to, steps int) RefluxStudy {
	if steps < 2 {
		steps = 2
	}
	study := RefluxStudy{}
	for i := 0; i <= steps; i++ {
		r := float64(from) + float64(i)*(float64(to-from))/float64(steps)
		n, err := TotalStagesAt(v, r, xD, xB, q, zF)
		study.Refluxes = append(study.Refluxes, r)
		if err != nil {
			study.Stages = append(study.Stages, 0)
		} else {
			study.Stages = append(study.Stages, n)
		}
	}
	return study
}

func MonotoneInReflux(study RefluxStudy) bool {
	prev := 0
	started := false
	for _, n := range study.Stages {
		if n <= 0 {
			continue
		}
		if started && n > prev {
			return false
		}
		prev = n
		started = true
	}
	return true
}

func BestRefluxInSweep(study RefluxStudy) (float64, int) {
	best := 0.0
	bestStages := math.MaxInt
	for i, r := range study.Refluxes {
		if study.Stages[i] > 0 && study.Stages[i] < bestStages {
			bestStages = study.Stages[i]
			best = r
		}
	}
	return best, bestStages
}

func TotalAnnualized(nStages, stageCost, reflux, energyCost float64) float64 {
	return float64(nStages)*stageCost + reflux*energyCost
}
