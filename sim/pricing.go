package sim

import (
	"context"
	"fmt"
	"time"
)

type PricingRequest struct {
	trade Trade
}

type PricingResponse struct {
	trade Trade
	price Matrix
}

type Pricer struct {
	market Market
	in     chan PricingRequest
	out    chan<- PricingResponse
}

func (pricer *Pricer) Init(workers int) {
	for i := 0; i < workers; i++ {
		go pricer.newPricingEngine(i)
	}
}

func (pricer *Pricer) newPricingEngine(name int) {
	var priced = 0
	for req := range pricer.in {
		//fmt.Printf("Pricing %v with %v\n", trd.Id, name)
		price, _ := pricer.market.Price(req.trade)
		priced = priced + 1
		if priced%100 == 0 {
			fmt.Printf("Pricer %v has done %v: %v\n", name, priced, time.Now())
		}
		pricer.out <- PricingResponse{req.trade, price.Matrix}
	}
}

func (pricer *Pricer) Price(ctx context.Context, trd Trade) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("Cancelled")
	default:
		pricer.in <- PricingRequest{trd}
		return nil
	}
}
