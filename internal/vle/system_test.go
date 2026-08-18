package vle

import (
	"math"
	"testing"
)

func TestSystemPartialPressure(t *testing.T) {
	s := NewSystem(2.5, 100)
	pA, pB := s.PartialPressure(0.4)
	if math.Abs(pA-100) > 1e-9 {
		t.Fatalf("pA want 100, got %v", pA)
	}
	if math.Abs(pB-60) > 1e-9 {
		t.Fatalf("pB want 60, got %v", pB)
	}
}

func TestSystemTotalPressure(t *testing.T) {
	s := NewSystem(2.5, 100)
	p := s.TotalPressure(0.4)
	if math.Abs(p-160) > 1e-9 {
		t.Fatalf("total pressure want 160, got %v", p)
	}
}

func TestSystemConsistency(t *testing.T) {
	s := NewSystem(2.5, 100)
	y := s.BubbleY(0.3)
	if !s.CheckThermodynamicConsistency(0.3, y) {
		t.Fatal("bubble y should be thermodynamically consistent")
	}
	x := s.DewX(y)
	if math.Abs(x-0.3) > 1e-9 {
		t.Fatalf("dew x want 0.3, got %v", x)
	}
}

func TestSystemKValue(t *testing.T) {
	s := NewSystem(2.5, 100)
	k := s.KValue(0.5)
	if k <= 0 || k > 2.5 {
		t.Fatalf("K value out of range: %v", k)
	}
}

func TestSystemVolatilityConstant(t *testing.T) {
	s := NewSystem(2.5, 100)
	if !s.Constant() {
		t.Fatal("constant alpha system should report constant volatility")
	}
	a := s.RelativeVolatilityAt(0.5)
	if math.Abs(a-2.5) > 1e-9 {
		t.Fatalf("volatility at x=0.5 want 2.5, got %v", a)
	}
}

func TestSystemVolatilityRatio(t *testing.T) {
	s := NewSystem(2.5, 100)
	r := s.VolatilityRatio(0.2, 0.8)
	if math.Abs(r-1) > 1e-9 {
		t.Fatalf("constant alpha ratio should be 1, got %v", r)
	}
}

func TestSystemBubbleTemp(t *testing.T) {
	s := NewSystem(2.5, 100)
	ratio := s.BubbleTemperature(250, 100)
	if math.Abs(ratio-2.5) > 1e-9 {
		t.Fatalf("temperature ratio want 2.5, got %v", ratio)
	}
}
