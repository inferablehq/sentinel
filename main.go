package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/inferablehq/sentinel/pkg/strategies"
)

var config strategies.Config

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func main() {
	config = getConfig()

	http.HandleFunc("/", router)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func router(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	var found *strategies.Route

	for _, route := range config.Routes {
		if r.Method == route.Method && route.PathMatcher(r.URL.Path) {
			log.Printf("Processing request: %s %s", r.Method, r.URL.Path)
			found = &route
			break
		}
	}

	if found == nil {
		log.Printf("Request not found: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
		return
	}

	targetURL := "https://api.inferable.ai" + r.URL.Path
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	// Read the request body
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	// Process the request body
	processedRequestBody, err := found.RequestBodyHandler(string(requestBody))
	if err != nil {
		log.Printf("Error handling request body: %v", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	// Create the proxy request
	proxyReq, err := http.NewRequest(r.Method, targetURL, io.NopCloser(strings.NewReader(processedRequestBody)))
	if err != nil {
		log.Printf("Error creating proxy request: %v", err)
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request to the proxy request
	copyHeader(proxyReq.Header, r.Header)

	// Set the Content-Length header
	proxyReq.ContentLength = int64(len(processedRequestBody))
	proxyReq.Header.Set("Content-Length", fmt.Sprintf("%d", proxyReq.ContentLength))

	proxyClient := &http.Client{}
	proxyResp, err := proxyClient.Do(proxyReq)
	if err != nil {
		log.Printf("Error proxying request: %v", err)
		http.Error(w, "Error proxying request", http.StatusInternalServerError)
		return
	}
	defer proxyResp.Body.Close()

	copyHeader(w.Header(), proxyResp.Header)
	w.WriteHeader(proxyResp.StatusCode)

	responseBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	// Process the response body
	processedResponseBody, err := found.ResponseBodyHandler(string(responseBody))
	if err != nil {
		log.Printf("Error processing response body: %v", err)
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(processedResponseBody))
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}
