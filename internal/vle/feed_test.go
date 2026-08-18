package vle

import (
	"math"
	"testing"
)

func TestFeedPhases(t *testing.T) {
	f := Feed{Rate: 100, Z: 0.5, Q: 0.7}
	if math.Abs(f.Vapor()-30) > 1e-9 {
		t.Fatalf("vapor want 30, got %v", f.Vapor())
	}
	if math.Abs(f.Liquid()-70) > 1e-9 {
		t.Fatalf("liquid want 70, got %v", f.Liquid())
	}
}

func TestProductLight(t *testing.T) {
	p := Product{Rate: 40, X: 0.9}
	if math.Abs(p.Light()-36) > 1e-9 {
		t.Fatalf("light want 36, got %v", p.Light())
	}
}

func TestFeedBalanced(t *testing.T) {
	f := Feed{Rate: 100, Z: 0.5, Q: 1}
	d := Product{Rate: 50, X: 0.95}
	b := Product{Rate: 50, X: 0.05}
	if !FeedBalanced(f, d, b) {
		t.Fatal("mass balance should close")
	}
	bad := Product{Rate: 60, X: 0.9}
	if FeedBalanced(f, bad, b) {
		t.Fatal("unbalanced products should fail")
	}
}

func TestSplitFraction(t *testing.T) {
	f := SplitFraction(0.5, 0.95, 0.05)
	if math.Abs(f-0.5) > 1e-9 {
		t.Fatalf("split want 0.5, got %v", f)
	}
}

func TestRecovery(t *testing.T) {
	d := Product{Rate: 50, X: 0.95}
	f := Product{Rate: 100, X: 0.5}
	r := Recovery(d, f)
	if math.Abs(r-0.95) > 1e-9 {
		t.Fatalf("recovery want 0.95, got %v", r)
	}
}

func TestTotalMaterial(t *testing.T) {
	if TotalMaterial(100, Product{Rate: 60, X: 0.9}, Product{Rate: 40, X: 0.05}) != 0 {
		t.Fatal("total material should be zero")
	}
}
