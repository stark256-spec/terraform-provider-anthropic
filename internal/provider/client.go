package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.anthropic.com"
const apiVersion = "2023-06-01"
const betaHeader = "workspaces-2025-04-14"

// AnthropicClient is a minimal HTTP client for the Anthropic admin API.
type AnthropicClient struct {
	apiKey  string
	orgID   string
	baseURL string
	http    *http.Client
}

func newClient(apiKey, orgID, baseURL string) *AnthropicClient {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &AnthropicClient{
		apiKey:  apiKey,
		orgID:   orgID,
		baseURL: baseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *AnthropicClient) do(ctx context.Context, method, path string, body any) ([]byte, int, error) {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request: %w", err)
		}
		buf = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, buf)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", apiVersion)
	req.Header.Set("anthropic-beta", betaHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}
	return data, resp.StatusCode, nil
}

// ── Workspace ──────────────────────────────────────────────────────────────

type Workspace struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

func (c *AnthropicClient) CreateWorkspace(ctx context.Context, name, displayName string) (*Workspace, error) {
	body := map[string]string{"name": name}
	if displayName != "" {
		body["display_name"] = displayName
	}
	data, _, err := c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/organizations/%s/workspaces", c.orgID), body)
	if err != nil {
		return nil, err
	}
	var w Workspace
	return &w, json.Unmarshal(data, &w)
}

func (c *AnthropicClient) GetWorkspace(ctx context.Context, id string) (*Workspace, error) {
	data, _, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/organizations/%s/workspaces/%s", c.orgID, id), nil)
	if err != nil {
		return nil, err
	}
	var w Workspace
	return &w, json.Unmarshal(data, &w)
}

func (c *AnthropicClient) UpdateWorkspace(ctx context.Context, id, displayName string) (*Workspace, error) {
	body := map[string]string{"display_name": displayName}
	data, _, err := c.do(ctx, http.MethodPatch, fmt.Sprintf("/v1/organizations/%s/workspaces/%s", c.orgID, id), body)
	if err != nil {
		return nil, err
	}
	var w Workspace
	return &w, json.Unmarshal(data, &w)
}

func (c *AnthropicClient) ArchiveWorkspace(ctx context.Context, id string) error {
	_, _, err := c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/organizations/%s/workspaces/%s/archive", c.orgID, id), nil)
	return err
}

// ── API Key ────────────────────────────────────────────────────────────────

type APIKey struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	WorkspaceID string  `json:"workspace_id"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	PartialKey  string  `json:"partial_key,omitempty"`
	Key         *string `json:"key,omitempty"` // only on creation
}

type CreateAPIKeyRequest struct {
	Name        string `json:"name"`
	WorkspaceID string `json:"workspace_id,omitempty"`
}

func (c *AnthropicClient) CreateAPIKey(ctx context.Context, name, workspaceID string) (*APIKey, error) {
	body := CreateAPIKeyRequest{Name: name, WorkspaceID: workspaceID}
	data, _, err := c.do(ctx, http.MethodPost, fmt.Sprintf("/v1/organizations/%s/api_keys", c.orgID), body)
	if err != nil {
		return nil, err
	}
	var k APIKey
	return &k, json.Unmarshal(data, &k)
}

func (c *AnthropicClient) GetAPIKey(ctx context.Context, id string) (*APIKey, error) {
	data, _, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/organizations/%s/api_keys/%s", c.orgID, id), nil)
	if err != nil {
		return nil, err
	}
	var k APIKey
	return &k, json.Unmarshal(data, &k)
}

func (c *AnthropicClient) UpdateAPIKey(ctx context.Context, id, name string) (*APIKey, error) {
	body := map[string]string{"name": name}
	data, _, err := c.do(ctx, http.MethodPatch, fmt.Sprintf("/v1/organizations/%s/api_keys/%s", c.orgID, id), body)
	if err != nil {
		return nil, err
	}
	var k APIKey
	return &k, json.Unmarshal(data, &k)
}

func (c *AnthropicClient) DeleteAPIKey(ctx context.Context, id string) error {
	_, _, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/v1/organizations/%s/api_keys/%s", c.orgID, id), nil)
	return err
}
