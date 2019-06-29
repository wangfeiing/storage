package main

import (
	"net/http"
	"storage/objects"
	"log"
	"storage/config"
)

func main() {
	http.HandleFunc("/objects/" , objects.Handler)
	log.Fatal(http.ListenAndServe( config.LISTEN_ADDRESS , nil))
}
