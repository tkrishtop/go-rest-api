package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func speak(w http.ResponseWriter, r *http.Request) {
	speakerName := getEnv("SPEAKER_NAME", "unknown")
	speakerSpeech := getEnv("SPEAKER_SPEECH", "UnknownSpeech")

	log.Println(speakerName, "got a request, going to tell:", speakerSpeech)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "[%s] said: %s \n", speakerName, speakerSpeech)
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
