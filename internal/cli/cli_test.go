package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunRateExample(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := RunRate([]string{"../../example/benzene.json"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errBuf.String())
	}
	for _, want := range []string{"distillate", "bottoms", "r_min", "n_min", "stages"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("output missing %q: %s", want, out.String())
		}
	}
}

func TestRunRateStdin(t *testing.T) {
	input := `{"feed":100,"feed_composition":0.5,"distillate_composition":0.95,
	          "bottoms_composition":0.05,"alpha":2.5,"q":1.0,"reflux":2.5}`
	old := stdinOverride
	stdinOverride = strings.NewReader(input)
	defer func() { stdinOverride = old }()

	var out, errBuf bytes.Buffer
	code := RunRate([]string{}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "balance_closed  true") {
		t.Fatalf("balance should close:\n%s", out.String())
	}
}

func TestRunRateInsufficientReflux(t *testing.T) {
	input := `{"feed":100,"feed_composition":0.5,"distillate_composition":0.95,
	          "bottoms_composition":0.05,"alpha":2.5,"q":1.0,"reflux":0.1}`
	old := stdinOverride
	stdinOverride = strings.NewReader(input)
	defer func() { stdinOverride = old }()

	var out, errBuf bytes.Buffer
	code := RunRate([]string{}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "infeasible") {
		t.Fatalf("expected infeasible output:\n%s", out.String())
	}
}

func TestRunRateRejectsBadComposition(t *testing.T) {
	input := `{"feed":100,"feed_composition":0.5,"distillate_composition":0.3,
	          "bottoms_composition":0.05,"alpha":2.5,"q":1.0,"reflux":2.0}`
	old := stdinOverride
	stdinOverride = strings.NewReader(input)
	defer func() { stdinOverride = old }()

	var out, errBuf bytes.Buffer
	code := RunRate([]string{}, &out, &errBuf)
	if code == 0 {
		t.Fatal("xD < zF should exit non-zero")
	}
}

func TestRunRateRejectsMissingFeed(t *testing.T) {
	input := `{"feed_composition":0.5,"distillate_composition":0.95,
	          "bottoms_composition":0.05,"alpha":2.5,"q":1.0,"reflux":2.0}`
	old := stdinOverride
	stdinOverride = strings.NewReader(input)
	defer func() { stdinOverride = old }()

	var out, errBuf bytes.Buffer
	code := RunRate([]string{}, &out, &errBuf)
	if code == 0 {
		t.Fatal("missing feed should exit non-zero")
	}
}

func TestRunRateRejectsAlphaOne(t *testing.T) {
	input := `{"feed":100,"feed_composition":0.5,"distillate_composition":0.95,
	          "bottoms_composition":0.05,"alpha":1.0,"q":1.0,"reflux":2.0}`
	old := stdinOverride
	stdinOverride = strings.NewReader(input)
	defer func() { stdinOverride = old }()

	var out, errBuf bytes.Buffer
	code := RunRate([]string{}, &out, &errBuf)
	if code == 0 {
		t.Fatal("alpha=1 should exit non-zero")
	}
}

func TestRunRateUnknownFile(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := RunRate([]string{"-f", "../../example/nope.json"}, &out, &errBuf)
	if code == 0 {
		t.Fatal("missing file should exit non-zero")
	}
}
