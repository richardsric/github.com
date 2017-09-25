package main

import (
	"io"
	"net/http"
	"fmt"
)

func root(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "YOU VISITED / route")

	fmt.Printf("Request URI is: %v\n", req.RequestURI)
}

func name(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "YOU VISITED /name route")
	fmt.Printf("Request URI is: %v\n", req.RequestURI)
}

func main(){
http.HandleFunc("/", root)
http.HandleFunc("/name", name)
http.ListenAndServe(":9292", nil)

}
