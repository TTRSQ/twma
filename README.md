# twma
In real time data stream, calculating moving average with time window.  
twma can be used for sensor data, stock trading information, and etc.  
<img src="https://user-images.githubusercontent.com/26806928/74592120-81ec0580-5061-11ea-9d6d-66b2d10368bb.png" width="70%" alt="btc-price in 2020-2-14">

## usage

```
package main

import (
	"fmt"
	"time"

	"github.com/TTRSQ/twma"
)

func main() {
	// initial item
	v1 := 1.0
	t1 := time.Now()

	ma := twma.NewTWMA(time.Second * 9)
	for i := 0; i < 10; i++ {
		// add with increment
		ma.Add(twma.Item{
			Value: v1 + float64(i),
			Time:  t1.Add(time.Second * time.Duration(i)),
		})
	}

	ave, _ := ma.Value()
	fmt.Printf("moving average: %.2f", ave) // moving average: 5.50
}
```
## explanation

![twma](https://user-images.githubusercontent.com/26806928/74354165-50273480-4dfe-11ea-8d2a-b22432d116ea.jpeg)
