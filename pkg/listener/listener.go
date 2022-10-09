package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type Speaker struct {
	Name string
	Url  string
	Port string
}

func callUrl(request_ctx context.Context, url string) string {
	log.Println("Calling URL:", url)

	request, err := http.NewRequestWithContext(request_ctx, "GET", url, nil)
	if err != nil {
		return "Error while constructing request " + url + ": " + err.Error() + "\n"
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "Error while calling url " + url + ": " + err.Error() + "\n"
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "Error while reading response from " + url + ": " + err.Error() + "\n"
	}

	return string(responseData)
}

// This is a pipeline function: gather input -> data treatement -> send output
// The idea is to use a Done channel for cancellation: https://blog.golang.org/pipelines
func gatherRequests(ctx context.Context, speakers []Speaker, w http.ResponseWriter, r *http.Request) {
	// Set up a done channel that's shared by the whole pipeline,
	// and close that channel when this pipeline exits, as a signal
	// for all the goroutines we started to exit.
	request_ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// collect replies from all speakers
	log.Println("List of speakers to call: ", speakers)
	wg := sync.WaitGroup{}
	wg.Add(len(speakers))
	responseChan := make(chan string, len(speakers))

	for _, speaker := range speakers {
		log.Println("Calling speaker", speaker.Name)
		urlString := "http://" + speaker.Url + speaker.Port
		go func(s string) {
			defer wg.Done()
			select {
			case responseChan <- callUrl(request_ctx, s):
				log.Println("Received reply for url: ", s)
			case <-request_ctx.Done():
				timeoutMessage := "Deadline exceeded for url: " + s
				log.Println(timeoutMessage)
				responseChan <- timeoutMessage
			}
		}(urlString)
	}

	// build concatenated response
	wg.Wait()
	concatenatedResponse := ""
	for range speakers {
		concatenatedResponse += <-responseChan
	}

	// send the requests
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Concatenated replies: \n%s \n", concatenatedResponse)
}

func readConfig() []Speaker {
	speakersFilePath := "/config/" + os.Getenv("CONFIG_FILE")

	log.Println("Reading config file: ", speakersFilePath)
	rawData, err := ioutil.ReadFile(speakersFilePath)
	if err != nil {
		log.Fatal("Failed reading config file: ", err)
	}

	log.Println("Unmarchalling configuration")
	var speakers []Speaker
	err = yaml.Unmarshal([]byte(rawData), &speakers)
	if err != nil {
		log.Fatal("Failed unmarchalling: ", err)
	}

	return speakers
}

func main() {
	// initialize listener
	speakers := readConfig()
	log.Println("Listener is active")

	// handle requests
	ctx := context.Background()
	http.HandleFunc("/speakers", func(w http.ResponseWriter, r *http.Request) {
		gatherRequests(ctx, speakers, w, r)
	})
	http.ListenAndServe(":3003", nil)
}
