package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

type Web1 struct {
}

func (w *Web1) GetIp(r *http.Request) string {
	ips := r.Header.Get("x-forwarded-for")
	if ips != "" {
		ips_list := strings.Split(ips, ",")
		if len(ips_list) > 0 && ips_list[0] != "" {
			return ips_list[0]
		}
	}
	return r.RemoteAddr
}

func (web *Web1) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	auth := request.Header.Get("Authorization")

	if auth == "" {
		writer.Header().Set("WWW-Authenticate", "Basic realm='please input your account an password.'")
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	auth_list := strings.Split(auth, " ")
	if len(auth_list) == 2 && auth_list[0] == "Basic" {
		res, err := base64.StdEncoding.DecodeString(auth_list[1])
		log.Printf(string(res))
		if err == nil && string(res) == "sx:123" {
			writer.Write([]byte("<h1>we1：</h1>" + web.GetIp(request)))
			return
		}
	}

	writer.Write([]byte("hello world web1!" + request.RequestURI))
}

type Web2 struct {
}

func (w *Web2) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hello world web2!" + request.RequestURI))
}

type Web3 struct {
}

func (w *Web3) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hello world web333!" + request.RequestURI))
}

func main() {
	c := make(chan os.Signal)
	//go func() {
	//	http.ListenAndServe("127.0.0.1:9091", &Web1{})
	//}()

	//go func() {
	//	http.ListenAndServe("127.0.0.1:9092", &Web2{})
	//}()

	go func() {
		http.ListenAndServe("127.0.0.1:9093", &Web3{})
	}()

	signal.Notify(c, os.Interrupt)
	s := <-c
	log.Println(s)
}
