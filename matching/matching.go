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

func (this *OrderSide) toString() string {
	switch *this {
	case Buy:
		return "Ask"
	case Sell:
		return "Bid"
	default:
		return ""
	}
}

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

func (this *MatchingEngine) placeBid(order *Order) error {
	this.orders.Bids = append(this.orders.Bids, *order)
	slices.SortFunc(this.orders.Bids, func(right, left Order) int {
		return cmp.Or(cmp.Compare(left.Price, right.Price), right.Timestamp.Compare(left.Timestamp))
	})
	log.Printf("Order id %d placed at %s", order.Id, order.OrderSide.toString())
	return nil
}

func (this *MatchingEngine) placeAsk(order *Order) error {
	this.orders.Asks = append(this.orders.Asks, *order)
	slices.SortFunc(this.orders.Asks, func(left, right Order) int {
		return cmp.Or(cmp.Compare(left.Price, right.Price), left.Timestamp.Compare(right.Timestamp))
	})
	log.Printf("Order id %d placed at %s", order.Id, order.OrderSide.toString())
	return nil
}

func (this *MatchingEngine) executeMarketSell(order *Order) (int, error) {
	sum := 0
	curAmm := order.Ammount
	for len(this.orders.Asks) != 0 && curAmm != 0 {
		if curAmm-this.orders.Asks[0].Ammount < 0 {
			log.Printf("Order id %d finished, order side %s", order.Id, order.OrderSide.toString())
			sum = sum + curAmm*this.orders.Asks[0].Price
			this.orders.Asks[0].Ammount = this.orders.Asks[0].Ammount - curAmm
			curAmm = 0
		} else if curAmm-this.orders.Asks[0].Ammount == 0 {
			log.Printf("Order id %d finished, order side %s", order.Id, order.OrderSide.toString())
			log.Printf("Order id %d finished, order side %s", this.orders.Asks[0].Id, this.orders.Asks[0].OrderSide.toString())
			sum = sum + curAmm*this.orders.Asks[0].Price
			this.orders.Asks = this.orders.Asks[1:]
			curAmm = 0
		} else {
			log.Printf("Order id %d finished, order side %s", this.orders.Asks[0].Id, this.orders.Asks[0].OrderSide.toString())
			sum = sum + curAmm*this.orders.Asks[0].Price
			curAmm = curAmm - this.orders.Asks[0].Ammount
			this.orders.Asks = this.orders.Asks[1:]
		}
	}
	if curAmm != 0 {
		this.placeBid(order)
	}
	return sum, nil
}

func (this *MatchingEngine) executeMarketBuy(order *Order) (int, error) {
	sum := 0
	curAmm := order.Ammount
	for len(this.orders.Bids) != 0 && curAmm != 0 {
		if curAmm-this.orders.Bids[0].Ammount < 0 {
			log.Printf("Order id %d finished, order side %s", order.Id, order.OrderSide.toString())
			sum = sum + curAmm*this.orders.Bids[0].Price
			this.orders.Bids[0].Ammount = this.orders.Bids[0].Ammount - curAmm
			curAmm = 0
		} else if curAmm-this.orders.Bids[0].Ammount == 0 {
			log.Printf("Order id %d finished, order side %s", order.Id, order.OrderSide.toString())
			log.Printf("Order id %d finished, order side %s", this.orders.Bids[0].Id, this.orders.Bids[0].OrderSide.toString())
			sum = sum + curAmm*this.orders.Bids[0].Price
			this.orders.Bids = this.orders.Bids[1:]
			curAmm = 0
		} else {
			log.Printf("Order id %d finished, order side %s", this.orders.Bids[0].Id, this.orders.Bids[0].OrderSide.toString())
			sum = sum + curAmm*this.orders.Bids[0].Price
			curAmm = curAmm - this.orders.Bids[0].Ammount
			this.orders.Bids = this.orders.Bids[1:]
		}
	}
	if curAmm != 0 {
		this.placeAsk(order)
	}
	return sum, nil
}

func (this *MatchingEngine) executeLimitSell(order *Order) (int, error) {
	sum := 0
	curAmm := order.Ammount
	if curAmm != 0 {
		this.placeBid(order)
	}
	return sum, nil
}

func (this *MatchingEngine) executeLimitBuy(order *Order) (int, error) {
	sum := 0
	curAmm := order.Ammount
	if curAmm != 0 {
		this.placeAsk(order)
	}
	return sum, nil
}

func (this *MatchingEngine) executeSell(order *Order) (int, error) {
	switch order.OrderType {
	case Market:
		return this.executeMarketSell(order)
	case Limit:
		return this.executeLimitSell(order)
	default:
		return 0, nil
	}
}

func (this *MatchingEngine) executeBuy(order *Order) (int, error) {
	switch order.OrderType {
	case Market:
		return this.executeMarketBuy(order)
	case Limit:
		return this.executeLimitBuy(order)
	default:
		return 0, nil
	}
}

func (this *MatchingEngine) PlaceOrder(order Order) error {
	switch order.OrderSide {
	case Buy:
		this.executeBuy(&order)
	case Sell:
		this.executeSell(&order)
	default:
	}
	return nil
}
