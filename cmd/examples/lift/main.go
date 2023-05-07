package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/kirillgrishin-tech/vector-go-sdk/pkg/vector"
	"github.com/kirillgrishin-tech/vector-go-sdk/pkg/vectorpb"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	flag.Parse()

	v, err := vector.NewEP(*serial)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = v.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			_, err := v.Conn.SetLiftHeight(
				ctx,
				&vectorpb.SetLiftHeightRequest{
					HeightMm:          250,
					MaxSpeedRadPerSec: 1,
					IdTag:             2000001,
					DurationSec:       .001,
				},
			)
			if err != nil {
				fmt.Println(err)
			}
			stop <- true
			return
		}
	}

}
