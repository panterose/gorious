package sim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNettingAggregation(t *testing.T) {
	trades := map[int]Trade{
		1: Trade{1, "id1", 1.0},
		2: Trade{2, "id2", 2.0},
	}
	netting := Netting{"n", trades}
	flow := NettingFlow{netting, make(chan *TradeSimulation)}
	result := make(chan float32)

	go NettingAggregation(result, flow)
	flow.channel <- &TradeSimulation{trades[1], NewRandomMatrix(1, 2, 0)}
	flow.channel <- &TradeSimulation{trades[2], NewRandomMatrix(1, 2, 1)}

	close(flow.channel)

	total := <-result
	assert.Equal(t, float32(1549.8564), total)
}
