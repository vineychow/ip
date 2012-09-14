package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	requestArgName = "ip"
)

var (
	mux     sync.Mutex
	service IpService
)

func init() {
	mux.Lock()
	defer mux.Unlock()

	if service != nil {
		return
	}

	// new
	newService, err := NewIpService("ip.route")
	if err != nil {
		panic("new ip service fail")
	}

	// success
	service = newService
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue(requestArgName)
	if len(ip) == 0 {
		w.Write([]byte("ip is null"))
		return
	}

	// query
	resp, err := service.Query(ip)
	if err != nil {
		println(err.Error())
		w.Write([]byte("service error"))
	}

	w.Write([]byte(resp))
}

func main() {
	fmt.Println("query ip start")

	http.HandleFunc("/query", ipHandler)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println("query ip end")
}
