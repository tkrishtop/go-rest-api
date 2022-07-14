package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Listener is active")
	url := "http://192.168.39.29:30100"
	for {
		log.Println("Sending a request")
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Got error: ", err)
		}
		defer resp.Body.Close()

		log.Println("Got response: ", resp.Body)
		time.Sleep(1 * time.Second)
	}
}
