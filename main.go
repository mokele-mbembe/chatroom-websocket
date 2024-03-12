package main

import (
	"flag"
	"log"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/home.html")
}

// start service:
func main() {
	port := flag.String("port", "5678", "http service port")
	flag.Parse()

	http.HandleFunc("/", serveHome)
	log.Printf("http serve on port %s\n", *port)

}
