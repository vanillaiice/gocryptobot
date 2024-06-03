package simple

import "testing"

func TestShouldBuy(t *testing.T) {
	botCfg := &BotCfg{
		PercentChangeBuy: 1.0,
	}

	cases := []struct {
		pChange float32
		expect  bool
	}{
		{pChange: percentChange(1000, 989), expect: true},
		{pChange: percentChange(1000, 990), expect: true},
		{pChange: percentChange(1000, 1000), expect: false},
		{pChange: percentChange(1000, 1001), expect: false},
		{pChange: percentChange(1000, 1010), expect: false},
	}

	for _, c := range cases {
		got := shouldBuy(botCfg, c.pChange)
		if got != c.expect {
			t.Errorf("got %v, want %v for %.5f", got, c.expect, c.pChange)
		}
	}
}

func TestShouldBuyStopEntry(t *testing.T) {
	botCfg := &BotCfg{
		StopEntryPrice:       1000,
		StopEntryPriceMargin: 0.1,
	}

	cases := []struct {
		price  float32
		expect bool
	}{
		{price: 1000.5, expect: true},
		{price: 1001, expect: true},
		{price: 1002, expect: false},
		{price: 999.5, expect: true},
		{price: 999, expect: true},
		{price: 998, expect: false},
	}

	for _, c := range cases {
		got := shouldBuyStopEntry(botCfg, c.price)
		if got != c.expect {
			t.Errorf("got %v, want %v for %.5f", got, c.expect, c.price)
		}
	}
}

func TestShouldSell(t *testing.T) {
	botCfg := &BotCfg{
		PercentChangeSell: 1.0,
	}

	cases := []struct {
		pChange float32
		expect  bool
	}{
		{pChange: percentChange(1000, 1010), expect: true},
		{pChange: percentChange(1000, 1011), expect: true},
		{pChange: percentChange(1000, 1200), expect: true},
		{pChange: percentChange(1000, 1000), expect: false},
		{pChange: percentChange(1000, 990), expect: false},
	}

	for _, c := range cases {
		got := shouldSell(botCfg, c.pChange)
		if got != c.expect {
			t.Errorf("got %v, want %v for %.5f", got, c.expect, c.pChange)
		}
	}
}

func TestShouldSellStopLoss(t *testing.T) {
	botCfg := &BotCfg{
		TrailingStopLossMargin: 1.0,
	}

	cases := []struct {
		pChange float32
		expect  bool
	}{
		{pChange: percentChange(1000, 989), expect: true},
		{pChange: percentChange(1000, 980), expect: true},
		{pChange: percentChange(1000, 1000), expect: false},
		{pChange: percentChange(1000, 1010), expect: false},
		{pChange: percentChange(1000, 1020), expect: false},
	}

	for _, c := range cases {
		got := shouldSellStopLoss(botCfg, c.pChange)
		if got != c.expect {
			t.Errorf("got %v, want %v for %.5f", got, c.expect, c.pChange)
		}
	}
}

func TestPercentChange(t *testing.T) {
	cases := []struct {
		previous float32
		current  float32
		expect   float32
	}{
		{previous: 1000, current: 1000, expect: 0},
		{previous: 1000, current: 1010, expect: -1},
		{previous: 1000, current: 1020, expect: -2},
		{previous: 1000, current: 990, expect: 1},
		{previous: 1000, current: 980, expect: 2},
	}

	for _, c := range cases {
		got := percentChange(c.previous, c.current)
		if got != c.expect {
			t.Errorf("got %v, want %v for %.5f %.5f", got, c.expect, c.previous, c.current)
		}
	}
}
