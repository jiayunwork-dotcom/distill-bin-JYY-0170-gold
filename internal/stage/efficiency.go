package stage

import (
	"math"

	"distill-bin/internal/vle"
)

type PlateEfficiency struct {
	Method string
	Value  float64
}

func MurphreeEfficiency(yActual, yEquilibrium, yIn float64) float64 {
	if yEquilibrium == yIn {
		return 1
	}
	return (yActual - yIn) / (yEquilibrium - yIn)
}

func OConnellEfficiency(alpha, mu float64) float64 {
	if alpha*mu <= 0 {
		return 0.5
	}
	return 0.492 / math.Pow(alpha*mu, 0.245)
}

func ActualPlates(theoretical, efficiency float64) int {
	if efficiency <= 0 {
		return 0
	}
	return int(math.Ceil(theoretical / efficiency))
}

func StagesWithEfficiency(v vle.VLE, r, xD, xB, q, zF, eff float64) (int, error) {
	res, err := StepOff(v, r, xD, xB, q, zF, 1000)
	if err != nil {
		return 0, err
	}
	return ActualPlates(float64(res.Stages), eff), nil
}

func EfficiencyToAchieve(theoretical, actual float64) float64 {
	if actual <= 0 {
		return 0
	}
	return theoretical / actual
}

func WeepingMargin(weep, op float64) float64 {
	return op / weep
}

func FloodingMargin(flood, op float64) float64 {
	return flood / op
}
