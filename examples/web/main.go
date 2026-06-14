package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//go:embed dist
var dist embed.FS

func mustBaseURL(env, fallback string) string {
	raw := os.Getenv(env)
	if raw == "" {
		raw = fallback
	}
	log.Printf("%s = %s", env, raw)
	return strings.TrimSuffix(raw, "/")
}

// proxyGet fetches upstreamURL and forwards its body and status code as JSON.
func proxyGet(w http.ResponseWriter, upstreamURL string) {
	resp, err := http.Get(upstreamURL)
	if err != nil {
		http.Error(w, "upstream unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	root, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Fatal(err)
	}

	backendURL := mustBaseURL("BACKEND_URL", "http://localhost:8080")
	warehouseURL := mustBaseURL("WAREHOUSE_URL", "http://localhost:8081")

	http.HandleFunc("GET /api/slug/{slug}", func(w http.ResponseWriter, r *http.Request) {
		proxyGet(w, backendURL+"/slug/"+url.PathEscape(r.PathValue("slug")))
	})
	http.HandleFunc("GET /warehouse/datasources/{id}", func(w http.ResponseWriter, r *http.Request) {
		proxyGet(w, warehouseURL+"/datasources/"+url.PathEscape(r.PathValue("id")))
	})

	fileServer := http.FileServer(http.FS(root))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/")
		if name == "" {
			name = "."
		}
		if _, err := root.Open(name); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}
		r2 := r.Clone(r.Context())
		r2.URL.Path = "/"
		fileServer.ServeHTTP(w, r2)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5174"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
