package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const FlathubServer = "https://dl.flathub.org"
const ProxyServer = "https://chn.flathub.cf"

func server() {
	server, _ := url.Parse(FlathubServer)
	proxy := httputil.NewSingleHostReverseProxy(server)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = server.Host
		r.URL.Scheme = server.Scheme
		r.Host = server.Host
		log.Println("proxy", r.URL.String())

		proxy.ServeHTTP(w, r)
	})
	http.HandleFunc("/repo/summary.sig", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	http.HandleFunc("/repo/summary", func(w http.ResponseWriter, r *http.Request) {
		log.Println("summary")
		r.URL.Host = server.Host
		r.URL.Scheme = server.Scheme
		resp, err := http.Get(r.URL.String())
		if err != nil {
			log.Println("error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if resp.StatusCode >= 400 {
			w.WriteHeader(resp.StatusCode)
			return
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data = bytes.Replace(data, []byte(FlathubServer+"/repo/"), []byte(ProxyServer+"/repo/"), 1)
		w.Write(data)
	})
	log.Panic(http.ListenAndServe(":18080", nil))
}

func main() {
	server()
}
