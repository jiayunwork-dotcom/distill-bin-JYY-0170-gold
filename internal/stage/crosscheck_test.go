package stage

import (
	"math"
	"testing"

	"distill-bin/internal/shortcut"
	"distill-bin/internal/vle"
)

func TestCompareStepVsGilliland(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	s, _ := shortcut.New(2.5, 0.45, 0.97, 0.03, 1.0)
	rMin, _ := s.UnderwoodMinReflux()
	res := CompareStepVsGilliland(v, s, rMin*2)
	if !res.Agree {
		t.Fatalf("step (%d) vs Gilliland (%.2f) disagree by %.2f", res.StepStages, res.Gilliland, res.Difference)
	}
}

func TestStepAtFullReflux(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	s, _ := shortcut.New(2.5, 0.45, 0.97, 0.03, 1.0)
	n, nMin := StepAtFullReflux(v, s)
	if math.Abs(float64(n)-nMin) > 2 {
		t.Fatalf("full reflux stages %d too far from Fenske %.2f", n, nMin)
	}
}

func TestRefluxFromStages(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	s, _ := shortcut.New(2.5, 0.45, 0.97, 0.03, 1.0)
	r := RefluxFromStages(v, s, 15)
	if r <= 0 {
		t.Fatalf("reflux from stages invalid: %v", r)
	}
	n, err := TotalStagesAt(v, r, s.XD, s.XB, s.Q, s.Z)
	if err != nil {
		t.Fatalf("TotalStagesAt: %v", err)
	}
	if n > 15 {
		t.Fatalf("stage count %d should be <= 15", n)
	}
}

func TestAgreementTolerance(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	s, _ := shortcut.New(2.5, 0.45, 0.97, 0.03, 1.0)
	tol := AgreementTolerance(v, s)
	if tol > 0.5 {
		t.Fatalf("step vs Gilliland tolerance too large: %v", tol)
	}
}
