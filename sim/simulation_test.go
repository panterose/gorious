package sim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrice(t *testing.T) {
	mkt := &Market{NewMatrix(1, 2)}
	mkt.array[0] = 10.0
	mkt.array[1] = 20.0

	trd := Trade{1, "id1", 3.0}
	prc, _ := mkt.Price(trd)

	tr, tc := prc.Dims()
	assert.Equal(t, tr, 1, "1 row")
	assert.Equal(t, tc, 2, "2 cols")
	assert.Equal(t, float32(30.0), prc.array[0], "1st element should be 30")
	assert.Equal(t, float32(60.0), prc.array[1], "2nd element should be 60")
}
