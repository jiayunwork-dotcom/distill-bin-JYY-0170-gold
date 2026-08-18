package vle

type Balance struct {
	F float64
	Z float64
	D float64
	X float64
	B float64
	W float64
}

type BalanceResult struct {
	Distillate float64
	Bottoms    float64
	Residual   float64
	Closed     bool
}

func SolveBalance(F, zF, xD, xB float64) (BalanceResult, error) {
	if F <= 0 {
		return BalanceResult{}, ErrInvalidFeed
	}
	if !inRange(zF) || !inRange(xD) || !inRange(xB) {
		return BalanceResult{}, ErrCompositionOutOfRange
	}
	if !(xD > zF && zF > xB) {
		return BalanceResult{}, ErrInvalidFeed
	}
	den := xD - xB
	if den == 0 {
		return BalanceResult{}, ErrInvalidFeed
	}
	D := F * (zF - xB) / den
	B := F - D
	if D < 0 || B < 0 {
		return BalanceResult{}, ErrInvalidFeed
	}
	check := D*xD + B*xB - F*zF
	res := BalanceResult{
		Distillate: D,
		Bottoms:    B,
		Residual:   math_abs(check),
		Closed:     math_abs(check) < 1e-9*math_abs(F*zF),
	}
	return res, nil
}

func inRange(c float64) bool {
	return c > 0 && c < 1
}

func math_abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func Recoveries(F, zF, xD, xB float64) (lightRec, heavyRec float64) {
	b, err := SolveBalance(F, zF, xD, xB)
	if err != nil {
		return 0, 0
	}
	if F*zF == 0 {
		return 0, 0
	}
	lightRec = b.Distillate * xD / (F * zF)
	heavyRec = b.Bottoms * (1 - xB) / (F * (1 - zF))
	return
}

func CheckOrder(zF, xD, xB float64) bool {
	return xD > zF && zF > xB
}
