package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("Listener is active")
	url := "http://" + os.Getenv("SPEAKER_URL") + ":3000"
	log.Println("Got Speaker URL", url)
	for {
		log.Println("Sending a request")
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Got error: ", err)
		} else {
			log.Println("Got response: ", resp.Body)
			defer resp.Body.Close()
		}

		time.Sleep(1 * time.Second)
	}
}
