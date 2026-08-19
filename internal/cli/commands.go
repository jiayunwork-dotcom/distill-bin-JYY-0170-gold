package cli

import (
	"fmt"
	"io"
	"math"

	"distill-bin/internal/stage"
	"distill-bin/internal/vle"
)

func RunGilliland(args []string, stdout, stderr io.Writer) int {
	var path string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f", "--file":
			if i+1 < len(args) {
				i++
				path = args[i]
			}
		case "-h", "--help":
			fmt.Fprintln(stdout, gillilandUsage)
			return 0
		default:
			if path == "" {
				path = args[i]
			}
		}
	}
	in, err := readInput(path)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	s, err := in.ToShortcut()
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	nMin := s.FenskeMinimumStages()
	fmt.Fprintf(stdout, "r_min      %.6f\n", rMin)
	fmt.Fprintf(stdout, "n_min      %.6f\n", nMin)
	for _, mult := range []float64{1.1, 1.2, 1.5, 2.0, 3.0} {
		r := rMin * mult
		n := s.GillilandStages(r, rMin, nMin)
		fmt.Fprintf(stdout, "R=%.3f (%.1fx Rmin)  N=%.2f\n", r, mult, n)
	}
	return 0
}

func RunStage(args []string, stdout, stderr io.Writer) int {
	var path string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f", "--file":
			if i+1 < len(args) {
				i++
				path = args[i]
			}
		case "-h", "--help":
			fmt.Fprintln(stdout, stageUsage)
			return 0
		default:
			if path == "" {
				path = args[i]
			}
		}
	}
	in, err := readInput(path)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	v, err := in.ToVLE()
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	profile, err := stage.CompositionProfile(v, in.Reflux, in.Distillate, in.Bottoms, in.Q, in.FeedComp, 1000)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "feed_tray    %d\n", profile.FeedTray)
	fmt.Fprintf(stdout, "rectifying   %d\n", stage.RectifyingTrays(profile))
	fmt.Fprintf(stdout, "stripping    %d\n", stage.StrippingTrays(profile))
	for _, tray := range profile.Trays {
		fmt.Fprintf(stdout, "tray %2d  x=%.4f  y=%.4f  %s\n", tray.Number, tray.X, tray.Y, tray.Section)
	}
	return 0
}

func readInput(path string) (Input, error) {
	if path != "" {
		return ReadFile(path)
	}
	return ReadStdin()
}

func balanceErrorText(res vle.BalanceResult) string {
	if res.Closed {
		return "closed"
	}
	return fmt.Sprintf("residual %.3e", res.Residual)
}

func stagesOrInf(n int, err error) float64 {
	if err != nil {
		return math.Inf(1)
	}
	return float64(n)
}

const gillilandUsage = `usage: distill-bin gilliland [-f <file>]

print Fenske Nmin, Underwood Rmin and Gilliland stage estimates at several
reflux ratios.`

const stageUsage = `usage: distill-bin stage [-f <file>]

McCabe-Thiele step-off: print the tray-by-tray liquid/vapor composition profile.`
