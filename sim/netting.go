package sim

import (
	"context"
	"fmt"
	"strconv"
	"time"

	ergp "golang.org/x/sync/errgroup"
)

type Netting struct {
	Name   string
	Trades map[int]Trade
}

func (n *Netting) Size() int {
	return len(n.Trades)
}

type NettingRequest struct {
	trade Trade
	price Matrix
}

type NettingGroup struct {
	Nettings map[string]*NettingEngine
	Prices   chan PricingResponse
	Results  chan float32
}

func (ng *NettingGroup) findNettingEngine(t Trade) (*NettingEngine, error) {
	id := t.Id
	for _, engine := range ng.Nettings {
		for _, trade := range engine.netting.Trades {
			if trade.Id == id {
				return engine, nil
			}
		}
	}
	var engine *NettingEngine
	return engine, fmt.Errorf("Can't find netting for %v", t)
}

func (ng *NettingGroup) findNettingEngine2(t Trade) (*NettingEngine, error) {
	id := t.Id
	n := id % len(ng.Nettings)
	if engine, ok := ng.Nettings["netting"+strconv.Itoa(n)]; ok {
		return engine, nil
	} else {
		var ne *NettingEngine
		return ne, fmt.Errorf("Couldn not find netting for: %v, tried %v", id, n)
	}
}

func (ng *NettingGroup) close() {
	fmt.Printf("Closing NettingGroup: %v \n", time.Now())
	for _, engine := range ng.Nettings {
		close(engine.in)
	}
}

// Init starts all the goroutine to ready the netting group to process message
func (ng *NettingGroup) Init(parent context.Context, nettings []Netting, nbrouter int, modulo int) {
	g1, ctx := ergp.WithContext(parent)
	for _, netting := range nettings {

		in := make(chan NettingRequest)
		total := NewMatrix(1000, 20)
		ne := &NettingEngine{netting: netting, Result: total, in: in, out: ng.Results}
		ng.Nettings[netting.Name] = ne
		g1.Go(ne.newNettingWorker(ctx, modulo))
	}
	go func() {
		err := g1.Wait()
		fmt.Printf("Netting is finished with %v: %v \n", err, time.Now())
		close(ng.Results)
	}()

	g2, ctx := ergp.WithContext(parent)
	for i := 0; i < nbrouter; i = i + 1 {
		g2.Go(ng.newNettingRouter(ctx, i, modulo))
	}
	go func() {
		err := g2.Wait()
		if err != nil {
			fmt.Printf("Routing is finished with %v: %v \n", err, time.Now())
		}
		ng.close()
	}()
}

func (ng *NettingGroup) newNettingRouter(ctx context.Context, name int, modulo int) routine {
	return func() error {
		var done = 0
		for price := range ng.Prices {
			engine, err := ng.findNettingEngine2(price.trade)
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return fmt.Errorf("Cancelled")
			case engine.in <- NettingRequest{price.trade, price.price}:
				done = done + 1
				if done%modulo == 0 {
					fmt.Printf("routed %v to %v: \n", price.trade.Id, engine.netting.Name)
				}
			}
		}
		return nil
	}
}

// NettingEngine store and process exposure for a specific Netting
type NettingEngine struct {
	netting Netting
	Result  Matrix
	in      chan NettingRequest
	out     chan float32
}

func (ne *NettingEngine) newNettingWorker(ctx context.Context, modulo int) routine {
	return func() error {
		var done = 0
		for nr := range ne.in {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				ne.aggregate(nr)
				done = done + 1
				if done%modulo == 0 {
					fmt.Printf("netting %v to %v: \n", nr.trade.Id, ne.netting.Name)
				}

			}
		}
		fmt.Printf("Aggregation for %v done, result=%v : %v \n", ne.netting.Name, ne.result(), time.Now())
		ne.out <- ne.result()
		return nil
	}
}

func (ne *NettingEngine) aggregate(nr NettingRequest) {
	//fmt.Printf("aggregate %v on %v: %v \n", nr.trade.Id, ne.netting.Name, time.Now())
	ne.Result.Add(nr.price)
}

func (ne *NettingEngine) result() float32 {
	value, _ := ne.Result.Get(0, 0)
	return value
}
