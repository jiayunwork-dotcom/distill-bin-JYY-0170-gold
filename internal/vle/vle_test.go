package vle

import (
	"errors"
	"math"
	"testing"
)

func TestConstantAlphaEquilibrium(t *testing.T) {
	v, err := NewConstantAlpha(2.5)
	if err != nil {
		t.Fatalf("NewConstantAlpha: %v", err)
	}
	x := 0.5
	want := 2.5 * 0.5 / (1 + 1.5*0.5)
	if math.Abs(v.Y(x)-want) > 1e-12 {
		t.Fatalf("y(0.5) want %v, got %v", want, v.Y(x))
	}
	if math.Abs(v.X(v.Y(x))-x) > 1e-9 {
		t.Fatalf("round trip failed: %v -> %v -> %v", x, v.Y(x), v.X(v.Y(x)))
	}
}

func TestRejectAlphaAtMostOne(t *testing.T) {
	if _, err := NewConstantAlpha(1.0); !errors.Is(err, ErrAlphaNonSeparating) {
		t.Fatalf("alpha=1 should fail, got %v", err)
	}
	if _, err := NewConstantAlpha(0.8); !errors.Is(err, ErrAlphaNonSeparating) {
		t.Fatalf("alpha=0.8 should fail, got %v", err)
	}
}

func TestCurveInterpolation(t *testing.T) {
	points := []CurvePoint{
		{X: 0.0, Y: 0.0},
		{X: 0.5, Y: 0.7},
		{X: 1.0, Y: 1.0},
	}
	v, err := NewCurve(points)
	if err != nil {
		t.Fatalf("NewCurve: %v", err)
	}
	if math.Abs(v.Y(0.25)-0.35) > 1e-9 {
		t.Fatalf("interpolated y want 0.35, got %v", v.Y(0.25))
	}
}

func TestRelativeVolatilityEstimate(t *testing.T) {
	points := []CurvePoint{
		{X: 0.3, Y: 0.6},
		{X: 0.5, Y: 0.75},
	}
	alpha := RelativeVolatilityFromPoints(points)
	if alpha <= 1 {
		t.Fatalf("estimated alpha should exceed 1, got %v", alpha)
	}
}

func TestSolveBalanceClosed(t *testing.T) {
	res, err := SolveBalance(100, 0.5, 0.95, 0.05)
	if err != nil {
		t.Fatalf("SolveBalance: %v", err)
	}
	if math.Abs(res.Distillate-50) > 1e-9 {
		t.Fatalf("D want 50, got %v", res.Distillate)
	}
	if math.Abs(res.Bottoms-50) > 1e-9 {
		t.Fatalf("B want 50, got %v", res.Bottoms)
	}
	if !res.Closed {
		t.Fatalf("balance should close, residual %v", res.Residual)
	}
}

func TestBalanceCompositionErrors(t *testing.T) {
	if _, err := SolveBalance(100, 0.5, 0.3, 0.05); !errors.Is(err, ErrInvalidFeed) {
		t.Fatalf("xD < zF should fail, got %v", err)
	}
	if _, err := SolveBalance(100, 0.5, 0.95, 0.6); !errors.Is(err, ErrInvalidFeed) {
		t.Fatalf("xB > zF should fail, got %v", err)
	}
	if _, err := SolveBalance(-5, 0.5, 0.95, 0.05); err == nil {
		t.Fatal("negative feed should fail")
	}
}

func TestRecoveries(t *testing.T) {
	light, heavy := Recoveries(100, 0.5, 0.95, 0.05)
	if light <= 0 || light > 1 {
		t.Fatalf("light recovery out of range: %v", light)
	}
	if heavy <= 0 || heavy > 1 {
		t.Fatalf("heavy recovery out of range: %v", heavy)
	}
}
