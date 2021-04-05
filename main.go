package main

import (
    "fmt"
    "net/http"
)

func home(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Hi there!")
}

func handleRequests() {
    http.HandleFunc("/", home)
    http.ListenAndServe(":3000", nil)
}

func main() {
    handleRequests()
}