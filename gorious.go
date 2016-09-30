package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/panterose/gorious/sim"
)

func main() {
	const rows = 1000
	const cols = 100

	nbNetting := 1000
	nbTrades := nbNetting * 1000
	tradePerNetting := nbTrades / nbNetting

	//create the output channel
	nettings := make([]sim.Netting, nbNetting)
	for n := 0; n < nbNetting; n++ {
		trades := make(map[int]sim.Trade, tradePerNetting)
		for t := 0; t < tradePerNetting; t++ {
			trades[t] = sim.NewTrade(n + t*nbNetting)
		}
		name := "netting" + strconv.Itoa(n)
		nettings[n] = sim.Netting{Name: name, Trades: trades}
	}

	//setup simulation and aggregation
	market := sim.Market{sim.NewRandomMatrix(rows, cols, 0)}

	trades := make(chan sim.PricingRequest)
	prices := make(chan sim.PricingResponse)
	results := make(chan float32)

	start := time.Now()

	//send some trades to be prices
	for i := 0; i < nbTrades; i++ {
		trades <- sim.NewTrade(i)
	}
	fmt.Printf("Finished sending trade request, closing the trade channel : %v \n", time.Now())
	close(trades)

	//collect result from channel
	for total := range results {
		fmt.Printf("Total is %v : %v\n", total, time.Now())
	}

	elapsed := time.Since(start)
	fmt.Printf("Simulation/Aggregation took %s or %v ns/op : %v\n", elapsed, elapsed.Nanoseconds()/int64(nbTrades), time.Now())

}
