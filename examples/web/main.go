package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var dist embed.FS

func mustProxy(env, fallback string) *httputil.ReverseProxy {
	raw := os.Getenv(env)
	if raw == "" {
		raw = fallback
	}
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("invalid %s: %v", env, err)
	}
	log.Printf("%s → %s", env, u)
	return httputil.NewSingleHostReverseProxy(u)
}

func main() {
	root, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/api/", http.StripPrefix("/api", mustProxy("BACKEND_URL", "http://localhost:8080")))
	http.Handle("/warehouse/", http.StripPrefix("/warehouse", mustProxy("WAREHOUSE_URL", "http://localhost:8081")))

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
