package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

// ── Kubernetes client ─────────────────────────────────────────────────────────

func buildK8sClient() (*kubernetes.Clientset, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			home, _ := os.UserHomeDir()
			kubeconfig = home + "/.kube/config"
		}
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("build config: %w", err)
		}
	}
	return kubernetes.NewForConfig(cfg)
}

// ── Core ──────────────────────────────────────────────────────────────────────

func findOrCreateDatasource(wh *warehouseClient, clusterName string) (datasource, error) {
	list, err := wh.listDatasources()
	if err != nil {
		return datasource{}, fmt.Errorf("list datasources: %w", err)
	}
	for _, ds := range list {
		if ds.Type == "k8s" && ds.Description == clusterName {
			log.Printf("found existing datasource %s", ds.ID)
			return ds, nil
		}
	}
	ds, err := wh.createDatasource("k8s", clusterName)
	if err != nil {
		return datasource{}, fmt.Errorf("create datasource: %w", err)
	}
	log.Printf("created datasource %s", ds.ID)
	return ds, nil
}

func sync(ctx context.Context, k8s *kubernetes.Clientset, wh *warehouseClient, ds datasource, namespace string) error {
	pods, err := k8s.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("list pods: %w", err)
	}

	existing, err := wh.listDatapoints(ds.ID)
	if err != nil {
		return fmt.Errorf("list datapoints: %w", err)
	}
	for _, dp := range existing {
		if err := wh.deleteDatapoint(ds.ID, dp.ID); err != nil {
			log.Printf("warn: delete datapoint %s: %v", dp.ID, err)
		}
	}

	pos := 0

	if err := wh.createDatapoint(ds.ID, "h1", ds.Description, "", pos); err != nil {
		return fmt.Errorf("create h1: %w", err)
	}
	pos++

	synced := time.Now().UTC().Format("2006-01-02 15:04:05 UTC")
	if err := wh.createDatapoint(ds.ID, "h2", "Synced "+synced, "", pos); err != nil {
		return fmt.Errorf("create h2: %w", err)
	}
	pos++

	for _, pod := range pods.Items {
		phase := string(pod.Status.Phase)
		if phase == "" {
			phase = "Unknown"
		}
		content := fmt.Sprintf("%s/%s (%s)", pod.Namespace, pod.Name, phase)
		desc := pod.Spec.NodeName
		if err := wh.createDatapoint(ds.ID, "li", content, desc, pos); err != nil {
			log.Printf("warn: pod %s: %v", pod.Name, err)
		}
		pos++
	}

	log.Printf("synced %d pods", len(pods.Items))
	return nil
}

func main() {
	warehouseURL := os.Getenv("WAREHOUSE_URL")
	if warehouseURL == "" {
		warehouseURL = "http://localhost:8081"
	}
	clusterName := os.Getenv("CLUSTER_NAME")
	if clusterName == "" {
		clusterName = "k8s-cluster"
	}
	namespace := os.Getenv("NAMESPACE") // empty = all namespaces

	interval := 30 * time.Second
	if v := os.Getenv("SYNC_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			interval = d
		}
	}

	wh := newWarehouseClient(warehouseURL)

	k8s, err := buildK8sClient()
	if err != nil {
		log.Fatalf("k8s client: %v", err)
	}

	ds, err := findOrCreateDatasource(wh, clusterName)
	if err != nil {
		log.Fatalf("datasource: %v", err)
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		if err := sync(ctx, k8s, wh, ds, namespace); err != nil {
			log.Printf("sync error: %v", err)
		}
		cancel()
		time.Sleep(interval)
	}
}
