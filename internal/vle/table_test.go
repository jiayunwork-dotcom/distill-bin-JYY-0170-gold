package vle

import (
	"math"
	"testing"
)

func TestEquilibriumTable(t *testing.T) {
	v, _ := NewConstantAlpha(2.5)
	rows := EquilibriumTable(v, 0, 10, 10)
	if len(rows) != 11 {
		t.Fatalf("want 11 rows, got %d", len(rows))
	}
	for _, r := range rows {
		if r.X < 0 || r.X > 1 {
			continue
		}
		if r.Y < r.X {
			t.Fatalf("light component y(%.3f)=%.3f should exceed x", r.X, r.Y)
		}
	}
}

func TestCompositionRange(t *testing.T) {
	v, _ := NewConstantAlpha(2.5)
	lo, hi := CompositionRange(v, 0.5)
	if lo >= hi {
		t.Fatalf("range invalid: %v %v", lo, hi)
	}
}

func TestPinchComposition(t *testing.T) {
	v, _ := NewConstantAlpha(2.5)
	pinch := PinchComposition(v, 1.0, 0.45)
	if pinch <= 0 || pinch >= 1 {
		t.Fatalf("pinch out of range: %v", pinch)
	}
}

func TestVaporPressureAntoine(t *testing.T) {
	p := VaporPressureAntoine(250, 8.1, 1730, 233)
	if p <= 0 || math.IsInf(p, 0) {
		t.Fatalf("vapor pressure invalid: %v", p)
	}
}

func TestBubbleDew(t *testing.T) {
	alpha := 2.5
	x := 0.4
	y := BubblePoint(alpha, x)
	xBack := DewPoint(alpha, y)
	if math.Abs(xBack-x) > 1e-9 {
		t.Fatalf("bubble/dew round trip failed: %v -> %v -> %v", x, y, xBack)
	}
}

func TestSeparationFactor(t *testing.T) {
	sf := SeparationFactor(0.97, 0.03)
	if sf <= 1 {
		t.Fatalf("separation factor should be large, got %v", sf)
	}
}

func TestVaporLiquidRates(t *testing.T) {
	if VaporRate(50, 2.5) != 175 {
		t.Fatalf("V = D(R+1) want 175, got %v", VaporRate(50, 2.5))
	}
	if LiquidRate(50, 2.5) != 125 {
		t.Fatalf("L = DR want 125, got %v", LiquidRate(50, 2.5))
	}
}

func TestEnrichmentRatio(t *testing.T) {
	v, _ := NewConstantAlpha(2.5)
	r := EnrichmentRatio(2.5, 0.3)
	if r <= 1 {
		t.Fatalf("enrichment ratio should exceed 1, got %v", r)
	}
	_ = v
}

func TestLogMeanSeparation(t *testing.T) {
	if LogMeanSeparation(2.5, 3) != math.Pow(2.5, 3) {
		t.Fatal("log mean separation wrong")
	}
}
