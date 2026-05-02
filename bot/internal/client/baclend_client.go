package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"receptor/bot/internal/dto"
)

type BackendClient interface {
	Register(ctx context.Context, email, password, name string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	Me(ctx context.Context, token string) (*dto.User, error)

	GetRecipes(ctx context.Context) ([]dto.RecipeResponse, error)
	GetRecipeByID(ctx context.Context, id int64) (*dto.RecipeResponse, error)
	GetRecipeDetails(ctx context.Context, id int64) (*dto.RecipeDetailsResponse, error)

	GetFavorites(ctx context.Context, token string) ([]dto.RecipeResponse, error)
	AddFavorite(ctx context.Context, token string, recipeID int64) error
	RemoveFavorite(ctx context.Context, token string, recipeID int64) error

	Match(ctx context.Context, ingredients []string) ([]dto.MatchResponse, error)
	UploadImage(ctx context.Context, image []byte) ([]dto.MatchResponse, error)
}

type backendClient struct {
	baseURL    *url.URL
	httpClient *http.Client
}

type BackendError struct {
	StatusCode int
	Message    string
}

func (e *BackendError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

type apiErrorResponse struct {
	Error string `json:"error"`
}

func parseBackendError(status int, body []byte) error {
	var e apiErrorResponse
	if len(body) > 0 {
		if err := json.Unmarshal(body, &e); err == nil && strings.TrimSpace(e.Error) != "" {
			return &BackendError{StatusCode: status, Message: strings.TrimSpace(e.Error)}
		}
	}
	msg := strings.TrimSpace(string(body))
	if msg == "" {
		msg = http.StatusText(status)
	}
	return &BackendError{StatusCode: status, Message: msg}
}

func NewBackendClient(baseURL string, httpClient *http.Client) (BackendClient, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid base url: %q", baseURL)
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	if httpClient.Timeout == 0 {
		httpClient.Timeout = 10 * time.Second
	}

	return &backendClient{
		baseURL:    parsed,
		httpClient: httpClient,
	}, nil
}

func (c *backendClient) UploadImage(ctx context.Context, image []byte) ([]dto.MatchResponse, error) {
	if ctx == nil {
		return nil, fmt.Errorf("ctx is nil")
	}
	if len(image) == 0 {
		return nil, fmt.Errorf("image is empty")
	}

	endpoint := *c.baseURL
	endpoint.Path = path.Join(endpoint.Path, "/upload")

	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)

	part, err := writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		_ = writer.Close()
		return nil, fmt.Errorf("create multipart field: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(image)); err != nil {
		_ = writer.Close()
		return nil, fmt.Errorf("write image: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bodyBuf)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<10))
		if resp.StatusCode == http.StatusBadRequest {
			log.Printf("backend upload 400 response: %s", strings.TrimSpace(string(b)))
		}
		return nil, parseBackendError(resp.StatusCode, b)
	}

	var payload struct {
		Data []dto.MatchResponse `json:"data"`
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if payload.Data == nil {
		return []dto.MatchResponse{}, nil
	}
	return payload.Data, nil
}

func (c *backendClient) GetRecipeByID(ctx context.Context, id int64) (*dto.RecipeResponse, error) {
	if ctx == nil {
		return nil, fmt.Errorf("ctx is nil")
	}
	if id <= 0 {
		return nil, fmt.Errorf("invalid recipe id: %d", id)
	}

	endpoint := *c.baseURL
	endpoint.Path = path.Join(endpoint.Path, "/recipes", fmt.Sprintf("%d", id))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<10))
		return nil, parseBackendError(resp.StatusCode, b)
	}

	var payload struct {
		Data dto.RecipeResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &payload.Data, nil
}

func (c *backendClient) GetRecipeDetails(ctx context.Context, id int64) (*dto.RecipeDetailsResponse, error) {
	if ctx == nil {
		return nil, fmt.Errorf("ctx is nil")
	}
	if id <= 0 {
		return nil, fmt.Errorf("invalid recipe id: %d", id)
	}

	endpoint := *c.baseURL
	endpoint.Path = path.Join(endpoint.Path, "/recipes", fmt.Sprintf("%d", id), "details")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<10))
		return nil, parseBackendError(resp.StatusCode, b)
	}

	var payload struct {
		Data dto.RecipeDetailsResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &payload.Data, nil
}

func (c *backendClient) GetRecipes(ctx context.Context) ([]dto.RecipeResponse, error) {
	if ctx == nil {
		return nil, fmt.Errorf("ctx is nil")
	}

	endpoint := *c.baseURL
	endpoint.Path = path.Join(endpoint.Path, "/recipes")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<10))
		return nil, parseBackendError(resp.StatusCode, b)
	}

	var payload struct {
		Data []dto.RecipeResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return []dto.RecipeResponse{}, nil
	}
	return payload.Data, nil
}

func (c *backendClient) GetFavorites(ctx context.Context, token string) ([]dto.RecipeResponse, error) {
	if ctx == nil {
		return nil, fmt.Errorf("ctx is nil")
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("jwt is empty")
	}

	endpoint := *c.baseURL
	endpoint.Path = path.Join(endpoint.Path, "/favorites")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<10))
		return nil, parseBackendError(resp.StatusCode, b)
	}

	var payload struct {
		Data []dto.RecipeResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return []dto.RecipeResponse{}, nil
	}
	return payload.Data, nil
}

func (c *backendClient) AddFavorite(ctx context.Context, token string, recipeID int64) error {
	if ctx == nil {
		return fmt.Errorf("ctx is nil")
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("jwt is empty")
	}
	if recipeID <= 0 {
		return fmt.Errorf("invalid recipe id: %d", recipeID)
	}

	endpoint := *c.baseURL
	endpoint.Path = path.Join(endpoint.Path, "/favorites", fmt.Sprintf("%d", recipeID))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<10))
		return parseBackendError(resp.StatusCode, b)
	}

	return nil
}

func (c *backendClient) RemoveFavorite(ctx context.Context, token string, recipeID int64) error {
	if ctx == nil {
		return fmt.Errorf("ctx is nil")
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("jwt is empty")
	}
	if recipeID <= 0 {
		return fmt.Errorf("invalid recipe id: %d", recipeID)
	}

	endpoint := *c.baseURL
	endpoint.Path = path.Join(endpoint.Path, "/favorites", fmt.Sprintf("%d", recipeID))

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<10))
		return parseBackendError(resp.StatusCode, b)
	}

	return nil
}

func (c *backendClient) Login(ctx context.Context, email, password string) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("ctx is nil")
	}
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)
	if email == "" || password == "" {
		return "", fmt.Errorf("email and password are required")
	}

	in := dto.LoginInput{Email: email, Password: password}
	b, err := json.Marshal(in)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	endpoint := *c.baseURL
	endpoint.Path = path.Join(endpoint.Path, "/auth/login")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<10))
		return "", parseBackendError(resp.StatusCode, body)
	}

	var payload struct {
		Data dto.AuthResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if strings.TrimSpace(payload.Data.Token) == "" {
		return "", fmt.Errorf("empty token")
	}

	return payload.Data.Token, nil
}

func (c *backendClient) Register(ctx context.Context, email, password, name string) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("ctx is nil")
	}
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)
	name = strings.TrimSpace(name)
	if email == "" || password == "" || name == "" {
		return "", fmt.Errorf("email, password and name are required")
	}

	in := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}{
		Email:    email,
		Password: password,
		Name:     name,
	}

	b, err := json.Marshal(in)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	endpoint := *c.baseURL
	endpoint.Path = path.Join(endpoint.Path, "/auth/register")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<10))
		return "", parseBackendError(resp.StatusCode, body)
	}

	var payload struct {
		Data dto.AuthResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if strings.TrimSpace(payload.Data.Token) == "" {
		return "", fmt.Errorf("empty token")
	}
	return payload.Data.Token, nil
}

func (c *backendClient) Me(ctx context.Context, token string) (*dto.User, error) {
	if ctx == nil {
		return nil, fmt.Errorf("ctx is nil")
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	endpoint := *c.baseURL
	endpoint.Path = path.Join(endpoint.Path, "/auth/me")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<10))
		return nil, parseBackendError(resp.StatusCode, body)
	}

	var payload struct {
		Data dto.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &payload.Data, nil
}

func (c *backendClient) Match(ctx context.Context, ingredients []string) ([]dto.MatchResponse, error) {
	if ctx == nil {
		return nil, fmt.Errorf("ctx is nil")
	}
	reqBody := struct {
		Ingredients []string `json:"ingredients"`
	}{
		Ingredients: ingredients,
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	endpoint := *c.baseURL
	endpoint.Path = path.Join(endpoint.Path, "/match/by-ingredients")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<10))
		return nil, parseBackendError(resp.StatusCode, bodyBytes)
	}

	var payload struct {
		Data []dto.MatchResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Data == nil {
		return []dto.MatchResponse{}, nil
	}
	return payload.Data, nil
}

