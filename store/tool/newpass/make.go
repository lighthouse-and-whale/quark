package main

import (
	"fmt"
	"lib/core"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Printf("%s\n", core.RandAz09(32))
		time.Sleep(time.Millisecond * 100)
	}
}
