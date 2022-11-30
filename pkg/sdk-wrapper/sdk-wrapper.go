package sdk_wrapper

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vector"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var Robot *vector.Vector
var bcAssumption bool = false
var ctx context.Context
var transCfg = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore SSL warnings
}
var eventStream vectorpb.ExternalInterface_EventStreamClient

func InitSDK(serial string) {
	var err error
	InitLanguages(LANGUAGE_ENGLISH)
	Robot, err = vector.NewEP(serial)
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.Background()
	eventStream, err = Robot.Conn.EventStream(ctx, &vectorpb.EventRequest{})
	RefreshSDKSettings()
}

func WaitForEvent() *vectorpb.Event {
	evt, err := eventStream.Recv()
	if err != nil {
		return nil
	}
	return evt.GetEvent()
}

func AssumeBehaviorControl(priority string) {
	var controlRequest *vectorpb.BehaviorControlRequest
	if priority == "high" {
		controlRequest = &vectorpb.BehaviorControlRequest{
			RequestType: &vectorpb.BehaviorControlRequest_ControlRequest{
				ControlRequest: &vectorpb.ControlRequest{
					Priority: vectorpb.ControlRequest_OVERRIDE_BEHAVIORS,
				},
			},
		}
	} else {
		controlRequest = &vectorpb.BehaviorControlRequest{
			RequestType: &vectorpb.BehaviorControlRequest_ControlRequest{
				ControlRequest: &vectorpb.ControlRequest{
					Priority: vectorpb.ControlRequest_DEFAULT,
				},
			},
		}
	}
	go func() {
		start := make(chan bool)
		stop := make(chan bool)
		bcAssumption = true
		go func() {
			// * begin - modified from official vector-go-sdk
			r, err := Robot.Conn.BehaviorControl(
				ctx,
			)
			if err != nil {
				log.Println(err)
				return
			}

			if err := r.Send(controlRequest); err != nil {
				log.Println(err)
				return
			}

			for {
				ctrlresp, err := r.Recv()
				if err != nil {
					log.Println(err)
					return
				}
				if ctrlresp.GetControlGrantedResponse() != nil {
					start <- true
					break
				}
			}

			for {
				select {
				case <-stop:
					if err := r.Send(
						&vectorpb.BehaviorControlRequest{
							RequestType: &vectorpb.BehaviorControlRequest_ControlRelease{
								ControlRelease: &vectorpb.ControlRelease{},
							},
						},
					); err != nil {
						log.Println(err)
						return
					}
					return
				default:
					continue
				}
			}
			// * end - modified from official vector-go-sdk
		}()
		for {
			select {
			case <-start:
				for {
					if bcAssumption {
						time.Sleep(time.Millisecond * 500)
					} else {
						break
					}
				}
				stop <- true
				return
			}
		}
	}()
}

func ReleaseBehaviorControl() {
	bcAssumption = false
}

func GetRobotSerial() string {
	return Robot.Cfg.SerialNo
}

/**********************************************************************************************************************/
/*                                              PRIVATE FUNCTIONS                                                     */
/**********************************************************************************************************************/

const shellToUse = "bash"

func shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(shellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func getTempFilename(extension string, fullpath bool) string {
	tmpPath := "/tmp/"
	tmpFile := GetRobotSerial() + fmt.Sprintf("_%d", time.Now().Unix()) + "." + extension
	if fullpath {
		tmpFile = tmpPath + tmpFile
	}
	return tmpFile
}

func getSDKSettings() []byte {
	resp, err := Robot.Conn.PullJdocs(ctx, &vectorpb.PullJdocsRequest{
		JdocTypes: []vectorpb.JdocType{vectorpb.JdocType_ROBOT_SETTINGS},
	})
	if err != nil {
		return []byte(err.Error())
	}
	json := resp.NamedJdocs[0].Doc.JsonDoc
	return []byte(json)
}
