package vle

import "math"

type EquilibriumRow struct {
	X    float64
	Y    float64
	Alpha float64
}

func EquilibriumTable(v VLE, from, to, steps int) []EquilibriumRow {
	if steps < 2 {
		steps = 2
	}
	rows := []EquilibriumRow{}
	for i := 0; i <= steps; i++ {
		x := float64(from) + float64(i)*(float64(to-from))/float64(steps)
		y := v.Y(x)
		a := 0.0
		if x > 0 && x < 1 && y > 0 && y < 1 {
			a = y * (1 - x) / (x * (1 - y))
		}
		rows = append(rows, EquilibriumRow{X: x, Y: y, Alpha: a})
	}
	return rows
}

func CompositionRange(v VLE, targetY float64) (float64, float64) {
	lo := 0.0
	hi := 1.0
	for i := 0; i < 200; i++ {
		mid := (lo + hi) / 2
		if v.Y(mid) < targetY {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo, hi
}

func PinchComposition(v VLE, q, zF float64) float64 {
	// pinch point is where q-line meets equilibrium curve
	lo := 0.0
	hi := zF
	if q > 1 {
		hi = zF
		lo = 0
	}
	f := func(x float64) float64 {
		yLine := QLineY(q, zF, x)
		return yLine - v.Y(x)
	}
	for i := 0; i < 200; i++ {
		mid := (lo + hi) / 2
		if f(mid) > 0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}

func QLineY(q, zF, x float64) float64 {
	if q == 1 {
		return zF
	}
	return q/(q-1)*x - zF/(q-1)
}

func VaporPressureAntoine(t, a, b, c float64) float64 {
	// Antoine: log10(P) = a - b/(t + c)
	return math.Pow(10, a-b/(t+c))
}

func RelativeVolatilityFromVP(pA, pB float64) float64 {
	if pB <= 0 {
		return 0
	}
	return pA / pB
}

func BubblePoint(alpha, x float64) float64 {
	if alpha <= 1 || x <= 0 || x >= 1 {
		return 0
	}
	return alpha * x / (1 + (alpha-1)*x)
}

func DewPoint(alpha, y float64) float64 {
	if alpha <= 1 || y <= 0 || y >= 1 {
		return 0
	}
	return y / (alpha - (alpha-1)*y)
}
