package sim

import "strconv"

//Trade represent a trade in the system
type Trade struct {
	Id   int
	Desc string
	Mtm  float32
}

type TradeSimulation struct {
	Trade
	Matrix
}

func NewTrade(i int) Trade {
	return Trade{i, "id" + strconv.Itoa(i), float32(i) * 0.5}
}
