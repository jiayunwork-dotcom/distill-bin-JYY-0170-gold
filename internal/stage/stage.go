package stage

import (
	"errors"
	"math"

	"distill-bin/internal/vle"
)

var (
	ErrInsufficientReflux = errors.New("stage: reflux below minimum, tower will not converge")
	ErrStepFailed         = errors.New("stage: stepping failed to reach bottoms composition")
)

type OperatingLine struct {
	Slope     float64
	Intercept float64
}

func RectifyingLine(r, xD float64) OperatingLine {
	if r <= 0 {
		return OperatingLine{Slope: 0, Intercept: xD}
	}
	return OperatingLine{
		Slope:     r / (r + 1),
		Intercept: xD / (r + 1),
	}
}

func QLine(q, zF float64) OperatingLine {
	if q == 1 {
		// vertical line x = zF; represented by slope +Inf convention
		return OperatingLine{Slope: math.Inf(1), Intercept: zF}
	}
	return OperatingLine{
		Slope:     q / (q - 1),
		Intercept: -zF / (q - 1),
	}
}

func StrippingLine(r, xD, q, zF, xB float64) OperatingLine {
	rect := RectifyingLine(r, xD)
	xi := intersection(rect, QLine(q, zF))
	slope := (xi.Y - xB) / (xi.X - xB)
	return OperatingLine{Slope: slope, Intercept: xB - slope*xB}
}

func intersection(a, b OperatingLine) (pt struct{ X, Y float64 }) {
	if math.IsInf(b.Slope, 1) {
		// b is the vertical q-line x = b.Intercept
		pt.X = b.Intercept
		pt.Y = a.Slope*pt.X + a.Intercept
		return
	}
	if math.Abs(a.Slope-b.Slope) < 1e-15 {
		pt.X = 0
		pt.Y = a.Intercept
		return
	}
	pt.X = (b.Intercept - a.Intercept) / (a.Slope - b.Slope)
	pt.Y = a.Slope*pt.X + a.Intercept
	return
}

type StepResult struct {
	Stages     int
	TrayCompositions []float64
	Converged  bool
	Reboiler   bool
}

func StepOff(v vle.VLE, r, xD, xB, q, zF float64, maxStages int) (StepResult, error) {
	if !v.ValidComposition(xD) || !v.ValidComposition(xB) {
		return StepResult{}, vle.ErrCompositionOutOfRange
	}
	rect := RectifyingLine(r, xD)
	xi := intersection(rect, QLine(q, zF))
	strip := StrippingLine(r, xD, q, zF, xB)

	// minimum reflux check: q-line intersection must lie under the equilibrium curve
	ye := v.Y(xi.X)
	if xi.Y > ye {
		return StepResult{}, ErrInsufficientReflux
	}

	x := v.X(xD)
	compositions := []float64{xD}
	stages := 1
	for stages <= maxStages {
		if x <= xB {
			return StepResult{
				Stages:           stages,
				TrayCompositions: compositions,
				Converged:        true,
				Reboiler:         true,
			}, nil
		}
		var y float64
		if x > xi.X {
			y = rect.Slope*x + rect.Intercept
		} else {
			y = strip.Slope*x + strip.Intercept
		}
		xNext := v.X(y)
		if xNext <= 0 || xNext > 1 || math.IsNaN(xNext) {
			return StepResult{}, ErrStepFailed
		}
		stages++
		compositions = append(compositions, xNext)
		x = xNext
	}
	return StepResult{}, ErrStepFailed
}

func TotalStagesAt(v vle.VLE, r, xD, xB, q, zF float64) (int, error) {
	res, err := StepOff(v, r, xD, xB, q, zF, 1000)
	if err != nil {
		return 0, err
	}
	return res.Stages, nil
}

func StagesWithoutReboiler(v vle.VLE, r, xD, xB, q, zF float64) (int, error) {
	res, err := StepOff(v, r, xD, xB, q, zF, 1000)
	if err != nil {
		return 0, err
	}
	if res.Stages > 0 {
		return res.Stages - 1, nil
	}
	return 0, nil
}
