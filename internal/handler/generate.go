package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bajra-manandhar17/logscope-v2/internal/generator"
)

// GenerateHandler handles POST /api/generate.
func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	var cfg generator.GenerateConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_config", "failed to parse request body")
		return
	}
	if err := cfg.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_config", err.Error())
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "internal_error", "streaming unsupported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	batches, errc := generator.Generate(r.Context(), cfg)

	total := 0
	for batch := range batches {
		total += len(batch)
		data, _ := json.Marshal(map[string][]string{"lines": batch})
		fmt.Fprintf(w, "event: batch\ndata: %s\n\n", data)
		flusher.Flush()
	}

	if err := <-errc; err != nil {
		// Client disconnected — stop silently.
		return
	}

	data, _ := json.Marshal(map[string]int{"totalLines": total})
	fmt.Fprintf(w, "event: done\ndata: %s\n\n", data)
	flusher.Flush()
}
