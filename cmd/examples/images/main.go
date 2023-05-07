package main

import (
	"context"
	"flag"
	"image/color"

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
			sdk_wrapper.GetEyeColor()
			sdk_wrapper.MoveHead(3.0)
			for l := 1; l < 64; l++ {
				sdk_wrapper.WriteText("TEXT", float64(l), true, 30, true)
			}
			sdk_wrapper.SayText("Text!")
			sdk_wrapper.WriteText("TEXT", 64, true, 1000, false)
			sdk_wrapper.WriteColoredText("This", 32, true, color.RGBA{255, 0, 0, 255}, 500, true)
			sdk_wrapper.WriteColoredText("is", 32, true, color.RGBA{255, 255, 0, 255}, 500, true)
			sdk_wrapper.WriteColoredText("some", 32, true, color.RGBA{0, 255, 0, 255}, 500, true)
			sdk_wrapper.WriteColoredText("text", 32, true, color.RGBA{0, 255, 255, 255}, 500, true)
			sdk_wrapper.SayText("Image display")
			sdk_wrapper.DisplayImage("data/images/birthday-cake.jpg", 5000, true)
			sdk_wrapper.SayText("Image transitions")
			sdk_wrapper.WriteText("SLIDE", 32, true, 2000, true)
			sdk_wrapper.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, sdk_wrapper.IMAGE_TRANSITION_SLIDE_LEFT, 10)
			sdk_wrapper.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, sdk_wrapper.IMAGE_TRANSITION_SLIDE_RIGHT, 10)
			sdk_wrapper.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, sdk_wrapper.IMAGE_TRANSITION_SLIDE_DOWN, 10)
			sdk_wrapper.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, sdk_wrapper.IMAGE_TRANSITION_SLIDE_UP, 10)
			sdk_wrapper.WriteText("FADE", 32, true, 2000, true)
			sdk_wrapper.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, sdk_wrapper.IMAGE_TRANSITION_FADE_IN, 10)
			sdk_wrapper.DisplayImageWithTransition("data/images/birthday-cake.jpg", 2000, sdk_wrapper.IMAGE_TRANSITION_FADE_OUT, 10)
			sdk_wrapper.SayText("Animated gifs")

			sdk_wrapper.SetBackgroundColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
			sdk_wrapper.DisplayAnimatedGif("data/images/animated.gif", sdk_wrapper.ANIMATED_GIF_SPEED_SLOW, 3, false)
			sdk_wrapper.DisplayAnimatedGif("data/images/animated2.gif", sdk_wrapper.ANIMATED_GIF_SPEED_NORMAL, 3, false)
			sdk_wrapper.DisplayAnimatedGif("data/images/animated3.gif", sdk_wrapper.ANIMATED_GIF_SPEED_FAST, 3, false)
			sdk_wrapper.WriteText("ENJOY!", 32, true, 5000, true)
			stop <- true
			return
		}
	}
}
