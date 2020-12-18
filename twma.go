package twma

import (
	"errors"
	"sort"
	"time"
)

// Item window item to Input
type Item struct {
	Value float64
	Time  time.Time
}

type item struct {
	value  float64
	time   float64 // unixsec
	weight float64 // average = sum(weight_1 * value_1) / windowSize
}

// TimeWindowedMovingAverage ..
type TimeWindowedMovingAverage struct {
	windowSize float64
	window     []item
	sum        float64
}

// NewTWMA ..
func NewTWMA(windowSize time.Duration) *TimeWindowedMovingAverage {
	twma := new(TimeWindowedMovingAverage)
	twma.windowSize = float64(windowSize) / float64(time.Second)
	return twma
}

// Apply .. add new item and get deleted item
func (tw *TimeWindowedMovingAverage) Apply(addItem Item) []Item {
	tItem := item{
		value: addItem.Value,
		time:  float64(addItem.Time.UnixNano()) / float64(time.Second),
	}

	if len(tw.window) == 0 {
		tw.window = append(tw.window, tItem)
		return []Item{}
	} else if tw.window[len(tw.window)-1].time < tItem.time {
		diff := tItem.time - tw.window[len(tw.window)-1].time
		tItem.weight = diff / 2
		return translateItem(tw.addLast(tItem))
	}

	tw.window = append(tw.window, tItem)
	sort.Slice(tw.window, func(i, j int) bool {
		return tw.window[i].time < tw.window[j].time
	})
	deletedList := tw.adjustWindow()
	tw.calcWeight()
	return translateItem(deletedList)
}

func translateItem(l []item) []Item {
	letList := []Item{}
	for _, v := range l {
		itime := int64(v.time)
		letList = append(letList, Item{
			Value: v.value,
			Time:  time.Unix(itime, int64((v.time-float64(itime))*1000000000)),
		})
	}
	return letList
}

func (tw *TimeWindowedMovingAverage) addLast(addItem item) []item {
	tw.window[len(tw.window)-1].weight += addItem.weight
	tw.sum += tw.window[len(tw.window)-1].value * addItem.weight

	tw.window = append(tw.window, addItem)
	tw.sum += addItem.value * addItem.weight

	return tw.adjustWindow()
}

// adjustWindow .. remove Items before (lastItem.Time - windowSize)
// and adjust weight of first Item.
func (tw *TimeWindowedMovingAverage) adjustWindow() []item {
	divider := 0
	lastTime := tw.window[len(tw.window)-1].time
	for idx := range tw.window {
		if lastTime-tw.window[idx].time < tw.windowSize {
			break
		}
		divider = idx
	}
	deleteList := tw.window[0:divider]
	tw.window = tw.window[divider:len(tw.window)]

	for idx := range deleteList {
		tw.sum -= deleteList[idx].value * deleteList[idx].weight
	}

	tw.sum -= tw.window[0].value * tw.window[0].weight
	tw.window[0].weight = (tw.window[1].time - tw.window[0].time) / 2
	tw.sum += tw.window[0].value * tw.window[0].weight

	return deleteList
}

func (tw *TimeWindowedMovingAverage) calcWeight() {
	firstIdx := 0
	lastIdx := len(tw.window) - 1
	tw.sum = 0.0
	for idx := range tw.window {
		diff := 0.0
		if idx == firstIdx {
			diff = tw.window[firstIdx+1].time - tw.window[firstIdx].time
		} else if idx == lastIdx {
			diff = tw.window[lastIdx].time - tw.window[lastIdx-1].time
		} else {
			diff += tw.window[idx+1].time - tw.window[idx].time
			diff += tw.window[idx].time - tw.window[idx-1].time
		}
		tw.window[idx].weight = diff / 2
		tw.sum += tw.window[idx].value * tw.window[idx].weight
	}
}

// Value calc TimeWindowedMovingAverage
func (tw *TimeWindowedMovingAverage) Value() (float64, error) {
	if len(tw.window) < 1 {
		return 0, errors.New("insufficient number of items")
	} else if float64(tw.window[len(tw.window)-1].time-tw.window[0].time) < tw.windowSize {
		return 0, errors.New("insufficient number of items")
	}
	return tw.sum / (tw.window[len(tw.window)-1].time - tw.window[0].time), nil
}
