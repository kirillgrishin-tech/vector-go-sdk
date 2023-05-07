package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	sdk_wrapper "github.com/kirillgrishin-tech/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/kirillgrishin-tech/vector-go-sdk/pkg/vectorpb"
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
			go func() {
				println("Listening for events...")
				isBusy := false
				for {
					evt := sdk_wrapper.WaitForEvent()
					if evt != nil {
						evtUserIntent := evt.GetUserIntent()
						if evtUserIntent != nil {
							println(fmt.Sprintf("Received intent %d", evtUserIntent.IntentId))
							println(evtUserIntent.JsonData)
							println(evtUserIntent.String())
						}
						evtRobotState := evt.GetRobotState()
						if evtRobotState != nil {
							if evtRobotState.Status&uint32(vectorpb.RobotStatus_ROBOT_STATUS_IS_PICKED_UP) > 0 ||
								evtRobotState.Status&uint32(vectorpb.RobotStatus_ROBOT_STATUS_IS_BEING_HELD) > 0 {
								println("I am being picked up.")
							}
							if evtRobotState.TouchData.IsBeingTouched == true && !isBusy {
								isBusy = true
								go func() {
									println("I am being touched.")
									//sdk_wrapper.PlayAnimation("anim_eyepose_angry", 0, false, false, false)
									_ = sdk_wrapper.PlaySound(sdk_wrapper.GetDataPath("audio/roar.wav"))
									time.Sleep(time.Duration(1000) * time.Millisecond)
									isBusy = false
								}()
							}
						}
					}
				}
			}()
			for {
				time.Sleep(time.Duration(100) * time.Millisecond)
			}
			stop <- true
			return
		}
	}
}
