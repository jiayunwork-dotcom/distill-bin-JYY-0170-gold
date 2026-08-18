package cli

import (
	"bytes"
	"strings"
	"testing"

	"distill-bin/internal/vle"
)

func vleBalanceClosed() vle.BalanceResult {
	res, _ := vle.SolveBalance(100, 0.5, 0.95, 0.05)
	return res
}

func TestRunGilliland(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := RunGilliland([]string{"-f", "../../example/benzene.json"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "r_min") || !strings.Contains(out.String(), "n_min") {
		t.Fatalf("missing shortcut output:\n%s", out.String())
	}
}

func TestRunStage(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := RunStage([]string{"-f", "../../example/benzene.json"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "tray") {
		t.Fatalf("missing tray profile:\n%s", out.String())
	}
}

func TestBalanceErrorText(t *testing.T) {
	if balanceErrorText(vleBalanceClosed()) != "closed" {
		t.Fatal("closed balance should report closed")
	}
}
