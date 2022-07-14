package main

import (
	"fmt"
	"log"
	"net/http"
)

func speak(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request, going to speak")
	fmt.Fprintln(w, "[Speaker] Hi there!")
}

func handleRequests() {
	http.HandleFunc("/", speak)
	http.ListenAndServe(":3000", nil)
}

func main() {
	log.Println("Speaker is active")
	handleRequests()
}
