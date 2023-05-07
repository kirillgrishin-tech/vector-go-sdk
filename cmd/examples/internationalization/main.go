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
			sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_HTGO)
			sdk_wrapper.SetRobotName("Augustus")
			loc := sdk_wrapper.GetLocale()
			println("Robot locale: " + loc)
			sdk_wrapper.SetLocale("it_IT")
			sdk_wrapper.SayText("Ciao mondo!")
			println("Robot locale: " + loc)
			sdk_wrapper.SetLocale("en_US")
			sdk_wrapper.SayText("Hello world!")
			sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_SPANISH)
			sdk_wrapper.SayText(sdk_wrapper.Translate("Hello, world!", sdk_wrapper.LANGUAGE_ENGLISH, sdk_wrapper.LANGUAGE_SPANISH))
			sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_FRENCH)
			sdk_wrapper.SayText(sdk_wrapper.Translate("Hello, world!", sdk_wrapper.LANGUAGE_ENGLISH, sdk_wrapper.LANGUAGE_FRENCH))
			sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_GERMAN)
			sdk_wrapper.SayText(sdk_wrapper.Translate("Hello, world!", sdk_wrapper.LANGUAGE_ENGLISH, sdk_wrapper.LANGUAGE_GERMAN))
			sdk_wrapper.SetLanguage(sdk_wrapper.LANGUAGE_JAPANESE)
			sdk_wrapper.SayText(sdk_wrapper.Translate("Hello, world!", sdk_wrapper.LANGUAGE_ENGLISH, sdk_wrapper.LANGUAGE_JAPANESE))
			stop <- true
			return
		}
	}
}
