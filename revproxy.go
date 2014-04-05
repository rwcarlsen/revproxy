package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var target = flag.String("target", "", "url to forward to")
var sub = flag.String("subdomain", "", "subdomain that should be forwarded")
var addr = flag.String("addr", "", "address to listen on")

var proxys = map[string]*httputil.ReverseProxy{}

func main() {
	flag.Parse()
	remote, err := url.Parse("http://" + *target)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxys[*sub] = proxy

	http.HandleFunc("/", handler)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	subdomain := strings.Split(r.Host, ".")[0]
	if proxy, ok := proxys[subdomain]; ok {
		proxy.ServeHTTP(w, r)
	}
}
