package stage

import (
	"math"

	"distill-bin/internal/vle"
)

type FeedQuality struct {
	Q      float64
	Type   string
	HVaporized float64
}

func FeedType(q float64) string {
	switch {
	case q < 0:
		return "superheated_vapor"
	case q == 0:
		return "saturated_vapor"
	case q < 1:
		return "partially_vaporized"
	case q == 1:
		return "saturated_liquid"
	default:
		return "subcooled_liquid"
	}
}

func VaporFraction(q float64) float64 {
	return 1 - q
}

func LiquidFraction(q float64) float64 {
	return q
}

func FeedEnthalpyChange(q float64, latentHeat float64) float64 {
	return (1 - q) * latentHeat
}

func QLineSlope(q float64) float64 {
	if q == 1 {
		return math.Inf(1)
	}
	return q / (q - 1)
}

func QLineIntercept(q, zF float64) float64 {
	return -zF / (q - 1)
}

func FeedQualityFromFractions(vaporFrac, liquidFrac float64) float64 {
	total := vaporFrac + liquidFrac
	if total <= 0 {
		return 1
	}
	return liquidFrac / total
}

func IntersectionX(rect, ql OperatingLine) float64 {
	if math.IsInf(ql.Slope, 1) {
		return ql.Intercept
	}
	if math.Abs(rect.Slope-ql.Slope) < 1e-15 {
		return 0
	}
	return (ql.Intercept - rect.Intercept) / (rect.Slope - ql.Slope)
}

func IntersectionY(rect, ql OperatingLine) float64 {
	x := IntersectionX(rect, ql)
	if math.IsInf(ql.Slope, 1) {
		return rect.Slope*x + rect.Intercept
	}
	return ql.Slope*x + ql.Intercept
}

func MinimumRefluxByIntersection(v vle.VLE, r, xD, zF, q float64) bool {
	rect := RectifyingLine(r, xD)
	ql := QLine(q, zF)
	xiX := IntersectionX(rect, ql)
	xiY := IntersectionY(rect, ql)
	return xiY <= v.Y(xiX)
}

func OptimumFeedTray(v vle.VLE, r, xD, xB, q, zF float64) int {
	profile, err := CompositionProfile(v, r, xD, xB, q, zF, 1000)
	if err != nil {
		return 0
	}
	return profile.FeedTray
}

func RectifyingRatio(r float64) float64 {
	return r / (r + 1)
}

func DistillateIntercept(xD, r float64) float64 {
	return xD / (r + 1)
}
