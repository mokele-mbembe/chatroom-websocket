package main

import (
	"flag"
	"log"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/home.html")
}

func main() {
	port := flag.String("port", "80", "http service port")
	flag.Parse() //parse variable
	hub := NewHub()
	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/chat", func(rw http.ResponseWriter, r *http.Request) {
		ServeWs(hub, rw, r)
	})
	log.Printf("http serve on port %s\n", *port)
	if err := http.ListenAndServe("127.0.0.1:"+*port, nil); err != nil { //this handler is the same as "/"
		log.Printf("start http service error: %s\n", err)
	}
}
