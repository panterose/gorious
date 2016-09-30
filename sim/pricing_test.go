package sim

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPricingEngine(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	trade1 := Trade{1, "id1", 1.0}
	mkt := Market{NewRandomMatrix(1, 2, 1)}
	in := make(chan PricingRequest)
	out := make(chan PricingResponse)

	pricer := Pricer{mkt, in, out}

	go pricer.newPricingWorker(ctx, 1, 1)()

	in <- PricingRequest{trade1}
	close(in)

	res := <-out
	close(out)

	assert.Equal(t, trade1, res.trade, "the price should be for trade1")
	assert.Equal(t, 1, res.price.rows, "The rows should be 1")
	assert.Equal(t, 2, res.price.cols, "The cols should be 2")
	assert.Equal(t, float32(604.6603), res.price.slice[0])
	assert.Equal(t, float32(940.5091), res.price.slice[1])
}

func TestInit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	mkt := Market{NewRandomMatrix(1, 2, 1)}
	in := make(chan PricingRequest)
	out := make(chan PricingResponse)

	pricer := Pricer{mkt, in, out}
	pricer.Init(ctx, 5, 1)

	go func() {
		for r := 0; r < 10; r++ {
			in <- PricingRequest{Trade{r, "id" + strconv.Itoa(r), float32(r) + 0.0}}
		}
		close(in)
	}()

	for item := range out {
		fmt.Printf("Price done %v: %v\n", item, time.Now())
	}
}
