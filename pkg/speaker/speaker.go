package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Speaker struct {
	Name   string
	Speech string
	Delay  string
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func speak(w http.ResponseWriter, r *http.Request) {
	s := Speaker{
		Name:   getEnv("SPEAKER_NAME", "unknown"),
		Speech: getEnv("SPEAKER_SPEECH", "UnknownSpeech"),
		Delay:  getEnv("SPEAKER_DELAY", "0"),
	}

	log.Println(s.Name, "( delay in seconds:", s.Delay, ") got a request, going to tell:", s.Speech)

	// this sleep is to simulate timeout
	intDelay, _ := strconv.Atoi(s.Delay)
	time.Sleep(time.Duration(intDelay) * time.Second)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "[%s] said: %s \n", s.Name, s.Speech)
}

func handleRequests() {
	speakerPort := getEnv("SPEAKER_PORT", ":3000")
	http.HandleFunc("/", speak)
	http.ListenAndServe(speakerPort, nil)
}

func main() {
	speakerName := getEnv("SPEAKER_NAME", "unknown")
	log.Println(speakerName, "is active")
	handleRequests()
}
