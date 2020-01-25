package model

import "testing"

var currencies = []string{"USDT", "BTC", "ETH", "RUB", "DASH", "LARGE"}

func TestString(t *testing.T) {
	for _, i := range currencies {
		for _, j := range currencies {
			pair := Pair{i, j}
			if pair.String() != i+"/"+j {
				t.Errorf("For %s, %s expected %s/%s, got %s", i, j, i, j, pair.String())
			}
		}
	}
}

func TestBinanceFormat(t *testing.T) {
	for _, i := range currencies {
		for _, j := range currencies {
			pair := Pair{i, j}
			if pair.ToBinanceFormat() != i+j {
				t.Errorf("For %s, %s expected %s, got %s", i, j, i+j, pair.ToBinanceFormat())
			}
		}
	}
}
