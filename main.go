package main

import (
	"fmt"
	"strconv"
	"time"
	"github.com/panterose/gorious/sim"
	"runtime"
)

func main() {
	const rows = 1000
	const cols = 100

	nbNetting := runtime.NumCPU() * 10
	nbTrades := nbNetting * 1000
	tradePerNetting := nbTrades / nbNetting

	//create the output channel
	trades := make(chan sim.Trade)
	nettings := make([]sim.Netting, nbNetting)
	for n := 0; n < nbNetting; n++ {
		trades := make(map[int]sim.Trade, tradePerNetting)
		for t := 0; t < tradePerNetting; t++ {
			trades[t] = sim.NewTrade(n + t*nbNetting)
		}
		nettings[n] = sim.Netting{"netting" + strconv.Itoa(n), trades}
	}

	//setup simulation and aggregation
	market := sim.Market{sim.NewRandomMatrix(rows, cols, 0)}
	prices, _ := sim.Simulate(market, trades)
	results, _ := sim.NettingRouter(nettings, prices)

	start := time.Now()

	//send some trades to be prices
	for i := 0; i < nbTrades; i++ {
		trades <- sim.NewTrade(i)
	}
	fmt.Printf("Finished sending trade request, closing the trade channel \n")
	close(trades)

	//collect result from channel
	for total := range results {
		fmt.Printf("Total is %v \n", total)
	}

	elapsed := time.Since(start)
	fmt.Printf("Simulation/Aggregation took %s or %v ns/op  \n", elapsed, elapsed.Nanoseconds() / int64(nbTrades) )

}
