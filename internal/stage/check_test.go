package stage

import (
	"math"
	"testing"

	"distill-bin/internal/shortcut"
)

func testShortcut() shortcut.Shortcut {
	s, _ := shortcut.New(2.5, 0.45, 0.97, 0.03, 1.0)
	return s
}

func isInf(v float64) bool {
	return math.IsInf(v, 1)
}

func infPos() float64 {
	return math.Inf(1)
}

func TestCheckReflux(t *testing.T) {
	ok := CheckReflux(2.0, 1.4)
	if !ok.OK {
		t.Fatal("R above Rmin should pass")
	}
	bad := CheckReflux(1.0, 1.4)
	if bad.OK {
		t.Fatal("R below Rmin should fail")
	}
}

func TestCheckCompositionOrder(t *testing.T) {
	if err := CheckCompositionOrder(0.45, 0.97, 0.03); err != nil {
		t.Fatalf("valid order should pass: %v", err)
	}
	if err := CheckCompositionOrder(0.45, 0.30, 0.03); err == nil {
		t.Fatal("xD < zF should fail")
	}
}

func TestSafeStagesInfinite(t *testing.T) {
	s := testShortcut()
	if got := SafeStages(s, 0.1); !isInf(got) {
		t.Fatalf("tiny reflux should give infinite stages, got %v", got)
	}
}

func TestRoundedStages(t *testing.T) {
	if RoundedStages(7.6) != 8 {
		t.Fatal("round(7.6) want 8")
	}
	if RoundedStages(infPos()) != 0 {
		t.Fatal("infinite should map to 0")
	}
}
