package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/inferablehq/sentinel/pkg/tokenizer"
)

type Route struct {
	Method              string
	PathMatcher         func(path string) bool
	RequestBodyHandler  func(in string) (string, error)
	ResponseBodyHandler func(in string) (string, error)
}

type Config struct {
	Routes []Route
}

func getConfig() Config {
	return Config{
		Routes: []Route{
			{
				Method:              "GET",
				PathMatcher:         func(path string) bool { return path == "/live" },
				RequestBodyHandler:  defaultPassthroughHandler,
				ResponseBodyHandler: defaultPassthroughHandler,
			},
			{
				Method:              "POST",
				PathMatcher:         func(path string) bool { return path == "/machines" },
				RequestBodyHandler:  defaultPassthroughHandler,
				ResponseBodyHandler: defaultPassthroughHandler,
			},
			{
				Method: "GET",
				PathMatcher: func(path string) bool {
					return regexp.MustCompile(`^/clusters/\w+/calls$`).MatchString(path)
				},
				RequestBodyHandler:  defaultPassthroughHandler,
				ResponseBodyHandler: detokenizableDataHandler,
			},
			{
				Method: "POST",
				PathMatcher: func(path string) bool {
					return regexp.MustCompile(`^/clusters/\w+/calls/\w+/result$`).MatchString(path)
				},
				RequestBodyHandler:  func(in string) (string, error) { return tokenizableDataHandler(in, []string{".resultType"}) },
				ResponseBodyHandler: defaultPassthroughHandler,
			},
			{
				Method: "POST",
				PathMatcher: func(path string) bool {
					return regexp.MustCompile(`^/clusters/\w+/runs$`).MatchString(path)
				},
				RequestBodyHandler: func(in string) (string, error) {
					return tokenizableDataHandler(in, []string{".result.schema", ".message"})
				},
				ResponseBodyHandler: defaultPassthroughHandler,
			},
			{
				Method:              "GET",
				PathMatcher:         func(path string) bool { return regexp.MustCompile(`^/clusters/\w+/runs/\w+$`).MatchString(path) },
				RequestBodyHandler:  defaultPassthroughHandler,
				ResponseBodyHandler: detokenizableDataHandler,
			},
		},
	}
}

func defaultPassthroughHandler(in string) (string, error) {
	return in, nil
}

func tokenizableDataHandler(in string, except []string) (string, error) {
	var unmarshalled interface{}
	err := json.Unmarshal([]byte(in), &unmarshalled)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		return "", fmt.Errorf("error unmarshalling request body: %v", err)
	}

	masked, err := tokenizer.Tokenizer(unmarshalled, "", except) // Pass except here
	if err != nil {
		log.Printf("Error tokenizing request body: %v", err)
		return "", fmt.Errorf("error tokenizing request body: %v", err)
	}

	marshalled, err := json.Marshal(masked)
	if err != nil {
		log.Printf("Error marshalling tokenized request body: %v", err)
		return "", fmt.Errorf("error marshalling tokenized request body: %v", err)
	}

	log.Printf("Tokenized request body: %s", string(marshalled))

	return string(marshalled), nil
}

func detokenizableDataHandler(in string) (string, error) {
	log.Printf("Detokenizing response body: %s", in)

	var unmarshalled interface{}
	err := json.Unmarshal([]byte(in), &unmarshalled)
	if err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		return "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	detokenized, err := tokenizer.Detokenizer(unmarshalled, "", []string{}) // Pass an empty except list here
	if err != nil {
		log.Printf("Error detokenizing response body: %v", err)
		return "", fmt.Errorf("error detokenizing response body: %v", err)
	}

	marshalled, err := json.Marshal(detokenized)
	if err != nil {
		log.Printf("Error marshalling detokenized response body: %v", err)
		return "", fmt.Errorf("error marshalling detokenized response body: %v", err)
	}

	log.Printf("Detokenized response body: %s", string(marshalled))

	return string(marshalled), nil
}
