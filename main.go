package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// global links variable
var links map[string]string

// Default HTTP handler. Looks up the path in the links map. 
// If it gets a result redirects, else returns a 404
func home(w http.ResponseWriter, r *http.Request) {
	// Strip any leading slashes from path
	strippedPath := strings.TrimPrefix(r.URL.Path, "/")

	// Get result from link map
	urlResult := links[strippedPath]

	slog.Info("Got request", slog.String("path", r.URL.Path), slog.String("urlResult", urlResult))

	// Check if we got a default value from the map
	if urlResult == "" {
		slog.Info("Returning Not Found")
		http.NotFound(w, r)
	} else {
		slog.Info("Redirecting")
		http.Redirect(w, r, urlResult, http.StatusFound)
	}

}

func main() {
	// Read flags
	addr := flag.String("addr", "localhost:8080", "What host:port to listen on")
	source := flag.String("src", "links.json", "JSON file for links")
	flag.Parse()

	// load links
	file, err := os.Open(*source)

	if err != nil {
		slog.Error("Failed to read links", "err", err)
		os.Exit(1)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&links)

	if err != nil {
		slog.Error("Failed to decode file", "err", err)
		os.Exit(1)
	}

	// Setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	slog.Info("Starting server..", slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, mux)
	slog.Error("Failed to start", "err", err)
}
