package helper

import (
	"testing"

	"github.com/lodthe/cpparserbot/model"
)

func TestFindPairInConfig(t *testing.T) {
	binancePairs = []model.Pair{
		{"BTC", "USDT"},
		{"ETH", "BTC"},
		{"XRP", "BTC"},
		{"DASH", "BTC"},
		{"ETC", "BTC"},
	}
	in := []model.Pair{
		{"BTC", "USDT"},
		{"ETH", "BTC"},
		{"XRP", "BTC"},
		{"BNB", "BTC"},
		{"BTC", "BTC"},
		{"NO", "SUCH"},
	}
	want := []*model.Pair{
		{"BTC", "USDT"},
		{"ETH", "BTC"},
		{"XRP", "BTC"},
		nil,
		nil,
		nil,
	}

	for i := range in {
		pair := FindPairInConfig(in[i].String())
		if (pair != want[i]) && ((pair != nil) && (want[i] != nil) && (*pair != *want[i])) {
			t.Errorf("For %s expected %s, got %s", in[i], want[i], pair)
		}
	}
}
