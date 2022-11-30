package main

import (
	"context"
	"flag"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"image/color"
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
			sdk_wrapper.MoveHead(3.0)
			sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})

			sdk_wrapper.UseVectorEyeColorInImages(true)
			sdk_wrapper.WriteText("Hello, World", 64, true, 5000, true)

			/*
				s1 := rand.NewSource(time.Now().UnixNano())
				r1 := rand.New(s1)
				die := r1.Intn(6) + 1
				dieImage := fmt.Sprintf("data/images/dice/%d.png", die)

				sdk_wrapper.DisplayAnimatedGif("data/images/dice/roll-the-dice.gif", sdk_wrapper.ANIMATED_GIF_SPEED_FASTEST, 1, false)
				sdk_wrapper.DisplayImage(dieImage, 100, false)
				sdk_wrapper.PlaySound(sdk_wrapper.SYSTEMSOUND_WIN)
				sdk_wrapper.SayText(fmt.Sprintf("You rolled a %d", die))
				sdk_wrapper.DisplayImageWithTransition(dieImage, 1000, sdk_wrapper.IMAGE_TRANSITION_FADE_OUT, 10)
			*/
			stop <- true
			return
		}
	}
}
