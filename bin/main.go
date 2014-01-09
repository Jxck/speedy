package main

import (
	"flag"
	"fmt"
	"github.com/Jxck/speedy"
	"log"
	"net/http"
)

// define your http.Handler struct
type Hello struct{}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

var httpHandler http.Handler = &Hello{}

// var httpHandler http.Handler = http.FileServer(http.Dir("."))

// your cert/key path
var (
	cert = "keys/cert.pem"
	key  = "keys/key.pem"
)

// go run bin/main.go :3000
// or debug mode
// DEBUG=DEBUG go run main.go :3000
func main() {
	flag.Parse()
	port := flag.Args()[0]

	// start speedy server
	err := speedy.ListenAndServe(port, cert, key, httpHandler)
	if err != nil {
		log.Fatal(err)
	}
}
