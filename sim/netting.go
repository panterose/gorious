package sim

import (
	"fmt"
	"sync"
)

type Netting struct {
	Name   string
	Trades map[int]Trade
}

func (n Netting) Size() int {
	return len(n.Trades)
}


type NettingFlow struct {
	netting Netting
	channel chan TradeSimulation
}

func NettingRouter(nettings []Netting, prices chan TradeSimulation) (chan float32, error) {
	// create the mapping between trade and netting
	// also create a channel and goroutine for each netting
	nbNettings := len(nettings)
	var allTradesForNettings sync.WaitGroup
	allTradesForNettings.Add(nbNettings)

	mapping := make(map[int]NettingFlow)
	results := make(chan float32);
	flows := make([]NettingFlow, nbNettings)

	for ind, netting := range nettings {
		//create the flow object
		flow := NettingFlow{netting, make(chan TradeSimulation)}
		flows[ind] = flow

		//add trade/channel mapping
		for _, trade := range netting.Trades {
			mapping[trade.Id] = flow
		}

		//create netting goroutine
		go NettingAggregation(results, flow, &allTradesForNettings)
	}

	go func() {
		allTradesForNettings.Wait()
		fmt.Println("All netting done, closing results channel")
		close(results)
	}()

	// range around the prices channel and reroute the simulation data to the appropiate channel
	go func() {
		for ts := range prices {
			flow := mapping[ts.Id]
			//fmt.Printf("received price %v , mapping to flow %v\n", ts.Id, fl.netting.Name)
			cha := flow.channel
			if (cha == nil) {
				fmt.Printf("channel nil for %v on %v\n",ts.Id, flow.netting.Name)
			}
			cha <- ts
		}
		fmt.Println("All prices done, closing sub channels")

		//close all netting channel
		for _, flow := range flows {
			close(flow.channel)
		}
	}()

	return results, nil

}

func NettingAggregation(results chan float32, flow NettingFlow, allTradesForNettings *sync.WaitGroup) {
	var total = NewMatrix(1000, 0)
	for ts := range flow.channel {
		//fmt.Printf("aggregating price %v for %v\n", ts.Id, flow.netting.Name)
		total.Add(ts.Matrix)
	}

	fmt.Printf("Finished aggregating %v: total = %v\n", flow.netting.Name, total.slice[0])
	results <- total.slice[0]
	allTradesForNettings.Done()
}
