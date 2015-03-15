package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func handler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprint(rw, "Hi there, i love %S!", req.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
