package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const openAIModelsURL = "https://api.openai.com/v1/models"

// ValidateAPIKey validates the given API key by making a request to OpenAI API.
func ValidateAPIKey(apiKey string) error {
	if apiKey == "" {
		return errors.New("API key cannot be empty")
	}

	httpRequest, err := createRequestWithHeaders("GET", openAIModelsURL, apiKey)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer httpResponse.Body.Close()

	if err := handleHTTPResponse(httpResponse); err != nil {
		return err
	}

	return nil
}

func createRequestWithHeaders(method, url, apiKey string) (*http.Request, error) {
	httpRequest, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("Authorization", "Bearer "+apiKey)
	return httpRequest, nil
}

type errorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Param   interface{} `json:"param"` // interface{} because it can be of different types
	Type    string      `json:"type"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

func handleHTTPResponse(response *http.Response) error {
	if response.StatusCode < 400 {
		return nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var parsedResponse errorResponse
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	return errors.New(parsedResponse.Error.Message)
}
