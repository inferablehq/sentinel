package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/inferablehq/sentinel/pkg/tokenizer"
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
			{
				Method: "POST",
				PathMatcher: func(path string) bool {
					return regexp.MustCompile(`^/clusters/\w+/runs$`).MatchString(path)
				},
				RequestBodyHandler: defaultPassthroughHandler,
			},
			{
				Method:             "GET",
				PathMatcher:        func(path string) bool { return regexp.MustCompile(`^/clusters/\w+/runs/\w+$`).MatchString(path) },
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

func maskedHandler(w http.ResponseWriter, r *http.Request, requestBody *string) {
	var unmarshalled map[string]interface{}
	err := json.Unmarshal([]byte(*requestBody), &unmarshalled)
	if err != nil {
		http.Error(w, "Error unmarshalling request body", http.StatusInternalServerError)
		return
	}

	masked, err := tokenizer.Tokenizer(unmarshalled, "", []string{})
	if err != nil {
		http.Error(w, "Error tokenizing request body", http.StatusInternalServerError)
		return
	}

	marshalled, err := json.Marshal(masked)
	if err != nil {
		http.Error(w, "Error marshalling tokenized request body", http.StatusInternalServerError)
		return
	}
	*requestBody = string(marshalled)
	w.Write(marshalled)
}
