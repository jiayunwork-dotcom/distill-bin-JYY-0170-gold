package cli

import (
	"fmt"
	"io"

	"distill-bin/internal/shortcut"
	"distill-bin/internal/stage"
	"distill-bin/internal/vle"
)

type RunResult struct {
	Distillate float64
	Bottoms    float64
	RMin       float64
	NMin       float64
	R          float64
	Stages     int
	Feasible   bool
	Trays      []float64
	Closed     bool
	Residual   float64
}

func RunDesign(in Input) (RunResult, error) {
	s, err := in.ToShortcut()
	if err != nil {
		return RunResult{}, err
	}
	v, err := in.ToVLE()
	if err != nil {
		return RunResult{}, err
	}
	bal, err := vle.SolveBalance(in.Feed, in.FeedComp, in.Distillate, in.Bottoms)
	if err != nil {
		return RunResult{}, err
	}
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		return RunResult{}, err
	}
	nMin := s.FenskeMinimumStages()
	if in.Reflux < rMin {
		return RunResult{
			Distillate: bal.Distillate,
			Bottoms:    bal.Bottoms,
			RMin:       rMin,
			NMin:       nMin,
			R:          in.Reflux,
			Feasible:   false,
			Closed:     bal.Closed,
			Residual:   bal.Residual,
		}, nil
	}
	maxStages := in.MaxStages
	if maxStages <= 0 {
		maxStages = 1000
	}
	step, err := stage.StepOff(v, in.Reflux, in.Distillate, in.Bottoms, in.Q, in.FeedComp, maxStages)
	if err != nil {
		return RunResult{
			Distillate: bal.Distillate,
			Bottoms:    bal.Bottoms,
			RMin:       rMin,
			NMin:       nMin,
			R:          in.Reflux,
			Feasible:   false,
			Closed:     bal.Closed,
			Residual:   bal.Residual,
		}, nil
	}
	return RunResult{
		Distillate: bal.Distillate,
		Bottoms:    bal.Bottoms,
		RMin:       rMin,
		NMin:       nMin,
		R:          in.Reflux,
		Stages:     step.Stages,
		Feasible:   true,
		Trays:      step.TrayCompositions,
		Closed:     bal.Closed,
		Residual:   bal.Residual,
	}, nil
}

func WriteResult(w io.Writer, res RunResult) {
	fmt.Fprintf(w, "distillate      %.6f kmol/h\n", res.Distillate)
	fmt.Fprintf(w, "bottoms         %.6f kmol/h\n", res.Bottoms)
	fmt.Fprintf(w, "r_min           %.6f\n", res.RMin)
	fmt.Fprintf(w, "n_min (Fenske)  %.6f\n", res.NMin)
	fmt.Fprintf(w, "reflux          %.6f\n", res.R)
	if !res.Feasible {
		fmt.Fprintf(w, "stages          infeasible (reflux below minimum)\n")
	} else {
		fmt.Fprintf(w, "stages          %d (incl. reboiler)\n", res.Stages)
	}
	fmt.Fprintf(w, "balance_closed  %v\n", res.Closed)
	fmt.Fprintf(w, "balance_resid   %.3e\n", res.Residual)
	if res.Feasible && len(res.Trays) > 1 {
		fmt.Fprintf(w, "tray_liquid     [")
		for i, c := range res.Trays {
			if i > 0 {
				fmt.Fprintf(w, " ")
			}
			fmt.Fprintf(w, "%.4f", c)
		}
		fmt.Fprintf(w, "]\n")
	}
}

func GillilandStages(w io.Writer, s shortcut.Shortcut, r float64) {
	rMin, _ := s.UnderwoodMinReflux()
	nMin := s.FenskeMinimumStages()
	if r < rMin {
		fmt.Fprintf(w, "reflux %.3f below Rmin %.3f: infeasible\n", r, rMin)
		return
	}
	n := s.GillilandStages(r, rMin, nMin)
	fmt.Fprintf(w, "gilliland stages at R=%.3f: %.2f (Nmin=%.2f)\n", r, n, nMin)
}
