package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"tailscale.com/tsnet"
)

var (
	hostname   = flag.String("hostname", "hello", "hostname for the tailnet")
	targetHost = flag.String("target", "crusty", "target hostname to proxy requests to")
)

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	flag.Parse()

	srv := new(tsnet.Server)
	srv.Hostname = *hostname

	defer srv.Close()

	ln, err := srv.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}

	proxy, err := NewProxy(*targetHost)
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	log.Fatal(http.Serve(ln, http.HandlerFunc(ProxyRequestHandler(proxy))))
}
