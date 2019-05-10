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

/*
docker exec -ti c1 /bin/sh
apk add git
apk add go
apk add libc-dev
go get github.com/Mendez9000/consul/discoveryServiceNode
cd ~/go/src/github.com/Mendez9000/consul/discoveryServiceNode
go run discoveryServiceNode.go



docker exec -ti c2 /bin/sh
apk add git
apk add go
apk add libc-dev
go get github.com/Mendez9000/consul/discoveryServiceNode
cd ~/go/src/github.com/Mendez9000/consul/discoveryServiceNode
go run discoveryServiceNode.go

docker exec -ti c3 /bin/sh
apk add git
apk add go
apk add libc-dev
go get github.com/Mendez9000/consul/discoveryServiceNode
cd ~/go/src/github.com/Mendez9000/consul/discoveryServiceNode
go run discoveryServiceNode.go



git pull
go run discoveryServiceNode.go


*/

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
