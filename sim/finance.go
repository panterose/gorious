package sim

import "math"

const (
	call string = "c"
	put  string = "p"
)

func unDef(f float64) bool {
	if math.IsNaN(f) {
		return true
	}
	if math.IsInf(f, 1) {
		return true
	}
	if math.IsInf(f, -1) {
		return true
	}
	return false
}

// Black Scholes (European put and call options)
// C Theoretical call premium (non-dividend paying stock)
// c = sn(d1) - ke^(-rt)N(d2)
// d1 = ln(s/k) + (r + v^2/2)t
// d2 = d1- vt^1/2
// k = Stock strike price
// s = Spot price
// t = time to expire in years
// r = risk free rate
// v = volitilaty (sigma)
// e math.E 2.7183
// putcall = "c" for a call or "p" for a put
func BlackScholes(s, k, t, r, v float64, putcall string) float64 {

	if k == 0 {
		return 0
	}
	if t <= 0 {
		return math.NaN()
	}
	if v == 0 {
		return math.NaN()
	}
	log := math.Log(s / k)
	if unDef(log) {
		return math.NaN()
	}
	vsq := math.Pow(v, 2)
	if unDef(vsq) {
		return math.NaN()
	}
	tsqr := math.Sqrt(t)
	if unDef(tsqr) {
		return math.NaN()
	}
	d1 := (log + ((r + (vsq / 2)) * t)) / (v * tsqr)
	d2 := d1 - (v * tsqr)
	emrt := math.Pow(math.E, (-1 * r * t))
	if unDef(emrt) {
		return math.NaN()
	}
	kemrt := k * emrt
	if putcall == call {
		cdfd1 := math.Erf(d1)
		if unDef(cdfd1) {
			return math.NaN()
		}

		cdfd2 := math.Erf(d2)
		if unDef(cdfd2) {
			return math.NaN()
		}
		return cdfd1*s - cdfd2*kemrt
	}
	if putcall == put {
		mcdfd2 := math.Erf(-1 * d2)
		if unDef(mcdfd2) {
			return math.NaN()
		}

		mcdfd1 := math.Erf(-1 * d1)
		if unDef(mcdfd1) {
			return math.NaN()
		}
		return kemrt - s + (mcdfd2*kemrt - mcdfd1*s)
	}
	return math.NaN()
}
