package cli

import (
	"fmt"
	"io"
)

func RunRate(args []string, stdout, stderr io.Writer) int {
	var path string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f", "--file":
			if i+1 < len(args) {
				i++
				path = args[i]
			}
		case "-h", "--help":
			fmt.Fprintln(stdout, usage)
			return 0
		default:
			if path == "" {
				path = args[i]
			}
		}
	}
	var in Input
	var err error
	if path != "" {
		in, err = ReadFile(path)
	} else {
		in, err = ReadStdin()
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	if err := in.Validate(); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	res, err := RunDesign(in)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	WriteResult(stdout, res)
	return 0
}

const usage = `usage: distill-bin rate [-f <file>]

binary distillation shortcut/stage rating from JSON on stdin or a file.
fields: feed, feed_composition, distillate_composition, bottoms_composition,
alpha, q, reflux. prints D, B, Rmin, Nmin, stage count and tray profile.`
