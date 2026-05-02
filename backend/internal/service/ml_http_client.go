package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type httpMLClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPMLClient(baseURL string) MLClient {
	return &httpMLClient{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}

func parseFastAPIDetail(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil && strings.TrimSpace(s) != "" {
		return strings.TrimSpace(s)
	}
	var items []struct {
		Msg string `json:"msg"`
	}
	if err := json.Unmarshal(raw, &items); err == nil && len(items) > 0 {
		parts := make([]string, 0, len(items))
		for _, it := range items {
			if t := strings.TrimSpace(it.Msg); t != "" {
				parts = append(parts, t)
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, "; ")
		}
	}
	return ""
}

func (c *httpMLClient) DetectIngredients(ctx context.Context, b []byte) ([]string, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("ML_URL is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "photo.jpg")
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(b); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/detect", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := ""
		var payload struct {
			Detail json.RawMessage `json:"detail"`
		}
		if json.Unmarshal(bodyBytes, &payload) == nil {
			msg = parseFastAPIDetail(payload.Detail)
		}
		if msg == "" && len(bodyBytes) > 0 && len(bodyBytes) < 2048 {
			msg = strings.TrimSpace(string(bodyBytes))
		}
		return nil, &MLServiceError{UpstreamStatus: resp.StatusCode, Message: msg}
	}

	var out struct {
		Ingredients []string `json:"ingredients"`
	}
	if err := json.Unmarshal(bodyBytes, &out); err != nil {
		return nil, err
	}
	return out.Ingredients, nil
}
