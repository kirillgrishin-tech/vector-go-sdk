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
			temperature := 16
			unit := sdk_wrapper.GetTemperatureUnit()
			if unit == sdk_wrapper.WEATHER_UNIT_FARANHEIT {
				temperature = (temperature * 9 / 5) + 32
			}
			sdk_wrapper.UseVectorEyeColorInImages(true)
			err1 := sdk_wrapper.DisplayTemperature(temperature, unit, 5000, true)
			if err1 != nil {
				println("ERROR " + err1.Error())
			}
			sdk_wrapper.UseVectorEyeColorInImages(false)
			sdk_wrapper.DisplayCondition("shower rain", "09n", 5000, true)
			stop <- true
			return
		}
	}
}
