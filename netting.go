package sim

import "fmt"

type Netting struct {
	Name   string
	Trades map[int]Trade
}

func (n Netting) Size() int {
	return len(n.Trades)
}

type NettingFlow struct {
	netting Netting
	channel chan *TradeSimulation
}

func NettingRouter(nettings []Netting, results chan float32, prices chan *TradeSimulation) {
	// create the mapping between trade and netting
	// also create a channel and goroutine for each netting
	mapping := make(map[int]NettingFlow)
	flows := make([]NettingFlow, len(nettings))
	for ind, netting := range nettings {
		//create the flow object
		flow := NettingFlow{netting, make(chan *TradeSimulation)}
		flows[ind] = flow

		//add trade/channel mapping
		for _, trade := range netting.Trades {
			mapping[trade.Id] = flow
		}

		//create netting goroutine
		go NettingAggregation(results, flow)
	}

	// range around the prices channel and reroute the simulation data to the appropiate channel
	go func() {
		for ts := range prices {
			fl := mapping[ts.Id]
			//fmt.Printf("received price %v , mapping to flow %v\n", ts.Id, fl.netting.Name)
			cha := fl.channel
			cha <- ts
		}
		fmt.Println("All prices done, closing sub channels")

		//close all netting channel
		for _, flow := range flows {
			close(flow.channel)
		}
	}()

}

func NettingAggregation(results chan float32, flow NettingFlow) {
	var total = NewMatrix(1000, 0)
	for ts := range flow.channel {
		//fmt.Printf("aggregating price %v for %v\n", ts.Id, flow.netting.Name)
		total.Add(&ts.Matrix)
	}

	fmt.Printf("Finished aggregating %v: total = %v\n", flow.netting.Name, total.array[0])
	results <- total.array[0]
}
