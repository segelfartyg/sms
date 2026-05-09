package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"sms-backend/internal/database"
	"sms-backend/internal/handlers"
	"sms-backend/internal/repository"
	"sms-backend/internal/warehouse"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/kms?sslmode=disable"
	}

	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	warehouseURL := os.Getenv("WAREHOUSE_URL")
	if warehouseURL == "" {
		warehouseURL = "http://localhost:8081"
	}

	allowedOriginsRaw := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsRaw == "" {
		allowedOriginsRaw = "http://localhost:5173,http://localhost:5174"
	}
	allowedOrigins := strings.Split(allowedOriginsRaw, ",")

	wh := warehouse.NewClient(warehouseURL)

	pages := repository.NewPageRepo(db)
	boxes := repository.NewBoxRepo(db)

	pageHandler := handlers.NewPageHandler(pages, boxes)
	boxHandler := handlers.NewBoxHandler(boxes)
	warehouseHandler := handlers.NewWarehouseHandler(wh)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /pages", pageHandler.List)
	mux.HandleFunc("POST /pages", pageHandler.Create)
	mux.HandleFunc("GET /pages/{pageID}", pageHandler.Get)
	mux.HandleFunc("PUT /pages/{pageID}", pageHandler.Update)
	mux.HandleFunc("DELETE /pages/{pageID}", pageHandler.Delete)

	mux.HandleFunc("GET /pages/{pageID}/boxes", boxHandler.List)
	mux.HandleFunc("POST /pages/{pageID}/boxes", boxHandler.Create)
	mux.HandleFunc("GET /pages/{pageID}/boxes/{boxID}", boxHandler.Get)
	mux.HandleFunc("PUT /pages/{pageID}/boxes/{boxID}", boxHandler.Update)
	mux.HandleFunc("DELETE /pages/{pageID}/boxes/{boxID}", boxHandler.Delete)

	mux.HandleFunc("GET /slug/{slug}", pageHandler.GetBySlug)
	mux.HandleFunc("GET /datasources", warehouseHandler.ListDatasources)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
