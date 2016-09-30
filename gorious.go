package main

import (
	"context"
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

		//fmt.Printf("Netting %v = %v  : %v \n", name, nettings[n], time.Now())
	}

	//setup simulation and aggregation
	market := sim.Market{Matrix: sim.NewRandomMatrix(rows, cols, 0)}

	trades := make(chan sim.PricingRequest)
	prices := make(chan sim.PricingResponse)
	results := make(chan float32)

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	pricer := sim.Pricer{Market: market, In: trades, Out: prices}
	netter := sim.NettingGroup{Nettings: make(map[string]*sim.NettingEngine), Prices: prices, Results: results}
	pricer.Init(ctx, 10, tradePerNetting)
	netter.Init(ctx, nettings, 10, tradePerNetting)
	start := time.Now()

	//send some trades to be prices
	for i := 0; i < nbTrades; i++ {
		trades <- sim.PricingRequest{Trade: sim.NewTrade(i)}
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
