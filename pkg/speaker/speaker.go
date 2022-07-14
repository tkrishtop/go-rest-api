package main

import (
	"fmt"
	"net/http"
)

func speak(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "[Speaker] Hi there!")
}

func handleRequests() {
	http.HandleFunc("/", speak)
	http.ListenAndServe(":3000", nil)
}

func main() {
	handleRequests()
}
