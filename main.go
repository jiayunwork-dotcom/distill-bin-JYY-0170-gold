package main

import (
	"fmt"
	"os"

	"distill-bin/internal/cli"
)

const usage = `distill-bin — binary distillation shortcut & McCabe-Thiele rating

usage:
  distill-bin rate [flags] [file]
    shortcut rating + stage count for a binary column
  distill-bin gilliland [flags] [file]
    Fenske/Underwood/Gilliland stage estimates at several refluxes
  distill-bin stage [flags] [file]
    McCabe-Thiele tray-by-tray composition profile

flags:
  -f, --file <path>   read JSON from a file (default: stdin)
  -h, --help          show this help

input JSON:
  { "feed": 100, "feed_composition": 0.5,
    "distillate_composition": 0.95, "bottoms_composition": 0.05,
    "alpha": 2.5, "q": 1.0, "reflux": 2.0 }

prints D, B, Rmin, Fenske Nmin, stages at the given reflux and tray profile.

examples:
  distill-bin rate example/benzene.json
  distill-bin gilliland example/benzene.json
  cat example/benzene.json | distill-bin stage`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]
	switch cmd {
	case "rate":
		os.Exit(cli.RunRate(args, os.Stdout, os.Stderr))
	case "gilliland":
		os.Exit(cli.RunGilliland(args, os.Stdout, os.Stderr))
	case "stage":
		os.Exit(cli.RunStage(args, os.Stdout, os.Stderr))
	case "-h", "--help", "help":
		fmt.Println(usage)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s\n", cmd, usage)
		os.Exit(2)
	}
}
