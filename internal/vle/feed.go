package vle

import "math"

type Feed struct {
	Rate float64
	Z    float64
	Q    float64
}

func (f Feed) Vapor() float64 {
	return f.Rate * (1 - f.Q)
}

func (f Feed) Liquid() float64 {
	return f.Rate * f.Q
}

func (f Feed) LightInFeed() float64 {
	return f.Rate * f.Z
}

func (f Feed) HeavyInFeed() float64 {
	return f.Rate * (1 - f.Z)
}

type Product struct {
	Rate float64
	X    float64
}

func (p Product) Light() float64 {
	return p.Rate * p.X
}

func (p Product) Heavy() float64 {
	return p.Rate * (1 - p.X)
}

func FeedBalanced(f Feed, d, b Product) bool {
	if f.Rate <= 0 {
		return false
	}
	massOK := math.Abs(f.Rate-(d.Rate+b.Rate)) < 1e-9*f.Rate
	lightOK := math.Abs(f.LightInFeed()-(d.Light()+b.Light())) < 1e-9*f.LightInFeed()
	return massOK && lightOK
}

func SplitFraction(zF, xD, xB float64) float64 {
	den := xD - xB
	if den == 0 {
		return 0
	}
	return (zF - xB) / den
}

func Recovery(d, f Product) float64 {
	if f.Light() <= 0 {
		return 0
	}
	return d.Light() / f.Light()
}

func TotalMaterial(f float64, d, b Product) float64 {
	return d.Rate + b.Rate - f
}

func LightKeyBalance(zF, xD, xB, recovery float64) float64 {
	return recovery * zF
}

func HeavyKeyBalance(zF, xD, xB float64) float64 {
	return (1 - zF)
}
