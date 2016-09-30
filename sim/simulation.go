package sim

//Market is the object representing some Market Data used for simulation
type Market struct {
	Matrix
}

//Price unexpoted
func (mkt *Market) Price(trd Trade) (TradeSimulation, error) {
	price, _ := mkt.Mult(trd.Mtm)
	return TradeSimulation{trd, price}, nil
}
