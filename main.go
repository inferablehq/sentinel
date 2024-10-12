package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var config Config

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

	var requestBody string
	processable := false

	for _, route := range config.Routes {
		if r.Method == route.Method && route.PathMatcher(r.URL.Path) {
			log.Printf("Processing request: %s %s", r.Method, r.URL.Path)
			route.RequestBodyHandler(w, r, &requestBody)
			processable = true
			break
		}
	}

	if !processable {
		log.Printf("Request not found: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
		return
	}

	targetURL := "https://api.inferable.ai" + r.URL.Path
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	proxyReq, err := http.NewRequest(r.Method, targetURL, strings.NewReader(requestBody))
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = r.Header

	proxyClient := &http.Client{}
	proxyResp, err := proxyClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error proxying request", http.StatusInternalServerError)
		return
	}
	defer proxyResp.Body.Close()

	copyHeader(w.Header(), proxyResp.Header)
	w.WriteHeader(proxyResp.StatusCode)
	io.Copy(w, proxyResp.Body)
	fmt.Println("requestBody")

	w.Write([]byte(requestBody))
}
