package telemetry

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Telemetry struct {
	apiKey  string
	baseUrl string
}

// NewTelemetry creates a new Telemetry instance.
func NewTelemetry() *Telemetry {
	return &Telemetry{
		baseUrl: "https://api.telemetry.sh",
	}
}

// Init initializes the Telemetry instance with the given API key.
func (t *Telemetry) Init(apiKey string) {
	t.apiKey = apiKey
}

// Log sends data to the specified table.
func (t *Telemetry) Log(table string, data map[string]interface{}) (map[string]interface{}, error) {
	if t.apiKey == "" {
		return nil, errors.New("API key is not initialized. Please call Init() with your API key.")
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": t.apiKey,
	}

	body := map[string]interface{}{
		"data":  data,
		"table": table,
	}

	return t.makeRequest("/log", headers, body)
}

// Query sends a query to the Telemetry service.
func (t *Telemetry) Query(query string) (map[string]interface{}, error) {
	if t.apiKey == "" {
		return nil, errors.New("API key is not initialized. Please call Init() with your API key.")
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": t.apiKey,
	}

	body := map[string]interface{}{
		"query":    query,
		"realtime": true,
		"json":     true,
	}

	return t.makeRequest("/query", headers, body)
}

// makeRequest is a helper function to make HTTP requests.
func (t *Telemetry) makeRequest(endpoint string, headers map[string]string, body map[string]interface{}) (map[string]interface{}, error) {
	url := t.baseUrl + endpoint

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	return responseData, nil
}
