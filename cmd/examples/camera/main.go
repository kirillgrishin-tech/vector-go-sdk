package main

import (
	"context"
	"flag"

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
			sdk_wrapper.PlayAnimation("anim_generic_look_up_01", 0, false, false, false)
			sdk_wrapper.SaveHiResCameraPicture("camera.jpg")
			stop <- true
			return
		}
	}
}
