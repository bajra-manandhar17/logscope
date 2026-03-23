package server

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"time"

	"github.com/bajra-manandhar17/logscope-v2/internal/handler"
)

// New builds and returns an http.Server. Pass devMode=true to enable CORS for localhost:5173.
func New(staticFS fs.FS, devMode bool) *http.Server {
	mux := http.NewServeMux()
	registerRoutes(mux, staticFS)

	var handler http.Handler = mux
	if devMode {
		handler = corsMiddleware(mux)
	}

	return &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func registerRoutes(mux *http.ServeMux, staticFS fs.FS) {
	mux.HandleFunc("GET /api/health", healthHandler)
	mux.HandleFunc("POST /api/analyze", handler.AnalyzeHandler)
	mux.HandleFunc("POST /api/generate", handler.GenerateHandler)

	// Serve embedded SPA; strip the "dist" prefix so index.html is at "/"
	distFS, err := fs.Sub(staticFS, "dist")
	if err != nil {
		panic("failed to sub dist from embedded FS: " + err.Error())
	}
	mux.Handle("GET /", http.FileServer(http.FS(distFS)))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
