package models

import "time"

type Datasource struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Content     string      `json:"content"`
	Datapoints  []Datapoint `json:"datapoints,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Datapoint struct {
	ID           string    `json:"id"`
	DatasourceID string    `json:"datasource_id"`
	Tag          string    `json:"tag"`
	Content      string    `json:"content"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
