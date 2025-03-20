package uqam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL = "https://monportail.uqam.ca"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

type LoginResponse struct {
	Token string `json:"token"`
}

type GradeResponse struct {
	Data struct {
		Resultats []struct {
			Programmes []struct {
				Activites []struct {
					Total float32 `json:"total"`
				} `json:"activites"`
			} `json:"programmes"`
		} `json:"resultats"`
	} `json:"data"`
}

func (c *Client) Login(username, password string) (string, error) {
	payload, err := json.Marshal(map[string]string{
		"identifiant": username,
		"motDePasse":  password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal login payload: %w", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/authentification", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read login response: %w", err)
	}

	var response LoginResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal login response: %w", err)
	}

	if response.Token == "" {
		return "", fmt.Errorf("empty token received")
	}

	return response.Token, nil
}

func (c *Client) GetGrade(token, semester, course string) (string, error) {
	url := fmt.Sprintf("%s/apis/resultatActivite/identifiant/%s/%s", baseURL, semester, course)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create grade request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send grade request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read grade response: %w", err)
	}

	var response GradeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal grade response: %w", err)
	}

	if len(response.Data.Resultats) == 0 || len(response.Data.Resultats[0].Programmes) == 0 || len(response.Data.Resultats[0].Programmes[0].Activites) == 0 {
		return "", fmt.Errorf("no grade found")
	}

	return fmt.Sprintf("%.2f", response.Data.Resultats[0].Programmes[0].Activites[0].Total), nil
}
