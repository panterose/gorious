package sim

import (
	"fmt"
	"sync"
)

//Market is the object representing some Market Data used for simulation
type Market struct {
	Matrix
}

func (mkt *Market) Price(trd Trade) (*TradeSimulation, error) {
	price, _ := mkt.Mult(trd.Mtm)
	return &TradeSimulation{trd, price}, nil
}

const rows = 1000
const cols = 100
const nbPricer = 10

//Simulate a market and price all trades with given ids and pass them to a givn channel
func Simulate(seed int64, trades chan Trade, prices chan *TradeSimulation) error {
	market := &Market{NewRandomMatrix(rows, cols, 0)}
	var wg sync.WaitGroup
	wg.Add(nbPricer)
	for i := 0; i < nbPricer; i++ {
		//fmt.Println("ready to price", id)
		go func(name int) {
			for trd := range trades {
				//fmt.Printf("Pricing %v with %v\n", trd.Id, name)
				tradePrice, _ := market.Price(trd)
				prices <- tradePrice
			}
			fmt.Printf("Pricer %v is done\n", name)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		fmt.Println("Closing prices channel")
		close(prices)
	}()
	return nil
}
