package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/inferablehq/inferable-go"
)

func TestLiveEndpoint(t *testing.T) {
	// Assuming the server is running on localhost:8080
	resp, err := http.Get("http://localhost:8080/live")
	if err != nil {
		t.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}
}

const (
	clusterId      = "01J7M489FSEY357S99YZ1H176E"
	consume_secret = "sk_cluster_consume_tmpjNA8YS3ElGjyGsxUiMrHZx1JXgL69MC5LCnBaKKE"
)

func TestInferableSetup(t *testing.T) {
	os.Setenv("STRATEGY", "e2e")
	os.Setenv("DATA_DIR", "./data")
	go main()

	// Initialize Inferable client
	client, err := inferable.New(inferable.InferableOptions{
		APIEndpoint: "http://localhost:8080",
		APISecret:   "sk_cluster_machine_eM0eP4cvmGdfdxLG3IIY9EJmpSQVpFei0EjBZscR8nA",
	})
	if err != nil {
		t.Fatalf("Failed to initialize Inferable client: %v", err)
	}

	type calcInput struct {
		Num1 int    `json:"num1"`
		Num2 int    `json:"num2"`
		Op   string `json:"op"`
	}

	err = client.Default.RegisterFunc(inferable.Function{
		Name:        "testFunction",
		Description: "A test function. Op can be add, sub, mul, or div.",
		Func: func(input calcInput) int {
			switch input.Op {
			case "add":
				return input.Num1 + input.Num2 + 10
			case "sub":
				return input.Num1 - input.Num2
			case "mul":
				return input.Num1 * input.Num2
			case "div":
				return input.Num1 / input.Num2
			}
			return 0
		},
	})

	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	err = client.Default.Start()
	defer client.Default.Stop()
	if err != nil {
		t.Fatalf("Failed to start client: %v", err)
	}

	if err != nil {
		t.Fatalf("Failed to register function: %v", err)
	}

	// Create a run to sum 1 and 2
	runID, err := createRun(consume_secret)

	if err != nil {
		t.Fatalf("Failed to create run: %v", err)
	}

	// Retrieve run results
	result, err := getRunResults(consume_secret, runID, 0)
	if err != nil {
		t.Fatalf("Failed to get run results: %v", err)
	}

	// Check if the result is correct
	if result != "13" {
		t.Errorf("Expected sum to be 13, got %s", result)
	}
}

func createRun(apiKey string) (string, error) {
	url := fmt.Sprintf("http://localhost:8080/clusters/%s/runs", clusterId)
	payload := strings.NewReader(`{
		"message": "Sum 1 and 2. Get the masked value.",
		"result": {
			"schema": {
				"type": "object",
				"properties": {
					"maskedValue": { "type": "string" }
				}
			}
		}
	}`)

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("Failed to create run: %v", res.Status)
	}

	body, _ := ioutil.ReadAll(res.Body)
	var response struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	fmt.Printf("Run created: %s", string(body))

	return response.ID, nil
}

func getRunResults(apiKey, runID string, attempts int) (string, error) {
	url := fmt.Sprintf("http://localhost:8080/clusters/%s/runs/%s", clusterId, runID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("authorization", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var response struct {
		Status string `json:"status"`
		Result struct {
			MaskedValue int `json:"maskedValue"`
		} `json:"result"`
	}
	fmt.Printf("==> Run status: %s\n", string(body))
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Status == "done" {
		return fmt.Sprintf("%d", response.Result.MaskedValue), nil
	} else if response.Status == "failed" {
		return "", fmt.Errorf("Run failed")
	} else {
		if attempts > 20 {
			return "", fmt.Errorf("Run not completed after 20 attempts")
		}

		time.Sleep(1 * time.Second)
		return getRunResults(apiKey, runID, attempts+1)
	}
}
