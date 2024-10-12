package main

import (
	"io"
	"net/http"
	"regexp"
)

type Route struct {
	Method             string
	PathMatcher        func(path string) bool
	RequestBodyHandler func(w http.ResponseWriter, r *http.Request, requestBody *string)
}

type Config struct {
	Routes []Route
}

func getConfig() Config {
	return Config{
		Routes: []Route{
			{
				Method:             "GET",
				PathMatcher:        func(path string) bool { return path == "/live" },
				RequestBodyHandler: defaultPassthroughHandler,
			},
			{
				Method:             "POST",
				PathMatcher:        func(path string) bool { return path == "/machines" },
				RequestBodyHandler: defaultPassthroughHandler,
			},
			{
				Method: "GET",
				PathMatcher: func(path string) bool {
					return regexp.MustCompile(`^/clusters/\w+/calls$`).MatchString(path)
				},
				RequestBodyHandler: defaultPassthroughHandler,
			},
			{
				Method: "POST",
				PathMatcher: func(path string) bool {
					return regexp.MustCompile(`^/clusters/\w+/calls/\w+/result$`).MatchString(path)
				},
				RequestBodyHandler: defaultPassthroughHandler,
			},
		},
	}
}

func defaultPassthroughHandler(w http.ResponseWriter, r *http.Request, requestBody *string) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	b := string(bodyBytes)
	*requestBody = b
}

// func defaultPassthroughHandler(w http.ResponseWriter, r *http.Request, requestBody *string) {
// 	targetURL := "https://api.inferable.ai" + r.URL.Path
// 	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
// 	if err != nil {
// 		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
// 		return
// 	}

// 	// Copy headers from the original request to the proxy request
// 	for header, values := range r.Header {
// 		for _, value := range values {
// 			proxyReq.Header.Add(header, value)
// 		}
// 	}

// 	// Send the proxy request
// 	client := &http.Client{}
// 	resp, err := client.Do(proxyReq)
// 	if err != nil {
// 		http.Error(w, "Error forwarding request", http.StatusBadGateway)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Copy the response headers back to the client
// 	for header, values := range resp.Header {
// 		for _, value := range values {
// 			w.Header().Add(header, value)
// 		}
// 	}

// 	// Set the status code
// 	w.WriteHeader(resp.StatusCode)

// 	// Copy the response body back to the client
// 	_, err = io.Copy(w, resp.Body)
// 	if err != nil {
// 		fmt.Printf("Error copying response body: %s\n", err)
// 	}
// }
