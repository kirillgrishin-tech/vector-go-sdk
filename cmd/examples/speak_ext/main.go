package main

import (
	"context"
	"flag"
	"math/rand"
	"time"

	sdk_wrapper "github.com/kirillgrishin-tech/vector-go-sdk/pkg/sdk-wrapper"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	var sentence = flag.String("sentence", "", "Sentence to say")
	flag.Parse()

	sdk_wrapper.InitSDK(*serial)

	sentences := [10]string{
		"Fortune favors the bold",
		"I think, therefore I am",
		"Time is money",
		"I came, I saw, I conquered",
		"When life gives you lemons, make lemonade",
		"Practice makes perfect",
		"Knowledge is power",
		"Have no fear of perfection, you'll never reach it",
		"No pain no gain",
		"That which does not kill us makes us stronger",
	}

	var phrase = ""

	if *sentence == "" {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		phrase = sentences[r1.Intn(10)]
	} else {
		phrase = *sentence
	}

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			sdk_wrapper.SayText(phrase)
			stop <- true
			return
		}
	}
}
