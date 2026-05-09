package warehouse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sms-backend/internal/models"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) ListDatasources(ctx context.Context) ([]models.Datasource, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/datasources", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()
	var out []models.Datasource
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return out, nil
}

