package shortcut

import (
	"math"
)

type RminStudy struct {
	Q        float64
	RMin     float64
	Theta    float64
}

func RminAtFeedQuality(s Shortcut, q float64) (RminStudy, error) {
	work := s
	work.Q = q
	theta, err := work.UnderwoodTheta()
	if err != nil {
		return RminStudy{}, err
	}
	rMin := work.Alpha*work.XD/(work.Alpha-theta) - 1
	return RminStudy{Q: q, RMin: rMin, Theta: theta}, nil
}

func RminSensitivity(s Shortcut) (float64, error) {
	if _, err := s.UnderwoodMinReflux(); err != nil {
		return 0, err
	}
	hi, err := RminAtFeedQuality(s, s.Q+0.1)
	if err != nil {
		return 0, err
	}
	lo, err := RminAtFeedQuality(s, s.Q-0.1)
	if err != nil {
		return 0, err
	}
	return (hi.RMin - lo.RMin) / 0.2, nil
}

func RminVsAlpha(alpha, zF, xD, xB, q float64) (float64, error) {
	s, err := New(alpha, zF, xD, xB, q)
	if err != nil {
		return 0, err
	}
	return s.UnderwoodMinReflux()
}

func NminVsSeparation(alpha, xD, xB float64) float64 {
	if alpha <= 1 {
		return math.Inf(1)
	}
	num := math.Log((xD/(1-xD))*((1-xB)/xB))
	return num / math.Log(alpha)
}

func CheckDesignPair(s Shortcut, r, n float64) bool {
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		return false
	}
	nMin := s.FenskeMinimumStages()
	computed := s.GillilandStages(r, rMin, nMin)
	return math.Abs(computed-n) <= 2.5
}

func RefluxRange(s Shortcut, maxR float64) (minR, feasibleMax float64, err error) {
	rMin, err := s.UnderwoodMinReflux()
	if err != nil {
		return 0, 0, err
	}
	return rMin, maxR, nil
}
