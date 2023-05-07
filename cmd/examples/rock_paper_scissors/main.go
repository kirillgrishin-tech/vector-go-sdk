package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	sdk_wrapper "github.com/kirillgrishin-tech/vector-go-sdk/pkg/sdk-wrapper"
)

/*
 Rock paper scissors icons found on flaticon.com

 <a href="https://www.flaticon.com/free-icons/rock-paper-scissors" title="rock paper scissors icons">Rock paper scissors icons created by Freepik - Flaticon</a>
*/

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	flag.Parse()

	sdk_wrapper.InitSDK(*serial)

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		//time.Sleep(time.Duration(10000) * time.Millisecond)
		_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			//playDemo()
			sdk_wrapper.SayText("Let's play!")
			playGame(10)
			sdk_wrapper.SayText("Ok, I think it's enough")
			stop <- true
			return
		}
	}
}

func playDemo() {
	sdk_wrapper.WriteText("OPENCV", 32, true, 10000, true)
	sdk_wrapper.WriteText("MEDIAPIPE", 32, true, 10000, true)
	sdk_wrapper.WriteText("0.", 64, true, 2000, true)
	sdk_wrapper.WriteText("0..", 64, true, 2000, true)
	sdk_wrapper.WriteText("0...", 64, true, 2000, true)
	sdk_wrapper.WriteText("0...5", 64, true, 2000, true)
	sdk_wrapper.WriteText("0", 64, true, 2000, true)
	sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("images/rock.png"), 2000, true)
	sdk_wrapper.WriteText("2", 64, true, 1000, true)
	sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("images/scissors.png"), 2000, true)
	sdk_wrapper.WriteText("5", 64, true, 1000, true)
	sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("images/paper.png"), 2000, true)
}

func playGame(numSteps int) {
	myScore := 0
	userScore := 0
	options := [3]string{
		"rock",
		"paper",
		"scissors",
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i <= numSteps; i++ {
		sdk_wrapper.SayText("one, two, three!")

		myMove := options[r1.Intn(len(options))]
		sdk_wrapper.DisplayImage(sdk_wrapper.GetDataPath("images/"+myMove+".png"), 5000, false)
		fName := sdk_wrapper.GetTemporaryFilename("rps", "jpg", true)
		err := sdk_wrapper.SaveHiResCameraPicture(fName)
		if err == nil {
			if err == nil {
				cmd := exec.Command("python", "hand_detection.py", fName)
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}

				var output string = out.String()
				output = regexp.MustCompile(`[^0-9]`).ReplaceAllString(output, "")

				var numFingers int = -1
				if len(output) > 0 {
					numFingers, _ = strconv.Atoi(output)
				}
				win := 0
				answer := ""
				userMove := ""

				fmt.Printf("num fingers %d, Output: %q\n", numFingers, output)

				switch numFingers {
				case 0:
					// User plays "rock"
					userMove = "rock"
					if myMove == "paper" {
						win = 1
					} else if myMove == "scissors" {
						win = -1
					}
					break
				case 2:
					// User plays "scissors"
					userMove = "scissors"
					if myMove == "rock" {
						win = 1
					} else if myMove == "paper" {
						win = -1
					}
					break
				case 5:
					// User plays "paper"
					userMove = "paper"
					if myMove == "scissors" {
						win = 1
					} else if myMove == "rock" {
						win = -1
					}
					break
				default:
					answer = "Sorry... I don't get it"
					sdk_wrapper.DisplayImage(fName, 5000, true)
					_ = os.Rename(fName, "/tmp/not_recognized_"+string(time.Now().Unix())+".jpg")
					break
				}

				if answer == "" {
					answer = "You put " + userMove + ". "

					switch win {
					case -1:
						answer = answer + "You win!"
						userScore++
						break
					case 1:
						answer = answer + "I win!"
						myScore++
						break
					default:
						answer = answer + "It's a draw!"
						break
					}
				}
				sdk_wrapper.SayText("I put " + myMove + "!")
				sdk_wrapper.SayText(answer)
				sdk_wrapper.WriteText(fmt.Sprintf("%d - %d", myScore, userScore), 64, true, 5000, true)
			} else {
				println("OPENCV Python script not found!")
			}
		}
	}
}
