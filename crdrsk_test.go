package sim

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	const nbTrades = 1000000
	const nbNetting = 1000
	const tradePerNetting = nbTrades / nbNetting

	//create the output channel
	results := make(chan float32)
	trades := make(chan Trade)
	prices := make(chan *TradeSimulation)

	nettings := make([]Netting, nbNetting)
	for n := 0; n < nbNetting; n++ {
		trades := make(map[int]Trade, tradePerNetting)
		for t := 0; t < tradePerNetting; t++ {
			trades[t] = NewTrade(n + t*nbNetting)
		}
		nettings[n] = Netting{"netting" + strconv.Itoa(n), trades}
	}

	//setup simulation and aggregation
	Simulate(0, trades, prices)
	NettingRouter(nettings, results, prices)

	start := time.Now()

	//send some trades to be prices
	for i := 0; i < nbTrades; i++ {
		trades <- NewTrade(i)
	}
	close(trades)

	//collect result from channel
	for i := 0; i < nbNetting; i++ {
		total := <-results
		fmt.Printf("Total is %v \n", total)
	}

	elapsed := time.Since(start)
	fmt.Printf("Simulation/Aggregation took %s \n", elapsed)

}
