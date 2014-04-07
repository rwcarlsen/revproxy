package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var addr = flag.String("addr", "", "address to listen on")

var proxys = map[string]*httputil.ReverseProxy{}

// positional args are [<subdomain> <target-addr>]...
func main() {
	flag.Parse()

	for i := 0; i < len(flag.Args()); i += 2 {
		subdomain := flag.Arg(i)
		target := flag.Arg(i + 1)
		remote, err := url.Parse("http://" + target)
		if err != nil {
			log.Fatal(err)
		}
		proxys[subdomain] = httputil.NewSingleHostReverseProxy(remote)
	}

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(*addr, nil)
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
