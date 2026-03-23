package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bajra-manandhar17/logscope-v2/internal/analyzer"
)

// MaxUploadBytes is the per-request body limit. Overridable in tests.
var MaxUploadBytes int64 = 100 << 20 // 100MB

var validFormats = map[string]bool{
	"":          true,
	"auto":      true,
	"json":      true,
	"plaintext": true,
}

// AnalyzeHandler handles POST /api/analyze.
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadBytes)

	format := r.URL.Query().Get("format")
	if format == "auto" {
		format = ""
	}
	if !validFormats[format] {
		writeError(w, http.StatusBadRequest, "invalid_format", "unsupported format: "+format)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(w, http.StatusRequestEntityTooLarge, "file_too_large", "file exceeds 100MB limit")
			return
		}
		writeError(w, http.StatusBadRequest, "invalid_request", "failed to parse multipart form")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "missing file field")
		return
	}
	defer file.Close()

	result, err := analyzer.Analyze(context.Background(), file, format)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "analysis failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg, "code": code})
}
