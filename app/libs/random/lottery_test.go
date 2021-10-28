package random_test

import (
	"gin_template/app/libs/random"
	"testing"
)

func TestLottery(t *testing.T) {
	pools := map[string]float64{
		"1": 10.000001,
		"2": 10.000001,
		"3": 10.000001,
		"4": 10.000001,
		"5": 10.000001,
		"6": 10.000001,
		"7": 10.000001,
		"8": 10.000001,
		"9": 10.000001,
	}

	for i := 0; i < 1; i++ {
		res := random.Lottery(pools, 10, true)

		t.Log(res)
	}
}
