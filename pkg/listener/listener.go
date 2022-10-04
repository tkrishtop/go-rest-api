package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func callUrl(url string) string {
	log.Println("[listener] Calling URL:", url)

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(responseData)
}

func gatherRequests(w http.ResponseWriter, r *http.Request) {
	speakers := []string{os.Getenv("WINNIE_URL"), os.Getenv("PIGLET_URL")}
	ports := []string{os.Getenv("WINNIE_PORT"), os.Getenv("PIGLET_PORT")}
	log.Println("[listener] Got a list of URLs", speakers)

	concatenatedResponse := ""
	for idx, urlString := range speakers {
		concatenatedResponse += callUrl("http://" + urlString + ports[idx])
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Concatenated replies: \n%s \n", concatenatedResponse)
}

func handleRequests() {
	http.HandleFunc("/speakers", gatherRequests)
	http.ListenAndServe(":3003", nil)
}

func main() {
	log.Println("Listener is active")
	handleRequests()
}
