package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/panterose/gorious/sim"
)

func main() {
	const nbTrades = 100000
	const nbNetting = 100
	const tradePerNetting = nbTrades / nbNetting

	//create the output channel
	results := make(chan float32)
	trades := make(chan sim.Trade)
	prices := make(chan *sim.TradeSimulation)

	nettings := make([]sim.Netting, nbNetting)
	for n := 0; n < nbNetting; n++ {
		trades := make(map[int]sim.Trade, tradePerNetting)
		for t := 0; t < tradePerNetting; t++ {
			trades[t] = sim.NewTrade(n + t*nbNetting)
		}
		nettings[n] = sim.Netting{"netting" + strconv.Itoa(n), trades}
	}

	//setup simulation and aggregation
	sim.Simulate(0, trades, prices)
	sim.NettingRouter(nettings, results, prices)

	start := time.Now()

	//send some trades to be prices
	for i := 0; i < nbTrades; i++ {
		trades <- sim.NewTrade(i)
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
