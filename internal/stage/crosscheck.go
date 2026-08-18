package stage

import (
	"math"

	"distill-bin/internal/shortcut"
	"distill-bin/internal/vle"
)

type CompareResult struct {
	StepStages  int
	Gilliland   float64
	Difference  float64
	Agree       bool
}

func CompareStepVsGilliland(v vle.VLE, s shortcut.Shortcut, r float64) CompareResult {
	step, err := TotalStagesAt(v, r, s.XD, s.XB, s.Q, s.Z)
	stepN := 0
	if err == nil {
		stepN = step
	}
	rMin, _ := s.UnderwoodMinReflux()
	nMin := s.FenskeMinimumStages()
	gil := s.GillilandStages(r, rMin, nMin)
	diff := math.Abs(float64(stepN) - gil)
	return CompareResult{
		StepStages: stepN,
		Gilliland:  gil,
		Difference: diff,
		Agree:      diff <= 3,
	}
}

func StepAtFullReflux(v vle.VLE, s shortcut.Shortcut) (int, float64) {
	n, err := TotalStagesAt(v, 1e9, s.XD, s.XB, s.Q, s.Z)
	if err != nil {
		return 0, 0
	}
	nMin := s.FenskeMinimumStages()
	return n, nMin
}

func RefluxFromStages(v vle.VLE, s shortcut.Shortcut, targetStages int) float64 {
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		return 0
	}
	for i := 1; i <= 2000; i++ {
		r := rMin * (1 + float64(i)/1000)
		n, err := TotalStagesAt(v, r, s.XD, s.XB, s.Q, s.Z)
		if err == nil && n <= targetStages {
			return r
		}
	}
	return 0
}

func AgreementTolerance(v vle.VLE, s shortcut.Shortcut) float64 {
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		return 0
	}
	step, gil := 0.0, 0.0
	n, err := TotalStagesAt(v, 2*rMin, s.XD, s.XB, s.Q, s.Z)
	if err == nil {
		step = float64(n)
	}
	gil = s.GillilandStages(2*rMin, rMin, s.FenskeMinimumStages())
	if gil == 0 {
		return 0
	}
	return math.Abs(step-gil) / gil
}
