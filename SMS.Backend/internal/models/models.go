package models

import "time"

type Page struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Boxes     []Box     `json:"boxes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Box struct {
	ID           string    `json:"id"`
	PageID       string    `json:"page_id"`
	Position     int       `json:"position"`
	DatasourceID *string   `json:"datasource_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Datasource mirrors the shape returned by SMS.Warehouse.
type Datasource struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Content     string      `json:"content"`
	Datapoints  []Datapoint `json:"datapoints,omitempty"`
}

type Datapoint struct {
	ID           string `json:"id"`
	DatasourceID string `json:"datasource_id"`
	Tag          string `json:"tag"`
	Content      string `json:"content"`
	Description  string `json:"description"`
}
