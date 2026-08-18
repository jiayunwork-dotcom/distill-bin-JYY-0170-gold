package shortcut

import (
	"errors"
	"math"
	"testing"
)

func TestFenskeMinimumStages(t *testing.T) {
	s, err := New(2.5, 0.45, 0.97, 0.03, 1.0)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	nMin := s.FenskeMinimumStages()
	num := math.Log((0.97/0.03)*(0.97/0.03))
	want := num / math.Log(2.5)
	if math.Abs(nMin-want) > 1e-9 {
		t.Fatalf("Nmin want %v, got %v", want, nMin)
	}
	if nMin <= 0 {
		t.Fatalf("Nmin should be positive, got %v", nMin)
	}
}

func TestUnderwoodMinReflux(t *testing.T) {
	s, err := New(2.5, 0.45, 0.97, 0.03, 1.0)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		t.Fatalf("UnderwoodMinReflux: %v", err)
	}
	if rMin <= 0 {
		t.Fatalf("Rmin should be positive, got %v", rMin)
	}
}

func TestGillilandCorrelation(t *testing.T) {
	s, _ := New(2.5, 0.45, 0.97, 0.03, 1.0)
	rMin, _ := s.UnderwoodMinReflux()
	nMin := s.FenskeMinimumStages()
	r := rMin * 2
	n := s.GillilandStages(r, rMin, nMin)
	if n <= nMin {
		t.Fatalf("Gilliland N=%.2f should exceed Nmin=%.2f", n, nMin)
	}
	if math.IsInf(n, 0) || math.IsNaN(n) {
		t.Fatalf("Gilliland N invalid: %v", n)
	}
}

func TestStagesMonotoneInReflux(t *testing.T) {
	s, _ := New(2.5, 0.45, 0.97, 0.03, 1.0)
	rMin, _ := s.UnderwoodMinReflux()
	n1 := s.StagesAtReflux(rMin * 1.2)
	n2 := s.StagesAtReflux(rMin * 2.0)
	n3 := s.StagesAtReflux(rMin * 5.0)
	if !(n1 >= n2 && n2 >= n3) {
		t.Fatalf("stages must not rise with reflux: %.2f %.2f %.2f", n1, n2, n3)
	}
}

func TestRejectNonSeparating(t *testing.T) {
	if _, err := New(1.0, 0.45, 0.97, 0.03, 1.0); !errors.Is(err, ErrNonSeparating) {
		t.Fatalf("alpha=1 should fail, got %v", err)
	}
	if _, err := New(0.7, 0.45, 0.97, 0.03, 1.0); !errors.Is(err, ErrNonSeparating) {
		t.Fatalf("alpha=0.7 should fail, got %v", err)
	}
}

func TestRejectBadOrder(t *testing.T) {
	if _, err := New(2.5, 0.45, 0.30, 0.03, 1.0); !errors.Is(err, ErrInvalidComposition) {
		t.Fatalf("xD < zF should fail, got %v", err)
	}
}

func TestDesignAtRefluxInfeasible(t *testing.T) {
	s, _ := New(2.5, 0.45, 0.97, 0.03, 1.0)
	rMin, _ := s.UnderwoodMinReflux()
	d, err := DesignAtReflux(s, rMin*0.5)
	if err != nil {
		t.Fatalf("DesignAtReflux: %v", err)
	}
	if d.Feasible {
		t.Fatal("reflux below Rmin should be infeasible")
	}
}

func TestDesignFeasible(t *testing.T) {
	s, _ := New(2.5, 0.45, 0.97, 0.03, 1.0)
	rMin, _ := s.UnderwoodMinReflux()
	d, err := DesignAtReflux(s, rMin*1.5)
	if err != nil {
		t.Fatalf("DesignAtReflux: %v", err)
	}
	if !d.Feasible {
		t.Fatal("reflux above Rmin should be feasible")
	}
	if d.N < d.NMin {
		t.Fatalf("N=%.2f below Nmin=%.2f", d.N, d.NMin)
	}
}

func TestTotalStagesRounding(t *testing.T) {
	s, _ := New(2.5, 0.45, 0.97, 0.03, 1.0)
	rMin, _ := s.UnderwoodMinReflux()
	n, err := TotalStagesAt(s, rMin*2)
	if err != nil {
		t.Fatalf("TotalStagesAt: %v", err)
	}
	if n < 2 {
		t.Fatalf("stages too few: %d", n)
	}
}
