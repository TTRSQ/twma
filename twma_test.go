package twma

import (
	"fmt"
	"time"
	"testing"
)

const allowableError = 0.000001

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}

func TestOrderedList(t *testing.T) {
	ma := NewTWMA(time.Second*10) // 10 sec window
	list := []Item{
		Item{
			Value: 1,
			Time: time.Date(2001, 5, 20, 23, 59, 0, 0, time.Local), 
		},
		Item{
			Value: 1,
			Time: time.Date(2001, 5, 20, 23, 59, 1, 0, time.Local), 
		},
		Item{
			Value: 1,
			Time: time.Date(2001, 5, 20, 23, 59, 3, 0, time.Local), 
		},
		Item{
			Value: 1,
			Time: time.Date(2001, 5, 20, 23, 59, 11, 0, time.Local), 
		},
	}
	expect := 1.0
	for _, item := range list {
		ma.Add(item)
	}
	result, _ := ma.Value()
	if absDiff(result, expect) > allowableError {
		fmt.Printf("%+v\n", ma)
		t.Errorf("result: %.2f, expected value: %.2f", result, expect)
	}
}

func TestUnOrderedList(t *testing.T) {
	ma := NewTWMA(time.Second*10) // 10 sec window
	list := []Item{
		Item{
			Value: 1,
			Time: time.Date(2001, 5, 20, 23, 59, 0, 0, time.Local), 
		},
		Item{
			Value: 1,
			Time: time.Date(2001, 5, 20, 23, 59, 3, 0, time.Local), 
		},
		Item{
			Value: 1,
			Time: time.Date(2001, 5, 20, 23, 59, 11, 0, time.Local), 
		},
		Item{
			Value: 1,
			Time: time.Date(2001, 5, 20, 23, 59, 1, 0, time.Local), 
		},
	}
	expect := 1.0
	for _, item := range list {
		ma.Add(item)
	}
	result, _ := ma.Value()
	if absDiff(result, expect) > allowableError {
		fmt.Printf("%+v\n", ma)
		t.Errorf("result: %.2f, expected value: %.2f", result, expect)
	}
}