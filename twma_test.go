package twma

import (
	"fmt"
	"time"
	"testing"
)

const allowableError = 0.000001

func TestNewTWMA(t *testing.T) {
	e := NewTWMA(time.Second)
	fmt.Printf("NewTWMA: %+v\n",e)
	if (e == nil) {
		t.Errorf("fail.")
	}
}