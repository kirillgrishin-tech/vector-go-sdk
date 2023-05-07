package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/digital-dream-labs/hugh/grpc/client"
	"github.com/kirillgrishin-tech/vector-go-sdk/pkg/vectorpb"
)

const (
	// These can be set to whatever you want
	sessionID  = "id01"
	clientName = "id02"
)

func main() {
	var robotName = flag.String("name", "", "Vector's name")
	var host = flag.String("host", "", "Vector's IP address")
	var serial = flag.String("serial", "", "Vector's serial number")
	var username = flag.String("username", "", "Anki account username")
	var password = flag.String("password", "", "Anki account password")
	flag.Parse()
	if *robotName == "" {
		log.Fatal("please use the -name argument and set it to your robot name")
	}
	if *host == "" {
		log.Fatal("please use the -host argument and set it to your robots IP address")
	}
	if *serial == "" {
		log.Fatal("please use the -serial argument and set it to your robots serial number")
	}
	if *username == "" {
		log.Fatal("please use the -username argument and set it to your anki account username")
	}
	if *password == "" {
		log.Fatal("please use the -password argument and set it to your anki account password")
	}

	var certFile = download_certificate(*serial, "./")
	var cert, _ = ioutil.ReadFile(certFile)
	var token = get_session_token(*username, *password)
	st := token["session"].(map[string]interface{})
	stt := []byte(fmt.Sprintf("%s", st["session_token"]))
	print("Session token: " + string(stt) + "\n\n")

	guid, err := user_authentication(stt, cert, *host, *robotName)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(guid)
}

func user_authentication(session_id []byte, cert []byte, ip string, name string) (string, error) {
	/*
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(cert) {
			log.Fatal("failed to add server CA's certificate")
		}*/

	c, err := client.New(
		client.WithTarget(
			fmt.Sprintf("%s:443", ip),
		),
		//client.WithCertPool(certPool),
		//client.WithOverrideServerName(name),
		client.WithInsecureSkipVerify(),
	)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	if err := c.Connect(); err != nil {
		log.Fatal(err)
		return "", err
	}

	vc := vectorpb.NewExternalInterfaceClient(c.Conn())

	response, err2 := vc.UserAuthentication(context.Background(),
		&vectorpb.UserAuthenticationRequest{
			UserSessionId: session_id,
			ClientName:    []byte(GetOutboundIP().String()),
		},
	)
	var guid string = string(response.GetClientTokenGuid())
	println("")
	println("GUID from robot: " + guid)

	return string(response.GetClientTokenGuid()), err2
}

func download_certificate(serial string, serialPath string) string {
	var fullURLFile = "https://session-certs.token.global.anki-services.com/vic/" + serial
	fileURL, err := url.Parse(fullURLFile)

	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	var fileName = segments[len(segments)-1]

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)
	defer file.Close()
	fmt.Printf("Downloaded a certificate file %s with size %d\n", fileName, size)

	return fileName
}

func get_session_token(username string, password string) map[string]interface{} {
	ret := ""
	api := "https://accounts.api.anki.com/1/sessions"

	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)

	req, err := http.NewRequest("POST", api, strings.NewReader(params.Encode()))
	req.Header.Set("User-Agent", "Vector-sdk/0.7.2")
	req.Header.Set("Anki-App-Key", "aung2ieCho3aiph7Een3Ei")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Get Session Token Response Status:\n", resp.Status)
	if resp.Status != "200 OK" {
		panic("Authentication error")
	} else {
		body, _ := io.ReadAll(resp.Body)
		ret = string(body)
	}
	print(ret)
	var jsonObj map[string]interface{}
	json.Unmarshal([]byte(ret), &jsonObj)
	return jsonObj
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
