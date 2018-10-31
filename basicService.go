package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var timestamp string
var webClient *http.Client
var port int

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, timestamp)
}

func main() {
	port = 8080
	timestamp = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	http.HandleFunc("/", handler)
	localIp := "http://" + GetOutboundIP().String()
	httpCheck := HttpCheck{Http: localIp + ":" + strconv.Itoa(port), Method: "POST", Interval: "4s"}
	service := Service{ID: "1455", Name: "Orders", Port: port, Address: localIp, Check: httpCheck}
	doRegistration(service)
	defer doDeregister(service)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("doDeregister: %+v", sig)
		doDeregister(service)
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func doRegistration(service Service) {
	endpointRegister := "http://localhost:8500/v1/agent/service/register"

	d, _ := json.Marshal(service)

	client := &http.Client{}
	client.Timeout = time.Second * 15

	body := bytes.NewBuffer(d)
	req, _ := http.NewRequest(http.MethodPut, endpointRegister, body)

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, _ := client.Do(req)

	defer resp.Body.Close()
}

func doDeregister(service Service) {
	endpointRegister := "http://localhost:8500/v1/agent/service/deregister/" + service.ID
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPut, endpointRegister, nil)
	client.Do(req)
}

type Service struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Port    int       `json:"port"`
	Address string    `json:"address"`
	Check   HttpCheck `json:"check"`
}

type HttpCheck struct {
	Http     string `json:"http, omitempty"`
	Method   string `json:"method, omitempty"`
	Interval string `json:"interval, omitempty"`
}

/*
type HttpCheck struct {
	ID            string `json:"id, omitempty"`
	Name          string `json:"name, omitempty"`
	Http          string `json:"http, omitempty"`
	TlsSkipVerify bool   `json:"tls_skip_verify, omitempty"`
	Method        string `json:"method, omitempty"`
	Header        string `json:"header, omitempty"`
	Interval      string `json:"interval, omitempty"`
	Timeout       string `json:"timeout, omitempty"`
}
*/
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
