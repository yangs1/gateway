package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Web1 struct {
}

func (w *Web1) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hello world web1!"))
}

type Web2 struct {
}

func (w *Web2) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hello world web2!"))
}

type myproxy struct {
}

func (prx *myproxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			log.Println(err)
		}
	}()

	//config
	proxyConfig := map[string]string{"/web1": "127.0.0.1:9090"}
	if url, ok := proxyConfig[r.URL.Path]; ok {
		RequestUrl(w, r, url)
	}

	w.Write([]byte("default index"))
}

func RequestUrl(w http.ResponseWriter, r *http.Request, url string) {
	newr, _ := http.NewRequest(r.Method, "http://127.0.0.1:9090", r.Body)
	for hk, hv := range r.Header {
		newr.Header.Set(hk, hv[0])
	}

	newresponse, _ := http.DefaultClient.Do(newr)
	defer newresponse.Body.Close()
	res_content, _ := ioutil.ReadAll(newresponse.Body)

	for hk, hv := range newresponse.Header {
		w.Header().Set(hk, hv[0])
	}
	w.WriteHeader(newresponse.StatusCode)

	w.Write([]byte(res_content))
	return
}

func main() {
	c := make(chan os.Signal)
	go func() {
		http.ListenAndServe("127.0.0.1:8080", &myproxy{})
	}()

	signal.Notify(c, os.Interrupt)
	s := <-c
	log.Println(s)
}
