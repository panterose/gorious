package sim

import (
	"fmt"
	"sync"
	"runtime"
)

//Market is the object representing some Market Data used for simulation
type Market struct {
	Matrix
}

func (mkt *Market) Price(trd Trade) (TradeSimulation, error) {
	price, _ := mkt.Mult(trd.Mtm)
	return TradeSimulation{trd, price}, nil
}



//Simulate a market and price all trades with given ids and pass them to a givn channel
func Simulate(market Market, trades chan Trade) (chan TradeSimulation, error) {
	nbPricers := runtime.NumCPU() * 20
	prices := make(chan TradeSimulation)
	var pricers sync.WaitGroup
	pricers.Add(nbPricers)
	for i := 0; i < nbPricers; i++ {
		//fmt.Println("ready to price", id)
		go func(name int) {
			for trd := range trades {
				//fmt.Printf("Pricing %v with %v\n", trd.Id, name)
				tradePrice, _ := market.Price(trd)
				prices <- tradePrice
			}
			fmt.Printf("Pricer %v is done\n", name)
			pricers.Done()
		}(i)
	}

	go func() {
		pricers.Wait()
		fmt.Println("Closing prices channel")
		close(prices)
	}()
	return prices, nil
}
