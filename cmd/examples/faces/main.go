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

			printFaces()

			println("Delete All")
			sdk_wrapper.FaceEnrollmentDeleteAll()

			printFaces()

			sdk_wrapper.FaceEnrollmentStart("Ambrogio")
			time.Sleep(time.Duration(30000) * time.Millisecond)
			stop <- true
			return
		}
	}
}

func printFaces() {
	faces := sdk_wrapper.FaceEnrollmentListAll()
	println("KNOWN FACES")
	for i := 0; i < len(faces); i++ {
		println(fmt.Sprintf("%d) %s", faces[i].FaceId, faces[i].Name))
	}
}
