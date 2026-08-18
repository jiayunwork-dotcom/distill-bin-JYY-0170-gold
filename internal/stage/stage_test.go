package stage

import (
	"errors"
	"math"
	"testing"

	"distill-bin/internal/vle"
)

func TestStepOffConverges(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	res, err := StepOff(v, 2.5, 0.97, 0.03, 1.0, 0.45, 1000)
	if err != nil {
		t.Fatalf("StepOff: %v", err)
	}
	if !res.Converged {
		t.Fatal("should converge")
	}
	if res.Stages < 3 {
		t.Fatalf("stages too few: %d", res.Stages)
	}
	if !res.Reboiler {
		t.Fatal("reboiler flag should be set")
	}
}

func TestStepOffInsufficientReflux(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	_, err := StepOff(v, 0.5, 0.97, 0.03, 1.0, 0.45, 1000)
	if !errors.Is(err, ErrInsufficientReflux) {
		t.Fatalf("want ErrInsufficientReflux, got %v", err)
	}
}

func TestFullRefluxApproachesFenske(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	res, err := StepOff(v, 1e9, 0.97, 0.03, 1.0, 0.45, 1000)
	if err != nil {
		t.Fatalf("StepOff full reflux: %v", err)
	}
	// Fenske Nmin for xD=0.97, xB=0.03, alpha=2.5
	num := math.Log((0.97/0.03)*(0.97/0.03))
	nMin := num / math.Log(2.5)
	diff := math.Abs(float64(res.Stages) - nMin)
	if diff > 2 {
		t.Fatalf("full reflux stages=%d too far from Fenske Nmin=%.2f", res.Stages, nMin)
	}
}

func TestStagesDoNotRiseWithReflux(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	n1, err := TotalStagesAt(v, 1.5, 0.97, 0.03, 1.0, 0.45)
	if err != nil {
		t.Fatalf("TotalStagesAt r=1.5: %v", err)
	}
	n2, err := TotalStagesAt(v, 3.0, 0.97, 0.03, 1.0, 0.45)
	if err != nil {
		t.Fatalf("TotalStagesAt r=3.0: %v", err)
	}
	if n2 > n1 {
		t.Fatalf("higher reflux needs %d stages, more than %d", n2, n1)
	}
}

func TestRectifyingLine(t *testing.T) {
	l := RectifyingLine(2.0, 0.9)
	if math.Abs(l.Slope-2.0/3.0) > 1e-12 {
		t.Fatalf("slope want 2/3, got %v", l.Slope)
	}
	if math.Abs(l.Intercept-0.9/3.0) > 1e-12 {
		t.Fatalf("intercept want 0.3, got %v", l.Intercept)
	}
}

func TestQLineSaturatedLiquid(t *testing.T) {
	l := QLine(1.0, 0.5)
	if !math.IsInf(l.Slope, 1) {
		t.Fatalf("q=1 should be a vertical line, got slope %v", l.Slope)
	}
	if math.Abs(l.Intercept-0.5) > 1e-12 {
		t.Fatalf("q=1 intercept want 0.5, got %v", l.Intercept)
	}
}

func TestStrippingLineEndsAtBottoms(t *testing.T) {
	l := StrippingLine(2.5, 0.97, 1.0, 0.45, 0.03)
	y := l.Slope*0.03 + l.Intercept
	if math.Abs(y-0.03) > 1e-9 {
		t.Fatalf("stripping line at xB should equal xB: got %v", y)
	}
}

func TestProfileMonotonic(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	p, err := CompositionProfile(v, 2.5, 0.97, 0.03, 1.0, 0.45, 1000)
	if err != nil {
		t.Fatalf("CompositionProfile: %v", err)
	}
	if len(p.Trays) == 0 {
		t.Fatal("empty profile")
	}
	if !MonotonicDecreasing(p) {
		t.Fatal("tray compositions should decrease from distillate to bottoms")
	}
}

func TestLastCompositionReachesBottoms(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	p, err := CompositionProfile(v, 2.5, 0.97, 0.03, 1.0, 0.45, 1000)
	if err != nil {
		t.Fatalf("CompositionProfile: %v", err)
	}
	last := LastComposition(p)
	if last > 0.03+1e-3 {
		t.Fatalf("last composition %.4f should be near xB=0.03", last)
	}
}

func TestMinimumTraysWithReboiler(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	n := MinimumTraysWithReboiler(v, 1e9, 0.97, 0.03, 1.0, 0.45)
	if n < 2 {
		t.Fatalf("minimum trays too few: %d", n)
	}
}
