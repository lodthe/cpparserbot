package api

import (
	"os"
	"testing"

	"github.com/lodthe/cpparserbot/model"
)

func TestBinancePairValidation(t *testing.T) {
	in := []*model.Pair{
		{"ETH", "BTC"},
		{"NONEXISTENT", "PAIR"},
	}
	want := []bool{
		false,
		true,
	}

	b := Binance{}
	b.Init(os.Getenv("BINANCE_API_KEY"), os.Getenv("BINANCE_SECRET_KEY"))

	for i, pair := range in {
		if _, err := b.GetKlines(pair); (err != nil) != want[i] {
			t.Errorf("For GetKlines(%s) expecting error: %v, got %v", pair, want[i], err != nil)
		}

		if _, err := b.GetPrice(pair); (err != nil) != want[i] {
			t.Errorf("For GetPrice(%s) expecting error: %v, got %v", pair, want[i], err != nil)
		}
	}

	if _, err := b.GetAllPrices(); err != nil {
		t.Errorf("For GetAllPrice() expecting no error, but got %s", err)
	}
}
