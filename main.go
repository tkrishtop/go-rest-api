package main

import (
	"fmt"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "[Worker 1] Hi there!")
}

func handleRequests() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":3000", nil)
}

func main() {
	handleRequests()
}
