package shortcut

import (
	"math"
)

type Cost struct {
	StageCost     float64
	EnergyCost    float64
	AnnualStages  float64
	AnnualEnergy  float64
	Total         float64
}

func AnnualCost(nStages, reflux float64, costPerStage, costPerReflux, operatingHours float64) Cost {
	stages := float64(nStages) * costPerStage
	energy := reflux * costPerReflux * operatingHours
	return Cost{
		StageCost:    stages,
		EnergyCost:   energy,
		AnnualStages: stages,
		AnnualEnergy: energy,
		Total:        stages + energy,
	}
}

func OptimumRefluxByCost(s Shortcut, stageCost, energyCost, hours float64) (float64, Cost, error) {
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		return 0, Cost{}, err
	}
	bestR := rMin
	bestCost := math.Inf(1)
	var bestCostStruct Cost
	for i := 1; i <= 200; i++ {
		r := rMin * (1 + float64(i)/100)
		n := s.GillilandStages(r, rMin, s.FenskeMinimumStages())
		if math.IsInf(n, 0) || math.IsNaN(n) {
			continue
		}
		c := AnnualCost(math.Round(n), r, stageCost, energyCost, hours)
		if c.Total < bestCost {
			bestCost = c.Total
			bestR = r
			bestCostStruct = c
		}
	}
	return bestR, bestCostStruct, nil
}

func BreakEvenReflux(rMin, nMin, stageCost, energyCost float64) float64 {
	// marginal cost equality: d(cost)/dR = 0 approximated by finite difference
	prev := math.Inf(1)
	for i := 1; i <= 500; i++ {
		r := rMin * (1 + float64(i)/200)
		n := sFrom(r, rMin, nMin)
		c := float64(int(math.Round(n)))*stageCost + r*energyCost
		if c > prev {
			return rMin * (1 + float64(i-1)/200)
		}
		prev = c
	}
	return rMin * 3.5
}

func sFrom(r, rMin, nMin float64) float64 {
	if r <= rMin {
		return math.Inf(1)
	}
	x := (r - rMin) / (r + 1)
	y := 0.75 * (1 - math.Pow(x, 0.5668))
	if y >= 1 {
		return math.Inf(1)
	}
	return (nMin + y) / (1 - y)
}

func Payback(nStages, costPerStage, savings float64) float64 {
	if savings <= 0 {
		return math.Inf(1)
	}
	return float64(nStages) * costPerStage / savings
}

func CostPerSeparatedMole(total, distillate float64) float64 {
	if distillate <= 0 {
		return 0
	}
	return total / distillate
}
