package sim

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	trade11  = Trade{1, "id1", 1.0}
	trade12  = Trade{2, "id2", 1.0}
	trades1  = map[int]Trade{1: trade11, 2: trade12}
	netting1 = Netting{"n1", trades1}

	trade21  = Trade{3, "id3", 1.0}
	trade22  = Trade{4, "id4", 1.0}
	trades2  = map[int]Trade{3: trade21, 4: trade22}
	netting2 = Netting{"n2", trades2}
)

func TestNettingEngineWorker(t *testing.T) {
	mat := NewMatrix(1, 2)
	in := make(chan NettingRequest)
	out := make(chan float32)

	ne := NettingEngine{netting1, mat, in, out}
	price := NewRandomMatrix(1, 2, 1)
	ctx := context.Background()
	go ne.newNettingWorker(ctx)()

	in <- NettingRequest{trade11, price}
	close(in)

	result := <-out
	close(out)

	assert.Equal(t, float32(604.6603), ne.Result.slice[0])
	assert.Equal(t, float32(940.5091), ne.Result.slice[1])
	assert.Equal(t, float32(604.6603), result)
}

func TestNettingGroup(t *testing.T) {

	fmt.Printf("Starting TestNettingGroup: %v\n", time.Now())
	in := make(chan PricingResponse)
	out := make(chan float32)
	nettingMap := make(map[string]*NettingEngine)
	ng := NettingGroup{nettingMap, in, out}
	ctx := context.Background()
	ng.Init(ctx, []Netting{netting1, netting2}, 1)

	ne1 := ng.nettings["n1"]
	ne2 := ng.nettings["n2"]
	assert.Equal(t, netting1, ne1.netting)
	assert.Equal(t, trade11, ne1.netting.Trades[1])
	assert.Equal(t, trade12, ne1.netting.Trades[2])
	assert.Equal(t, netting2, ne2.netting)
	assert.Equal(t, trade21, ne2.netting.Trades[3])
	assert.Equal(t, trade22, ne2.netting.Trades[4])

	go func() {
		in <- PricingResponse{trade11, NewRandomMatrix(1, 2, 1)}
		in <- PricingResponse{trade12, NewRandomMatrix(1, 2, 2)}
		in <- PricingResponse{trade21, NewRandomMatrix(1, 2, 3)}
		in <- PricingResponse{trade22, NewRandomMatrix(1, 2, 4)}

		close(in)
	}()

	<-out
	<-out

	assert.Equal(t, float32(771.9569), ne1.result())
	assert.Equal(t, float32(963.3544), ne2.result())
}
