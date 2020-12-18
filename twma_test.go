package twma

import (
	"fmt"
	"testing"
	"time"
)

const allowableError = 0.000001

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}

func TestOne(t *testing.T) {
	items := []Item{
		{
			Value: 1,
			Time:  time.Date(2001, 5, 20, 23, 59, 0, 0, time.Local),
		},
		{
			Value: 1,
			Time:  time.Date(2001, 5, 20, 23, 59, 1, 0, time.Local),
		},
		{
			Value: 1,
			Time:  time.Date(2001, 5, 20, 23, 59, 3, 0, time.Local),
		},
		{
			Value: 1,
			Time:  time.Date(2001, 5, 20, 23, 59, 11, 0, time.Local),
		},
	}
	windowSizeSec := 10
	expect := 1.0

	// order by time ask.
	testItems(items, windowSizeSec, expect, 1, t)
	// order by time desc.
	testItemsDesc(items, windowSizeSec, expect, 1, t)
}

func TestLinear(t *testing.T) {
	v1 := 0.0
	t1 := time.Now()

	items := []Item{}
	for i := 0; i <= 10; i++ {
		items = append(items, Item{
			Value: v1 + float64(i),
			Time:  t1.Add(time.Second * time.Duration(i)),
		})
	}
	windowSizeSec := 10
	expect := 5.0

	// order by time ask.
	testItems(items, windowSizeSec, expect, 0, t)
	// order by time desc.
	testItemsDesc(items, windowSizeSec, expect, 0, t)
}

func TestPopItems(t *testing.T) {
	v1 := 0.0
	t1 := time.Now()

	items := []Item{}
	for i := 0; i <= 10; i++ {
		items = append(items, Item{
			Value: v1 + float64(i),
			Time:  t1.Add(time.Second * time.Duration(i)),
		})
	}
	// expect 3 item will pop that values (0.0, 1.0, 2.0)
	windowSizeSec := 7
	expect := float64(3.0/2+4+5+6+7+8+9+10.0/2) / 7
	expectPopCount := 3

	// order by time ask.
	testItems(items, windowSizeSec, expect, expectPopCount, t)
	// order by time desc.
	testItemsDesc(items, windowSizeSec, expect, expectPopCount, t)
}

func testItems(items []Item, windowSizeSec int, expect float64, popCount int, t *testing.T) {
	ma := NewTWMA(time.Second * time.Duration(windowSizeSec))
	popItems := []Item{}
	for _, item := range items {
		popItems = append(popItems, ma.Apply(item)...)
	}
	result, _ := ma.Value()
	if len(popItems) != popCount {
		t.Errorf("poped items count is not valid: %d", len(popItems))
	}
	if absDiff(result, expect) > allowableError {
		fmt.Printf("%+v\n", ma)
		t.Errorf("result: %.2f, expected value: %.2f", result, expect)
	}
}

func testItemsDesc(items []Item, windowSizeSec int, expect float64, popCount int, t *testing.T) {
	ma := NewTWMA(time.Second * time.Duration(windowSizeSec))
	popItems := []Item{}
	for idx := range items {
		rIdx := len(items) - 1 - idx
		popItems = append(popItems, ma.Apply(items[rIdx])...)
	}
	result, _ := ma.Value()
	if len(popItems) != popCount {
		t.Errorf("poped items count is not valid: %d", len(popItems))
	}
	if absDiff(result, expect) > allowableError {
		fmt.Printf("%+v\n", ma)
		t.Errorf("result: %.2f, expected value: %.2f", result, expect)
	}
}
