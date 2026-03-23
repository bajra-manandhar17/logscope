# 01 — Project Scaffolding

**Complexity:** Simple
**Phase:** 1 — Scaffolding
**Blocked by:** None
**Blocks:** All other tasks

## Objective

Initialize Go module, frontend scaffolding (Vite + React + TS + Tailwind + shadcn + React Router + @tanstack/react-virtual), and Makefile with build targets.

## Scope

- `go mod init github.com/bajra-manandhar17/logscope-v2`
- Vite + React + TypeScript project in `frontend/`
- Install and configure: Tailwind CSS, shadcn/ui, React Router, @tanstack/react-virtual, Zustand, Recharts
- Makefile with `build`, `dev-backend`, `dev-frontend` targets
- `vite.config.ts` with `/api` proxy to `:8080`
- `tsconfig.json` with strict mode

## Acceptance Criteria

- [x] `go mod tidy` succeeds
- [x] `cd frontend && npm install && npm run build` succeeds
- [x] `make build` target defined (may fail until Go entry point exists)
- [x] Vite dev server starts on `:5173`
- [x] Tailwind + shadcn configured and rendering
