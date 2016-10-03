package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/panterose/gorious/sim"
)

func main() {
	const rows = 1000
	const cols = 100

	nbNetting := 1000
	nbTrades := nbNetting * 1000
	nbworkers := 10
	modulo := (nbTrades / nbworkers) / 10
	//modulo := 1
	//	tradePerNetting := nbTrades / nbNetting

	// create a random mapping
	r := rand.New(rand.NewSource(1))
	tradeMapping := make(map[int][]int, nbTrades)
	for t := 0; t < nbTrades; t++ {
		net := r.Intn(nbNetting)
		tradeMapping[net] = append(tradeMapping[net], t)
	}
	//fmt.Printf("tradeMapping %v \n", tradeMapping)

	// create the Netting based on the mapping
	nettings := make([]*sim.Netting, nbNetting)
	for n, tradeIds := range tradeMapping {
		trades := make(map[int]sim.Trade, len(tradeIds))
		for _, t := range tradeIds {
			trade := sim.NewTrade(t)
			trades[t] = trade
		}
		name := "netting" + strconv.Itoa(n)
		nettings[n] = sim.NewNetting(name, trades)

		fmt.Printf("Netting %v = %v \n", name, nettings[n].Size())
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
	pricer.Init(ctx, nbworkers, modulo)
	netter.Init(ctx, nettings, nbworkers, modulo)
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
