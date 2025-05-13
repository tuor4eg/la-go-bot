package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	PATCH  Method = "PATCH"
)

const (
	UserInfoEndpoint       = "/telegram/user"
	ClosestCamerasEndpoint = "/telegram/camera/closest"
)

type ApiClient struct {
	secretKey string
	baseURL   string
	client    *http.Client
}

func NewApiClient(secretKey string, baseURL string) *ApiClient {
	return &ApiClient{
		secretKey: secretKey,
		baseURL:   baseURL,
		client:    &http.Client{},
	}
}

func (a *ApiClient) MakeRequest(method Method, endpoint string, data any) (string, error) {
	var body *bytes.Buffer
	if data != nil && (method == POST || method == PUT || method == PATCH) {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("error marshaling data: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	fullURL := a.baseURL + endpoint
	log.Printf("Making %s request to: %s", method, fullURL)

	req, err := http.NewRequest(string(method), fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.secretKey)
	if body != nil && body.Len() > 0 {
		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(body)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	// Читаем ответ
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return string(respBody), nil
}

func (a *ApiClient) Get(endpoint string) (string, error) {
	return a.MakeRequest(GET, endpoint, nil)
}

func (a *ApiClient) Post(endpoint string, data any) (string, error) {
	return a.MakeRequest(POST, endpoint, data)
}

func (a *ApiClient) GetUserInfo(userId string) (string, error) {
	return a.Get(fmt.Sprintf("%s/%s", UserInfoEndpoint, userId))
}

func (a *ApiClient) getClosestCameras(latitude float64, longitude float64, distance float64) (string, error) {
	return a.Post(ClosestCamerasEndpoint, map[string]any{
		"latitude":  latitude,
		"longitude": longitude,
		"distance":  distance,
	})

}
