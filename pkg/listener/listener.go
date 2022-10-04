package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func callUrl(url string) string {
	log.Println("[listener] Calling URL:", url)

	response, err := http.Get(url)
	if err != nil {
		log.Println("[listener] There is an error to call url, ignore it", err)
		return ""
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("[listener] There is an error convert response data, ignore it", err)
		return ""
	}

	return string(responseData)
}

// func parallelCall(ctx, url, wg) {
//     ctx, cancel := context.WithTimeout(ctx, time.Second)
//     defer wg.Done()
//     defer cancel()

//     callUrl(ctx, url)
// }

func gatherRequests(w http.ResponseWriter, r *http.Request) {
	speakers := []string{os.Getenv("WINNIE_URL"), os.Getenv("PIGLET_URL")}
	ports := []string{os.Getenv("WINNIE_PORT"), os.Getenv("PIGLET_PORT")}
	log.Println("[listener] Got a list of URLs", speakers)

	// collect replies from all speakers
	wg := sync.WaitGroup{}
	wg.Add(len(speakers))
	concatenatedResponse := ""

	for idx, urlBase := range speakers {
		urlString := "http://" + urlBase + ports[idx]
		go func(s string) {
			defer wg.Done()
			concatenatedResponse += callUrl(s)
		}(urlString)
	}

	wg.Wait()
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
