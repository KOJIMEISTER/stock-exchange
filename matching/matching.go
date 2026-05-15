package matching

import (
	"cmp"
	"log"
	"slices"
	"time"
)

type OrderType int

const (
	TypeUnknow OrderType = iota
	Market
	Limit
)

type OrderSide int

const (
	SideUknow OrderSide = iota
	Buy
	Sell
)

type Order struct {
	Id        int
	Timestamp time.Time
	OrderType OrderType
	OrderSide OrderSide
	Ammount   int
	Price     int
}

type orderBook struct {
	Asks []Order
	Bids []Order
}

type MatchingEngine struct {
	orders orderBook
}

func NewMatchingEngine() *MatchingEngine {
	return &MatchingEngine{orderBook{make([]Order, 0, 20), make([]Order, 0, 20)}}
}

func (this *MatchingEngine) PlaceOrder(order Order) error {
	switch order.OrderSide {
	case Buy:
		{
			this.orders.Asks = append(this.orders.Asks, order)
			slices.SortFunc(this.orders.Asks, func(left, right Order) int {
				return cmp.Or(cmp.Compare(left.Price, right.Price), left.Timestamp.Compare(right.Timestamp))
			})
			log.Printf("Order id %d placed at asks", order.Id)
		}
	case Sell:
		{
			this.orders.Bids = append(this.orders.Bids, order)
			slices.SortFunc(this.orders.Bids, func(right, left Order) int {
				return cmp.Or(cmp.Compare(left.Price, right.Price), right.Timestamp.Compare(left.Timestamp))
			})
			log.Printf("Order id %d placed at bids", order.Id)
		}
	default:
	}
	return nil
}
