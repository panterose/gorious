package sim

import (
	"context"
	"fmt"
	"time"

	ergp "golang.org/x/sync/errgroup"
)

type routine func() error

type PricingRequest struct {
	Trade Trade
}

type PricingResponse struct {
	trade Trade
	price Matrix
}

type Pricer struct {
	Market Market
	In     chan PricingRequest
	Out    chan<- PricingResponse
}

func (pricer *Pricer) Init(parent context.Context, workers int, modulo int) {
	g, ctx := ergp.WithContext(parent)
	for i := 0; i < workers; i++ {
		g.Go(pricer.newPricingWorker(ctx, i, modulo))
	}

	go func() {
		defer close(pricer.Out)
		error := g.Wait()
		fmt.Printf("Pricing finished  %v : %v\n", error, time.Now())
	}()
}

func (pricer *Pricer) newPricingWorker(ctx context.Context, name int, modulo int) routine {
	return func() error {
		var priced = 0
		fmt.Printf("Starting pricing worker %v \n", name)
		for req := range pricer.In {
			price, err := pricer.Market.Price(req.Trade)
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return fmt.Errorf("Aborting pricing %v", name)
			case pricer.Out <- PricingResponse{req.Trade, price.Matrix}:
				priced = priced + 1
				if priced%modulo == 0 {
					fmt.Printf("Pricer %v has done %v: %v\n", name, priced, time.Now())
				}
			}
		}
		return nil
	}
}
