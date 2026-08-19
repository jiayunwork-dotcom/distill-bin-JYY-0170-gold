package stage

import (
	"math"

	"distill-bin/internal/vle"
)

type Tray struct {
	Number   int
	X        float64
	Y        float64
	Section  string
}

type Profile struct {
	Trays []Tray
	FeedTray int
}

func CompositionProfile(v vle.VLE, r, xD, xB, q, zF float64, maxStages int) (Profile, error) {
	rect := RectifyingLine(r, xD)
	xi := intersection(rect, QLine(q, zF))
	strip := StrippingLine(r, xD, q, zF, xB)
	x := v.X(xD)
	trays := []Tray{}
	feedTray := 0
	for n := 0; n < maxStages; n++ {
		if x <= xB {
			break
		}
		section := "rectifying"
		if x <= xi.X {
			section = "stripping"
			if feedTray == 0 {
				feedTray = n
			}
		}
		var y float64
		if section == "rectifying" {
			y = rect.Slope*x + rect.Intercept
		} else {
			y = strip.Slope*x + strip.Intercept
		}
		xNext := v.X(y)
		if xNext <= 0 || xNext > 1 || math.IsNaN(xNext) {
			return Profile{}, ErrStepFailed
		}
		trays = append(trays, Tray{Number: n + 1, X: xNext, Y: y, Section: section})
		x = xNext
	}
	return Profile{Trays: trays, FeedTray: feedTray}, nil
}

func RectifyingTrays(p Profile) int {
	n := 0
	for _, t := range p.Trays {
		if t.Section == "rectifying" {
			n++
		}
	}
	return n
}

func StrippingTrays(p Profile) int {
	n := 0
	for _, t := range p.Trays {
		if t.Section == "stripping" {
			n++
		}
	}
	return n
}

func LastComposition(p Profile) float64 {
	if len(p.Trays) == 0 {
		return 0
	}
	return p.Trays[len(p.Trays)-1].X
}

func MonotonicDecreasing(p Profile) bool {
	prev := math.Inf(1)
	for _, t := range p.Trays {
		if t.X >= prev {
			return false
		}
		prev = t.X
	}
	return true
}

func MinimumTraysWithReboiler(v vle.VLE, r, xD, xB, q, zF float64) int {
	res, err := StepOff(v, r, xD, xB, q, zF, 1000)
	if err != nil {
		return 0
	}
	return res.Stages
}
