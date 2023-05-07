package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	sdk_wrapper "github.com/kirillgrishin-tech/vector-go-sdk/pkg/sdk-wrapper"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	flag.Parse()

	sdk_wrapper.InitSDK(*serial)

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			println("")
			println("ANIMATION LIST:")
			aList := sdk_wrapper.LoadAnimationList()
			if nil != aList {
				println(fmt.Sprintf("%d animations found.", len(aList)))
				for i := 0; i < len(aList); i++ {
					println(aList[i])
					sdk_wrapper.PlayAnimation(aList[i], 1, false, false, false)
					time.Sleep(time.Duration(5000) * time.Millisecond)
				}
			} else {
				println("Could not load animation list :(")
			}
			stop <- true
			return
		}
	}
}
