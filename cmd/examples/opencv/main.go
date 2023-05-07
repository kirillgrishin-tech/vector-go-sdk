package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"os/exec"
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
			//doOpenCVStuff(50)
			doOpenCVStuffWithImageServer()
			stop <- true
			return
		}
	}
}

func doOpenCVStuff(numSteps int) {
	for i := 0; i <= numSteps; i++ {
		fName := fmt.Sprintf("/tmp/camera%02d.jpg", i)
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
				fmt.Printf("FRAME %d, Output: %s\n", i, out.String())
				sdk_wrapper.SayText(out.String())
				time.Sleep(time.Duration(2000) * time.Millisecond)
			} else {
				println("OPENCV Python script not found!")
			}
		}
	}
}

func doOpenCVStuffWithImageServer() {
	var client *http.Client
	//setup a mocked http client.
	println("")
	println("Setup HTTP client")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", b)
	}))
	defer ts.Close()
	client = ts.Client()

	println("Starting Vector")
	sdk_wrapper.MoveHead(3.0)
	sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
	sdk_wrapper.UseVectorEyeColorInImages(true)
	sdk_wrapper.EnableCameraStream()
	for {
		image := sdk_wrapper.ProcessCameraStream()
		if image != nil {
			var handInfo map[string]interface{}
			jsonData := sendImageToImageServer(client, &image)
			json.Unmarshal([]byte(jsonData), &handInfo)
			text := "?"
			fingers := -1
			fingers = int(handInfo["raisedfingers"].(float64))
			fx := int(handInfo["index_x"].(float64))
			fy := int(handInfo["index_y"].(float64))
			if fingers >= 0 {
				text = fmt.Sprintf("%d", fingers)
			}
			text += fmt.Sprintf(" %d,%d", fx, fy)
			sdk_wrapper.WriteText(text, 32, true, 100, false)
		}
	}
}

func sendImageToImageServer(client *http.Client, img *image.Image) string {
	//println("Encoding new frame")
	// Convert image to jpg and obtain the bytes
	var imageBuf bytes.Buffer
	_ = jpeg.Encode(&imageBuf, *img, nil)

	// Prepare the reader instances to encode
	values := map[string]io.Reader{
		"file": bytes.NewReader(imageBuf.Bytes()),
	}

	// Upload and get back the json response
	resp, err := Upload(client, "http://192.168.43.65:8090", values)
	if err != nil {
		println("Response error!")
		return ""
	}

	// Return json string
	//println("Response received: " + resp)
	return resp
}

func Upload(client *http.Client, url string, values map[string]io.Reader) (response string, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return "", err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return "", err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return "", err
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	//println("Encoded data")
	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		println("Error when performing HTTP request")
		return "", err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	//println("POSTing...")
	res, err := client.Do(req)
	if err != nil {
		println(err.Error())
		return "", err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		println(fmt.Errorf("bad status: %s", res.Status))
	} else {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		return bodyString, nil
	}
	return "", err
}
