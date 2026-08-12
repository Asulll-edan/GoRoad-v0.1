# GoRoad — AGENTS.md

Motorcycle-touring community platform. Monorepo with 5 services: Go backend (primary), 3 Python microservices, Flutter mobile app, Next.js admin panel, plus infra (Postgres/PostGIS+TimescaleDB, Redis, NATS, MinIO, LiveKit, OSRM, nginx, Prometheus/Grafana/Loki). All docs (README, MAIN-PLAN) are written in Indonesian.

## Service layout

- `go-road-backend/` — primary backend, Go 1.23, module `go-road-backend`. Three binaries: `cmd/api` (Fiber v3 REST, port 8080), `cmd/ws` (gorilla/websocket + NATS, port 8081), `cmd/worker` (asynq tasks). Domain-first pattern: `internal/domain/{x}/` holds entity + repo + service interfaces; implementations live in `internal/repository/`, `internal/service/`, `internal/handler/`. Uses GORM, go-redis, NATS JetStream, MinIO, Viper, Zap, uber-fx DI.
- `py-ai-service/` — FastAPI + Gemini AI (REST 8001, gRPC 50051)
- `py-analytics-service/` — FastAPI analytics (REST 8002, gRPC 50052)
- `py-data-pipeline/` — Celery workers + beat, Redis broker
- `py-shared/` — shared Python utils (redis_client, nats_client, logging); not a package, import by path
- `go_road_app_flutter/` — Flutter (Riverpod, GoRouter, Dio, Drift)
- `go-road-admin/` — Next.js 14 (NOT 15 as README claims), TanStack Query/Table, Tailwind
- `proto/` — gRPC .proto sources; generated Python stubs are committed under `py-*/app/proto/`

## Commands (exact, CI-verified)

```bash
# Go — CI runs all of these from go-road-backend/
go mod download
go vet ./...
go build ./...
go test ./... -cover -race -timeout 120s

# Python (each service dir; CI runs this)
pip install -r requirements.txt
pip install pytest pytest-asyncio pytest-cov
pytest tests/ -v --cov=app --cov-report=term --ignore=tests/integration   # integration/ dir is skipped in CI

# Flutter
flutter pub get
flutter analyze
flutter test --coverage

# Admin
npm ci
npm run lint    # next lint
npm run build

# k6 load test (targets localhost:8080, override with BASE_URL)
k6 run loadtest/k6_scenario.js
```

Run services manually: `go run ./cmd/{api,ws,worker}/main.go`; `uvicorn app.main:app --reload --port 8001`; `celery -A app.tasks.celery_app worker|beat`. Prereqs (postgres, redis, nats, minio) must be up first.

## Setup & env

- Copy `.env.example` → `.env`; required keys: `JWT_SECRET`, `GEMINI_API_KEY`, `OPENWEATHER_API_KEY`, `LIVEKIT_API_KEY`, `LIVEKIT_API_SECRET`. `.env` is gitignored.
- Dev DB creds are hardcoded everywhere (postgres / sandiandasalah / goroad_dbvi1) — see `.env-information.md`, `.env.example`, docker-compose.yml.
- Ports: API 8080, WS 8081, admin 3000, postgres 5432, redis 6379, nats 4222, minio 9000/9001, AI gRPC 50051, analytics gRPC 50052.

## Verified gotchas

- **`docker compose up --build` fails**: the `admin` service in `docker-compose.yml` uses build context `./admin`, but the real directory is `go-road-admin` (docker-compose.dev.yml correctly uses `./go-road-admin`). CI's `docker compose build --parallel` hits this too.
- **Migrations**: applied automatically when the postgres container starts via `./migrations` mounted at `/docker-entrypoint-initdb.d` — only run on a fresh (empty) volume; there is no migration runner in code despite golang-migrate being in go.mod. Add new migrations as `0000NN_{name}.{up,down}.sql` pairs; up-files are executed in filename order at init.
- **gRPC is one-way today**: Python services have committed generated stubs (`py-*/app/proto/*_pb2*.py`, `*_pb2_grpc.py`) and run gRPC servers; the Go backend has grpc deps in go.mod but **no generated `.pb.go` code** — no Go gRPC client exists yet. `MAIN-PLAN.md` describes the aspirational v3.0 architecture, not the current state; trust the code and CI over it.
- **Flutter codegen is unused so far**: `riverpod_generator`/`json_serializable`/`build_runner` are dev deps but no `@Riverpod` annotations or committed `*.g.dart` files exist. If you add codegen, run `dart run build_runner build` and commit the outputs.
- **Changes to `docker/nginx.conf`** are only picked up via docker; manual runs of Go/Python services do not go through nginx.
