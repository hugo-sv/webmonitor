package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// ServerDown returns a 500 Status Code
func ServerDown(w http.ResponseWriter) {
	http.Error(w, "The server is down.", http.StatusInternalServerError)
}

// ServerUp returns a 200 Status Code
func ServerUp(w http.ResponseWriter) {
	fmt.Fprintf(w, "The server is up.")
}

// Handler handles HTTP requests on this server
func Handler(w http.ResponseWriter, r *http.Request) {
	// Adding a random reponse time
	time.Sleep(time.Second * time.Duration(rand.Float64()))
	switch r.URL.Path[1:] {
	// 0% availability
	case "down":
		ServerDown(w)
		return
	// random based 80% availability
	case "random":
		if rand.Intn(100) >= 80 {
			ServerDown(w)
			return
		}
	// This route should availability over 2 min varies from 100% to 60%
	case "alert":
		if time.Now().Unix()%168 < 48 {
			ServerDown(w)
			return
		}
	}
	// By default, 100% availability
	ServerUp(w)
}

func main() {
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
