package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// ── Warehouse types ───────────────────────────────────────────────────────────

type datasource struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type datapoint struct {
	ID      string `json:"id"`
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

// ── Warehouse client ──────────────────────────────────────────────────────────

type warehouseClient struct {
	base string
	http *http.Client
}

func newWarehouseClient(base string) *warehouseClient {
	return &warehouseClient{base: base, http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *warehouseClient) request(method, path string, body any) (*http.Response, error) {
	var buf *bytes.Buffer
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(data)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req, err := http.NewRequest(method, c.base+path, buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.http.Do(req)
}

func (c *warehouseClient) listDatasources() ([]datasource, error) {
	resp, err := c.request("GET", "/datasources", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out []datasource
	return out, json.NewDecoder(resp.Body).Decode(&out)
}

func (c *warehouseClient) createDatasource(dsType, description string) (datasource, error) {
	body := map[string]string{"type": dsType, "description": description, "content": ""}
	resp, err := c.request("POST", "/datasources", body)
	if err != nil {
		return datasource{}, err
	}
	defer resp.Body.Close()
	var out datasource
	return out, json.NewDecoder(resp.Body).Decode(&out)
}

func (c *warehouseClient) listDatapoints(dsID string) ([]datapoint, error) {
	resp, err := c.request("GET", "/datasources/"+dsID+"/datapoints", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out []datapoint
	return out, json.NewDecoder(resp.Body).Decode(&out)
}

func (c *warehouseClient) createDatapoint(dsID, tag, content, description string, position int) error {
	body := map[string]any{
		"tag":         tag,
		"content":     content,
		"description": description,
		"position":    position,
	}
	resp, err := c.request("POST", "/datasources/"+dsID+"/datapoints", body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("create datapoint: status %d", resp.StatusCode)
	}
	return nil
}

func (c *warehouseClient) deleteDatapoint(dsID, dpID string) error {
	resp, err := c.request("DELETE", "/datasources/"+dsID+"/datapoints/"+dpID, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *warehouseClient) findOrCreateDatasource(name string) (datasource, error) {
	list, err := c.listDatasources()
	if err != nil {
		return datasource{}, fmt.Errorf("list datasources: %w", err)
	}
	for _, ds := range list {
		if ds.Type == "fastfetch" && ds.Description == name {
			return ds, nil
		}
	}
	return c.createDatasource("fastfetch", name)
}

// ── Fastfetch ─────────────────────────────────────────────────────────────────

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[mGKHF]`)

func runFastfetch() ([]string, error) {
	out, err := exec.Command("fastfetch", "--logo", "none", "--color", "black").Output()
	if err != nil {
		return nil, err
	}

	var props []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(ansiEscape.ReplaceAllString(line, ""))
		if strings.Contains(line, ": ") {
			props = append(props, line)
		}
	}
	return props, nil
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	warehouseURL := flag.String("warehouse", "http://localhost:8081", "warehouse URL")
	name := flag.String("name", "", "datasource name (required)")
	flag.Parse()

	if *name == "" {
		log.Fatal("-name is required")
	}

	props, err := runFastfetch()
	if err != nil {
		log.Fatalf("fastfetch: %v", err)
	}

	wh := newWarehouseClient(*warehouseURL)

	ds, err := wh.findOrCreateDatasource(*name)
	if err != nil {
		log.Fatalf("datasource: %v", err)
	}

	existing, err := wh.listDatapoints(ds.ID)
	if err != nil {
		log.Fatalf("list datapoints: %v", err)
	}
	for _, dp := range existing {
		if err := wh.deleteDatapoint(ds.ID, dp.ID); err != nil {
			log.Printf("warn: delete datapoint %s: %v", dp.ID, err)
		}
	}

	pos := 0
	if err := wh.createDatapoint(ds.ID, "h1", *name, "", pos); err != nil {
		log.Fatalf("create h1: %v", err)
	}
	pos++

	synced := time.Now().UTC().Format("2006-01-02 15:04:05 UTC")
	if err := wh.createDatapoint(ds.ID, "h2", "Fetched "+synced, "", pos); err != nil {
		log.Fatalf("create h2: %v", err)
	}
	pos++

	for _, prop := range props {
		if err := wh.createDatapoint(ds.ID, "li", prop, "", pos); err != nil {
			log.Printf("warn: %s: %v", prop, err)
		}
		pos++
	}

	log.Printf("synced %d properties to %q", len(props), *name)
}
