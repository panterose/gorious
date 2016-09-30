package sim

import (
	"context"
	"fmt"
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
	nettings map[string]*NettingEngine
	prices   chan PricingResponse
	results  chan float32
}

// CreateNettingEngines allows to create the netting/engine map needed to create a NettingGroup
func CreateNettingEngines(nettings []Netting, out chan float32) map[string]*NettingEngine {
	nettingMap := make(map[string]*NettingEngine)
	for _, netting := range nettings {

		in := make(chan NettingRequest)
		total := NewMatrix(1000, 20)
		ne := NettingEngine{netting, total, in, out}
		nettingMap[netting.Name] = &ne
	}
	return nettingMap
}

func (ng *NettingGroup) findNettingEngine(t Trade) (*NettingEngine, error) {
	id := t.Id
	for _, engine := range ng.nettings {
		for tradeID := range engine.netting.Trades {
			if tradeID == id {
				return engine, nil
			}
		}
	}
	var engine *NettingEngine
	return engine, fmt.Errorf("Can't find netting for %v", t)
}

func (ng *NettingGroup) close() {
	fmt.Printf("Closing NettingGroup: %v \n", time.Now())
	for _, engine := range ng.nettings {
		close(engine.in)
	}
}

// Init starts all the goroutine to ready the netting group to process message
func (ng *NettingGroup) Init(parent context.Context, nettings []Netting, nbrouter int) {
	g1, ctx := ergp.WithContext(parent)
	for _, netting := range nettings {

		in := make(chan NettingRequest)
		total := NewMatrix(1000, 20)
		ne := &NettingEngine{netting: netting, Result: total, in: in, out: ng.results}
		ng.nettings[netting.Name] = ne
		g1.Go(ne.newNettingWorker(ctx))
	}
	go func() {
		g1.Wait()
	}()

	g2, ctx := ergp.WithContext(parent)
	for i := 0; i < nbrouter; i = i + 1 {
		g2.Go(ng.newNettingRouter(ctx, i))
	}
	go func() {
		g2.Wait()
		ng.close()
	}()
}

func (ng *NettingGroup) newNettingRouter(ctx context.Context, name int) routine {
	return func() error {
		for price := range ng.prices {
			engine, err := ng.findNettingEngine(price.trade)
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return fmt.Errorf("Cancelled")
			case engine.in <- NettingRequest{price.trade, price.price}:
				fmt.Printf("routed %v to %v: \n", price.trade.Id, engine.netting.Name)
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

func (ne *NettingEngine) newNettingWorker(ctx context.Context) routine {
	return func() error {
		for nr := range ne.in {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				ne.aggregate(nr)
			}
		}
		fmt.Printf("Aggregation for %v done, result=%v : %v \n", ne.netting.Name, ne.result(), time.Now())
		ne.out <- ne.result()
		return nil
	}
}

func (ne *NettingEngine) aggregate(nr NettingRequest) {
	fmt.Printf("aggregate %v on %v: %v \n", nr.trade.Id, ne.netting.Name, time.Now())
	ne.Result.Add(nr.price)
}

func (ne *NettingEngine) result() float32 {
	value, _ := ne.Result.Get(0, 0)
	return value
}
