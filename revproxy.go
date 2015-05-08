package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var addr = flag.String("addr", "", "address to listen on")
var target = flag.String("target", "", "address to forward to")

func main() {
	flag.Parse()
	remote, err := url.Parse("http://" + *target)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		fmt.Println("serving request ", req.URL)
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
	}

	http.Handle("/", proxy)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
