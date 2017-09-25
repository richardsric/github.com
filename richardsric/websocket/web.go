package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/test", test)
	http.ListenAndServe(":3333", nil)
}

func test(w http.ResponseWriter, r *http.Request) {
	URL := "wss://api.poloniex.com"
	var dialer *websocket.Dialer

	_, _, err := dialer.Dial(URL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("connected succesfully")
}

func index(w http.ResponseWriter, r *http.Request) {
	URL := "wss://api.poloniex.com"
	var dialer *websocket.Dialer

	_, _, err := dialer.Dial(URL, nil)
	if err == websocket.ErrBadHandshake {
		log.Printf("handshake failed with status %d", w.Response.Status)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("connected succesfully")
}
