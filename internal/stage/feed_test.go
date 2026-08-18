package stage

import (
	"math"
	"testing"

	"distill-bin/internal/vle"
)

func TestFeedType(t *testing.T) {
	if FeedType(0.0) != "saturated_vapor" {
		t.Fatalf("q=0 want saturated_vapor, got %s", FeedType(0))
	}
	if FeedType(1.0) != "saturated_liquid" {
		t.Fatalf("q=1 want saturated_liquid, got %s", FeedType(1))
	}
	if FeedType(1.5) != "subcooled_liquid" {
		t.Fatalf("q=1.5 want subcooled_liquid, got %s", FeedType(1.5))
	}
	if FeedType(-0.5) != "superheated_vapor" {
		t.Fatalf("q<0 want superheated_vapor, got %s", FeedType(-0.5))
	}
}

func TestMinimumRefluxByIntersection(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	ok := MinimumRefluxByIntersection(v, 2.5, 0.97, 0.45, 1.0)
	if !ok {
		t.Fatal("R=2.5 should be above minimum at q=1")
	}
	bad := MinimumRefluxByIntersection(v, 0.1, 0.97, 0.45, 1.0)
	if bad {
		t.Fatal("R=0.1 should be below minimum")
	}
}

func TestOConnellEfficiency(t *testing.T) {
	e := OConnellEfficiency(2.5, 0.3)
	if e <= 0 || e > 1 {
		t.Fatalf("O'Connell efficiency out of range: %v", e)
	}
}

func TestActualPlates(t *testing.T) {
	if ActualPlates(10, 0.7) != 15 {
		t.Fatalf("ceil(10/0.7)=15, got %d", ActualPlates(10, 0.7))
	}
}

func TestRefluxSweepMonotone(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	study := RefluxSweep(v, 0.97, 0.03, 1.0, 0.45, 1, 8, 4)
	if !MonotoneInReflux(study) {
		t.Fatalf("stages should not rise with reflux: %+v", study.Stages)
	}
}

func TestBestRefluxInSweep(t *testing.T) {
	study := RefluxStudy{
		Refluxes: []float64{1, 2, 3, 4},
		Stages:   []int{0, 15, 12, 11},
	}
	r, n := BestRefluxInSweep(study)
	if r != 4 || n != 11 {
		t.Fatalf("best reflux want 4/11, got %v/%d", r, n)
	}
}

func TestMurphreeEfficiency(t *testing.T) {
	e := MurphreeEfficiency(0.7, 0.8, 0.4)
	want := (0.7 - 0.4) / (0.8 - 0.4)
	if math.Abs(e-want) > 1e-9 {
		t.Fatalf("Murphree want %v, got %v", want, e)
	}
}

func TestOperatingPointAtInfeasible(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	op := OperatingPointAt(v, 0.05, 0.97, 0.03, 1.0, 0.45, 1.4, 7.6)
	if op.Feasible {
		t.Fatal("tiny reflux should be infeasible")
	}
}

func TestCrossOverReflux(t *testing.T) {
	v, _ := vle.NewConstantAlpha(2.5)
	r := CrossOverReflux(v, 0.97, 0.03, 1.0, 0.45, 1.4)
	if r < 1.4 {
		t.Fatalf("crossover reflux should exceed Rmin, got %v", r)
	}
}
