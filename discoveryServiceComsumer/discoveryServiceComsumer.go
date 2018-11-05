package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Dentro de la maquina local (no docker)
// dig @172.17.0.2 -p 8600 orders.service.dc1.consul SRV
// dig +noall +answer -t aaaa @172.17.0.2 -p 8600 orders.service.dc1.consul SRV

func main() {
	resp, _ := http.Get("http://localhost:8500/v1/catalog/service/orders")
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var address []Address
	err := json.Unmarshal(bodyBytes, &address)

	fmt.Print(err)
	fmt.Print(string(bodyBytes))
}

type Address struct {
	ServiceAddress string
	ServicePort    uint16
}
