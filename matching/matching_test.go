package matching

import (
	"testing"
	"time"
)

func fixTime(sec int) time.Time {
	return time.Date(2026, 5, 15, 21, 30, 30, 0, time.UTC).Add(time.Second * time.Duration(sec))
}

func TestShouldBeSorted(t *testing.T) {
	tests := []struct {
		name   string
		side   OrderSide
		orders func() []Order
		sorted []int
	}{
		{
			name: "Buy orderbook sorted",
			side: Buy,
			orders: func() []Order {
				return []Order{
					{Id: 0, OrderType: Limit, Ammount: 1, OrderSide: Buy, Price: 100, Timestamp: fixTime(1)},
					{Id: 1, OrderType: Limit, Ammount: 1, OrderSide: Buy, Price: 105, Timestamp: fixTime(2)},
					{Id: 3, OrderType: Limit, Ammount: 1, OrderSide: Buy, Price: 105, Timestamp: fixTime(1)},
					{Id: 4, OrderType: Limit, Ammount: 1, OrderSide: Buy, Price: 105, Timestamp: fixTime(3)},
					{Id: 5, OrderType: Limit, Ammount: 1, OrderSide: Buy, Price: 95, Timestamp: fixTime(1)},
					{Id: 6, OrderType: Limit, Ammount: 1, OrderSide: Buy, Price: 101, Timestamp: fixTime(1)},
				}
			},
			sorted: []int{5, 0, 6, 3, 1, 4},
		},
		{
			name: "Sell orderbook sorted",
			side: Sell,
			orders: func() []Order {
				return []Order{
					{Id: 0, OrderType: Limit, Ammount: 1, OrderSide: Sell, Price: 100, Timestamp: fixTime(1)},
					{Id: 1, OrderType: Limit, Ammount: 1, OrderSide: Sell, Price: 105, Timestamp: fixTime(2)},
					{Id: 3, OrderType: Limit, Ammount: 1, OrderSide: Sell, Price: 105, Timestamp: fixTime(1)},
					{Id: 4, OrderType: Limit, Ammount: 1, OrderSide: Sell, Price: 105, Timestamp: fixTime(3)},
					{Id: 5, OrderType: Limit, Ammount: 1, OrderSide: Sell, Price: 95, Timestamp: fixTime(1)},
					{Id: 6, OrderType: Limit, Ammount: 1, OrderSide: Sell, Price: 101, Timestamp: fixTime(1)},
				}
			},
			sorted: []int{3, 1, 4, 6, 0, 5},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matchEngine := NewMatchingEngine()
			for _, order := range tc.orders() {
				matchEngine.PlaceOrder(order)
			}
			var curSide *[]Order
			switch tc.side {
			case Buy:
				curSide = &matchEngine.orders.Asks
			case Sell:
				curSide = &matchEngine.orders.Bids
			}
			if len(*curSide) != len(tc.sorted) {
				t.Errorf("Size mismatch, expected %d but have %d", len(tc.sorted), len(*curSide))
			}
			for i, order := range *curSide {
				if tc.sorted[i] != order.Id {
					t.Errorf("Id %d at index %d, but expected at %d", order.Id, i, tc.sorted[i])
				}
			}
		})
	}
}

func TestShouldSellMarket(t *testing.T) {
	tests := []struct {
		name string
		ask  Order
		bids func() []Order
		sum  int
	}{
		{
			name: "First order close",
			ask:  Order{Ammount: 10, OrderType: Market, OrderSide: Buy},
			bids: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
				}
			},
			sum: 1000,
		},
		{
			name: "More than one order close",
			ask:  Order{Ammount: 30, OrderType: Market, OrderSide: Buy},
			bids: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
				}
			},
			sum: 3000,
		},
		{
			name: "More than one and equal two orders",
			ask:  Order{Ammount: 40, OrderType: Market, OrderSide: Buy},
			bids: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
				}
			},
			sum: 4000,
		},
		{
			name: "Zero orders",
			ask:  Order{Ammount: 40, OrderType: Market, OrderSide: Buy},
			bids: func() []Order {
				return []Order{}
			},
			sum: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchEngine := NewMatchingEngine()
			for _, order := range tt.bids() {
				matchEngine.PlaceOrder(order)
			}
			sum, err := matchEngine.PlaceOrder(tt.ask)
			if err != nil {
				t.Errorf("Error %e", err)
			}
			if sum != tt.sum {
				t.Errorf("Sum is %d, but expected %d", sum, tt.sum)
			}
		})
	}
}

func TestShouldBuyMarket(t *testing.T) {
	tests := []struct {
		name string
		asks func() []Order
		bid  Order
		sum  int
	}{
		{
			name: "First order close",
			bid:  Order{Ammount: 10, OrderType: Market, OrderSide: Sell},
			asks: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
				}
			},
			sum: 1000,
		},
		{
			name: "More than one order close",
			bid:  Order{Ammount: 30, OrderType: Market, OrderSide: Sell},
			asks: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
				}
			},
			sum: 3000,
		},
		{
			name: "More than one and equal two orders",
			bid:  Order{Ammount: 40, OrderType: Market, OrderSide: Sell},
			asks: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
				}
			},
			sum: 4000,
		},
		{
			name: "Zero orders",
			bid:  Order{Ammount: 40, OrderType: Market, OrderSide: Buy},
			asks: func() []Order {
				return []Order{}
			},
			sum: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchEngine := NewMatchingEngine()
			for _, order := range tt.asks() {
				matchEngine.PlaceOrder(order)
			}
			sum, err := matchEngine.PlaceOrder(tt.bid)
			if err != nil {
				t.Errorf("Error %e", err)
			}
			if sum != tt.sum {
				t.Errorf("Sum is %d, but expected %d", sum, tt.sum)
			}
		})
	}
}

func TestShouldSellLimit(t *testing.T) {
	tests := []struct {
		name string
		ask  Order
		bids func() []Order
		sum  int
	}{
		{
			name: "First order close",
			ask:  Order{Ammount: 10, Price: 110, OrderType: Limit, OrderSide: Buy},
			bids: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
				}
			},
			sum: 1000,
		},
		{
			name: "More than one order close",
			ask:  Order{Ammount: 30, Price: 110, OrderType: Limit, OrderSide: Buy},
			bids: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
				}
			},
			sum: 3000,
		},
		{
			name: "More than one and equal two orders",
			ask:  Order{Ammount: 40, Price: 110, OrderType: Limit, OrderSide: Buy},
			bids: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
				}
			},
			sum: 4000,
		},
		{
			name: "No one is lower your limit price",
			ask:  Order{Ammount: 40, Price: 90, OrderType: Limit, OrderSide: Buy},
			bids: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
				}
			},
			sum: 0,
		},
		{
			name: "Two is lower your limit price, your order is not finished",
			ask:  Order{Ammount: 60, Price: 99, OrderType: Limit, OrderSide: Buy},
			bids: func() []Order {
				return []Order{
					{Ammount: 20, Price: 90, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 95, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
				}
			},
			sum: 3700,
		},
		{
			name: "Two is lower your limit price, your order is finished",
			ask:  Order{Ammount: 40, Price: 95, OrderType: Limit, OrderSide: Buy},
			bids: func() []Order {
				return []Order{
					{Ammount: 20, Price: 90, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 95, OrderType: Limit, OrderSide: Sell, Timestamp: fixTime(0)},
				}
			},
			sum: 3700,
		},
		{
			name: "Zero orders",
			ask:  Order{Ammount: 40, Price: 110, OrderType: Limit, OrderSide: Buy},
			bids: func() []Order {
				return []Order{}
			},
			sum: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchEngine := NewMatchingEngine()
			for _, order := range tt.bids() {
				matchEngine.PlaceOrder(order)
			}
			sum, err := matchEngine.PlaceOrder(tt.ask)
			if err != nil {
				t.Errorf("Error %e", err)
			}
			if sum != tt.sum {
				t.Errorf("Sum is %d, but expected %d", sum, tt.sum)
			}
		})
	}
}

func TestShouldBuyLimit(t *testing.T) {
	tests := []struct {
		name string
		bid  Order
		asks func() []Order
		sum  int
	}{
		{
			name: "First order close",
			bid:  Order{Ammount: 10, Price: 90, OrderType: Limit, OrderSide: Sell},
			asks: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
				}
			},
			sum: 1000,
		},
		{
			name: "More than one order close",
			bid:  Order{Ammount: 30, Price: 90, OrderType: Limit, OrderSide: Sell},
			asks: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
				}
			},
			sum: 3000,
		},
		{
			name: "More than one and equal two orders",
			bid:  Order{Ammount: 40, Price: 90, OrderType: Limit, OrderSide: Sell},
			asks: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
				}
			},
			sum: 4000,
		},
		{
			name: "No one is greater your limit price",
			bid:  Order{Ammount: 40, Price: 101, OrderType: Limit, OrderSide: Sell},
			asks: func() []Order {
				return []Order{
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
				}
			},
			sum: 0,
		},
		{
			name: "Two is greater your limit price, your order is not finished",
			bid:  Order{Ammount: 60, Price: 95, OrderType: Limit, OrderSide: Sell},
			asks: func() []Order {
				return []Order{
					{Ammount: 20, Price: 90, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 95, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
				}
			},
			sum: 3900,
		},
		{
			name: "Two is lower than your limit price, your order is finished",
			bid:  Order{Ammount: 40, Price: 95, OrderType: Limit, OrderSide: Sell},
			asks: func() []Order {
				return []Order{
					{Ammount: 20, Price: 90, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 100, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
					{Ammount: 20, Price: 95, OrderType: Limit, OrderSide: Buy, Timestamp: fixTime(0)},
				}
			},
			sum: 3900,
		},
		{
			name: "Zero orders",
			bid:  Order{Ammount: 40, Price: 110, OrderType: Limit, OrderSide: Sell},
			asks: func() []Order {
				return []Order{}
			},
			sum: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchEngine := NewMatchingEngine()
			for _, order := range tt.asks() {
				matchEngine.PlaceOrder(order)
			}
			sum, err := matchEngine.PlaceOrder(tt.bid)
			if err != nil {
				t.Errorf("Error %e", err)
			}
			if sum != tt.sum {
				t.Errorf("Sum is %d, but expected %d", sum, tt.sum)
			}
		})
	}
}
