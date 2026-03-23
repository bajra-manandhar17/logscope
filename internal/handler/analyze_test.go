package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bajra-manandhar17/logscope-v2/internal/handler"
)

func multipartBody(t *testing.T, content string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", "test.log")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.WriteString(fw, content)
	if err != nil {
		t.Fatal(err)
	}
	w.Close()
	return &buf, w.FormDataContentType()
}

func TestAnalyzeHandler_ReturnsAnalysisJSON(t *testing.T) {
	logContent := `{"level":"info","msg":"started","ts":"2024-01-01T00:00:00Z"}` + "\n"
	body, ct := multipartBody(t, logContent)

	req := httptest.NewRequest(http.MethodPost, "/api/analyze", body)
	req.Header.Set("Content-Type", ct)
	rec := httptest.NewRecorder()

	handler.AnalyzeHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("expected application/json, got %q", ct)
	}
	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, ok := result["format_detected"]; !ok {
		t.Error("response missing format_detected field")
	}
}

func TestAnalyzeHandler_FormatQueryParam(t *testing.T) {
	logContent := "2024-01-01 00:00:00 INFO started\n"
	body, ct := multipartBody(t, logContent)

	req := httptest.NewRequest(http.MethodPost, "/api/analyze?format=plaintext", body)
	req.Header.Set("Content-Type", ct)
	rec := httptest.NewRecorder()

	handler.AnalyzeHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAnalyzeHandler_InvalidFormat(t *testing.T) {
	body, ct := multipartBody(t, "anything\n")

	req := httptest.NewRequest(http.MethodPost, "/api/analyze?format=csv", body)
	req.Header.Set("Content-Type", ct)
	rec := httptest.NewRecorder()

	handler.AnalyzeHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["code"] != "invalid_format" {
		t.Errorf("expected code=invalid_format, got %q", resp["code"])
	}
}

func TestAnalyzeHandler_FileTooLarge(t *testing.T) {
	// Lower the limit so we don't need to write 100MB in tests.
	orig := handler.MaxUploadBytes
	handler.MaxUploadBytes = 512
	t.Cleanup(func() { handler.MaxUploadBytes = orig })

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	ct := mw.FormDataContentType()

	go func() {
		defer pw.Close()
		fw, _ := mw.CreateFormFile("file", "big.log")
		chunk := bytes.Repeat([]byte("x"), 128)
		for i := 0; i < 10; i++ { // 1280 bytes > 512 limit
			if _, err := fw.Write(chunk); err != nil {
				return
			}
		}
		mw.Close()
	}()

	req := httptest.NewRequest(http.MethodPost, "/api/analyze", pr)
	req.Header.Set("Content-Type", ct)
	rec := httptest.NewRecorder()

	handler.AnalyzeHandler(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", rec.Code)
	}
	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["code"] != "file_too_large" {
		t.Errorf("expected code=file_too_large, got %q", resp["code"])
	}
}

func TestAnalyzeHandler_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/analyze", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.AnalyzeHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
