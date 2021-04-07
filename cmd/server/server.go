package main

import (
	"github.com/alexsuslov/godotenv"
	"github.com/srvchain/hooker"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var httpAddr = "0.0.0.0:80"

type Clients map[int]string

func (Clients Clients) IsAllow(id int) bool {
	_, ok := Clients[id]
	return ok
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warrning load env", err)
	}

	data, err := ioutil.ReadFile("clients.yml")
	if err != nil {
		panic(err)
	}
	clients := &Clients{}
	err = yaml.Unmarshal(data, clients)
	if err != nil {
		panic(err)
	}

	//get client
	http.HandleFunc("/", handle)
	hook:=hooker.Handle(clients)
	http.HandleFunc("/hook", hook)

	if os.Getenv("LISTEN") != "" {
		httpAddr = os.Getenv("LISTEN")
	}
	log.Println("listen", httpAddr)
	server := http.Server{Addr: httpAddr}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	log.Printf("body=%v", string(data))
	w.WriteHeader(200)
	w.Write([]byte("ok"))
	return
}
