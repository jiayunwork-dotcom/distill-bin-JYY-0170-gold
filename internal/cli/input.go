package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"distill-bin/internal/shortcut"
	"distill-bin/internal/vle"
)

type Input struct {
	Feed     float64  `json:"feed"`
	FeedComp float64  `json:"feed_composition"`
	Distillate float64 `json:"distillate_composition"`
	Bottoms  float64  `json:"bottoms_composition"`
	Alpha    float64  `json:"alpha"`
	Q        float64  `json:"q"`
	Reflux   float64  `json:"reflux"`
	MaxStages int     `json:"max_stages,omitempty"`
}

var (
	ErrEmptyInput   = errors.New("cli: empty input")
	ErrInvalidJSON  = errors.New("cli: invalid json")
	ErrMissingFeed  = errors.New("cli: feed rate missing")
)

var stdinOverride io.Reader

func ParseInput(r io.Reader) (Input, error) {
	var in Input
	dec := json.NewDecoder(r)
	if err := dec.Decode(&in); err != nil {
		if err == io.EOF {
			return in, ErrEmptyInput
		}
		return in, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	return in, nil
}

func ReadFile(path string) (Input, error) {
	f, err := os.Open(path)
	if err != nil {
		return Input{}, err
	}
	defer f.Close()
	return ParseInput(f)
}

func ReadStdin() (Input, error) {
	if stdinOverride != nil {
		return ParseInput(stdinOverride)
	}
	return ParseInput(os.Stdin)
}

func (in Input) Validate() error {
	if in.Feed <= 0 {
		return ErrMissingFeed
	}
	return nil
}

func (in Input) ToShortcut() (shortcut.Shortcut, error) {
	return shortcut.New(in.Alpha, in.FeedComp, in.Distillate, in.Bottoms, in.Q)
}

func (in Input) ToVLE() (vle.VLE, error) {
	return vle.NewConstantAlpha(in.Alpha)
}
