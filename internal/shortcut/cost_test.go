package shortcut

import (
	"math"
	"testing"
)

func TestAnnualCost(t *testing.T) {
	c := AnnualCost(10, 2.0, 1000, 500, 8000)
	if c.Total != 10*1000+2.0*500*8000 {
		t.Fatalf("cost wrong: %v", c.Total)
	}
}

func TestOptimumRefluxByCost(t *testing.T) {
	s, _ := New(2.5, 0.45, 0.97, 0.03, 1.0)
	r, c, err := OptimumRefluxByCost(s, 1000, 100, 8000)
	if err != nil {
		t.Fatalf("OptimumRefluxByCost: %v", err)
	}
	if r <= 0 || c.Total <= 0 {
		t.Fatalf("optimum invalid: r=%v cost=%+v", r, c)
	}
}

func TestBreakEvenReflux(t *testing.T) {
	r := BreakEvenReflux(1.4, 7.6, 1000, 100)
	if r < 1.4 {
		t.Fatalf("break-even reflux should exceed Rmin, got %v", r)
	}
}

func TestRminVsFeedQuality(t *testing.T) {
	s, _ := New(2.5, 0.45, 0.97, 0.03, 1.0)
	study, err := RminAtFeedQuality(s, 0.8)
	if err != nil {
		t.Fatalf("RminAtFeedQuality: %v", err)
	}
	if study.RMin <= 0 {
		t.Fatalf("Rmin should be positive: %v", study.RMin)
	}
}

func TestRminSensitivity(t *testing.T) {
	s, _ := New(2.5, 0.45, 0.97, 0.03, 1.0)
	sen, err := RminSensitivity(s)
	if err != nil {
		t.Fatalf("RminSensitivity: %v", err)
	}
	if math.Abs(sen) > 1e6 {
		t.Fatalf("sensitivity absurd: %v", sen)
	}
}

func TestNminVsSeparation(t *testing.T) {
	n := NminVsSeparation(2.5, 0.97, 0.03)
	if n <= 0 {
		t.Fatalf("Nmin should be positive: %v", n)
	}
}

func TestCheckDesignPair(t *testing.T) {
	s, _ := New(2.5, 0.45, 0.97, 0.03, 1.0)
	rMin, _ := s.UnderwoodMinReflux()
	n := s.GillilandStages(rMin*2, rMin, s.FenskeMinimumStages())
	if !CheckDesignPair(s, rMin*2, n) {
		t.Fatal("design pair should be consistent")
	}
}

func TestPayback(t *testing.T) {
	p := Payback(10, 1000, 5000)
	if math.Abs(p-2) > 1e-9 {
		t.Fatalf("payback want 2 years, got %v", p)
	}
}

func TestCostPerSeparatedMole(t *testing.T) {
	if CostPerSeparatedMole(100, 50) != 2 {
		t.Fatal("cost per mole wrong")
	}
}
