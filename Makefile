.PHONY: build dev-backend dev-frontend

build:
	cd frontend && npm run build
	go build -o bin/logscope ./cmd/logscope

dev-backend:
	go run ./cmd/logscope

dev-frontend:
	cd frontend && npm run dev
