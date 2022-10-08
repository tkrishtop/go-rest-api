package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Speaker struct {
	Name string
	Url  string
	Port string
}

func callUrl(url string) string {
	log.Println("Calling URL:", url)

	response, err := http.Get(url)
	if err != nil {
		log.Println("There is an error while calling url, ignore it and return empty reply:", err)
		return ""
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("There is an error while converting response data, ignore it and return empty reply:", err)
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

func gatherRequests(speakers []Speaker) string {
	log.Println("List of speakers to call: ", speakers)

	// collect replies from all speakers
	wg := sync.WaitGroup{}
	wg.Add(len(speakers))
	responseChan := make(chan string, len(speakers))

	for _, speaker := range speakers {
		log.Println("Calling speaker", speaker.Name)
		urlString := "http://" + speaker.Url + speaker.Port
		go func(s string) {
			defer wg.Done()
			responseChan <- callUrl(s)
		}(urlString)
	}

	wg.Wait()
	// build concatenated response
	concatenatedResponse := ""
	for range speakers {
		concatenatedResponse += <-responseChan
	}
	return concatenatedResponse
}

func callSpeakers(w http.ResponseWriter, r *http.Request) {
	concatenatedResponse := ""

	speakersFilePath := "/config/" + os.Getenv("CONFIG_FILE")
	log.Println("Going to read config file: ", speakersFilePath)

	rawData, err := ioutil.ReadFile(speakersFilePath)
	if err != nil {
		log.Println("Failed reading config file: ", err, "Going to stop here.")
	} else {
		log.Println("Going to unmarchall configuration")

		var speakers []Speaker
		err = yaml.Unmarshal([]byte(rawData), &speakers)
		if err != nil {
			log.Println("Failed unmarchalling: ", err, "Going to stop here.")
		} else {
			concatenatedResponse = gatherRequests(speakers)
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Concatenated replies: \n%s \n", concatenatedResponse)
}

func handleRequests() {
	http.HandleFunc("/speakers", callSpeakers)
	http.ListenAndServe(":3003", nil)
}

func main() {
	log.Println("Listener is active")
	handleRequests()
}
