package matching

import (
	"log"
	"time"
)

type OrderType int

const (
	TypeUnknow OrderType = iota
	TypeBuy
	TypeSell
)

type Order struct {
	Id        int
	Timestamp time.Time
	OrderType OrderType
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
	return &MatchingEngine{orderBook{make([]Order, 20), make([]Order, 20)}}
}

func (this *MatchingEngine) PlaceOrder(order Order) error {
	switch order.OrderType {
	case TypeBuy:
		{
			this.orders.Asks = append(this.orders.Asks, order)
			log.Printf("Order id %d placed at asks", order.Id)
			break
		}
	case TypeSell:
		{
			this.orders.Bids = append(this.orders.Bids, order)
			log.Printf("Order id %d placed at bids", order.Id)
			break
		}
	default:
	}
	return nil
}
