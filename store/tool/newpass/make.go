package main

import (
	"fmt"
	"github.com/lighthouse-and-whale/quark/core"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Printf("%s\n", core.RandAz09(16))
		time.Sleep(time.Millisecond * 100)
	}
}

/*
8EAsMybv1jhX_A1tc
*/
