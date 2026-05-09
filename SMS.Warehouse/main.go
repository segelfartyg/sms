package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"sms-warehouse/internal/database"
	"sms-warehouse/internal/handlers"
	"sms-warehouse/internal/repository"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5433/warehouse?sslmode=disable"
	}

	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	allowedOriginsRaw := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsRaw == "" {
		allowedOriginsRaw = "http://localhost:5173,http://localhost:5174"
	}
	allowedOrigins := strings.Split(allowedOriginsRaw, ",")

	dsRepo := repository.NewDatasourceRepo(db)
	dpRepo := repository.NewDatapointRepo(db)

	ds := handlers.NewDatasourceHandler(dsRepo, dpRepo)
	dp := handlers.NewDatapointHandler(dpRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /datasources", ds.List)
	mux.HandleFunc("POST /datasources", ds.Create)
	mux.HandleFunc("GET /datasources/{id}", ds.Get)
	mux.HandleFunc("PUT /datasources/{id}", ds.Update)
	mux.HandleFunc("DELETE /datasources/{id}", ds.Delete)

	mux.HandleFunc("GET /datasources/{id}/datapoints", dp.List)
	mux.HandleFunc("POST /datasources/{id}/datapoints", dp.Create)
	mux.HandleFunc("GET /datasources/{id}/datapoints/{dpID}", dp.Get)
	mux.HandleFunc("PUT /datasources/{id}/datapoints/{dpID}", dp.Update)
	mux.HandleFunc("DELETE /datasources/{id}/datapoints/{dpID}", dp.Delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("listening on :%s (origins: %s)", port, allowedOriginsRaw)
	if err := http.ListenAndServe(":"+port, cors(allowedOrigins, mux)); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func cors(origins []string, next http.Handler) http.Handler {
	allowed := make(map[string]bool, len(origins))
	for _, o := range origins {
		allowed[strings.TrimSpace(o)] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
