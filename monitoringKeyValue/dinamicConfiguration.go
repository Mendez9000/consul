package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"net/http"
)

var c1_IP = "172.17.0.1"

func main() {
	loadConfiguration()
	http.HandleFunc("/watch-conf", watchConfHandler)

	log.Fatal(http.ListenAndServe(":8091", nil))
}

func loadConfiguration() {
	config := consulapi.DefaultConfig()
	config.Address = c1_IP + ":8500"
	consul, err := consulapi.NewClient(config)

	if err != nil {
		fmt.Println(err)
		return
	}

	kv := consul.KV()
	kvp, _, err1 := kv.Get("app1/config/secureModeLevel", nil)
	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println(string(kvp.Value))
	}
}

func watchConfHandler(w http.ResponseWriter, r *http.Request) {
	loadConfiguration()
}
