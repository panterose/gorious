package sim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlackScholes(te *testing.T) {
	s := 1.0
	k := 2.0
	t := 1.0
	r := 0.02
	v := 0.1

	p := BlackScholes(s, k, t, r, v, "p")
	assert.Equal(te, 1.9207946932270215, p, "Value of the PUT option")

	c := BlackScholes(s, k, t, r, v, "c")
	assert.Equal(te, 0.9603973466135107, c, "Value of the CALL option")
}
