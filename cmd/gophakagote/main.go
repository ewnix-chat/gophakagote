package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Doist/unfurlist"
)

const (
	port = 8080
)

func allowCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "https://dev.chat.ewnix.net/" || origin == "https://chat.ewnix.net/" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
}

func embedURLHandler(fetcher http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowCORS(w, r)

		if r.Method != http.MethodGet {
			errMsg := "Invalid request method"
			log.Println("Error:", errMsg)
			http.Error(w, errMsg, http.StatusMethodNotAllowed)
			return
		}

		queryParams := r.URL.Query()
		url := queryParams.Get("url")
		if url == "" {
			errMsg := "URL parameter is required"
			log.Println("Error:", errMsg)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		r.URL.Path = fmt.Sprintf("/fetch/%s", url)
		fetcher.ServeHTTP(w, r)
	}
}

func main() {
	timeout := 5 * time.Second
	client := &http.Client{Timeout: timeout}

	fetcher := unfurlist.New(unfurlist.WithHTTPClient(client))

	http.HandleFunc("/embed", embedURLHandler(fetcher))
	log.Printf("Server listening on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
}

