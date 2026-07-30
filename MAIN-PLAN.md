# Go Road v3.0 — Implementation Plan

**Versi:** 1.0  
**Tanggal:** 2026-07-30  
**Scope:** Full-stack implementation dengan Python microservices, Redis caching, lazy loading, dan pagination

---

## Ringkasan

Go Road adalah platform multi-platform (mobile + web) untuk komunitas touring motor. Plan ini menambahkan **Python microservices** untuk AI/ML, analytics, dan data pipeline — melengkapi Go backend utama. Setiap layer mengimplementasikan **Redis caching**, **lazy loading**, dan **cursor-based pagination** secara konsisten.

---

## User Review Required

> [!IMPORTANT]
> **Python Microservices Integration** — Plan ini menambahkan 3 Python microservice di samping Go backend utama:
> 1. `py-ai-service` — AI Assistant (Gemini), itinerary generator, cost estimator
> 2. `py-analytics-service` — Rider statistics, reporting, data aggregation, badge evaluation
> 3. `py-data-pipeline` — ETL jobs, weather data processing, POI cache refresh, smart detection ML models
>
> Komunikasi Go ↔ Python via **gRPC** (low-latency) dan **NATS** (event-driven).

> [!WARNING]
> **Breaking Change dari System Brief v3.0:**
> - System brief asli hanya menyebutkan Go untuk backend. Plan ini menambahkan Python sebagai complementary services, **bukan** menggantikan Go. Go tetap sebagai primary backend (REST API, WebSocket, core business logic). Python menangani compute-heavy tasks (AI, ML, analytics, data pipeline).

> [!IMPORTANT]
> **Caching Strategy** — Redis digunakan sebagai multi-layer cache:
> - **L1 (In-Process)**: Go in-memory cache (bigcache) untuk hot data (TTL 30s)
> - **L2 (Redis)**: Distributed cache untuk semua services (TTL bervariasi per data type)
> - **L3 (CDN)**: Static assets dan map tiles
> - Python services juga menggunakan Redis sebagai shared cache layer

---

## Open Questions

> [!IMPORTANT]
> 1. **Python Framework Preference** — Apakah ada preferensi antara **FastAPI** (async, modern) vs **Flask** (simpler) untuk Python microservices? Plan ini merekomendasikan **FastAPI** karena async support yang cocok untuk AI streaming dan gRPC.

> [!IMPORTANT]
> 2. **gRPC vs REST untuk inter-service communication** — Plan ini merekomendasikan **gRPC** untuk Go ↔ Python communication (type-safe, streaming support, lebih cepat). Alternatif: REST via internal network. Mana yang diprefer?

> [!IMPORTANT]
> 3. **Celery vs Dramatiq untuk Python task queue** — Plan ini merekomendasikan **Celery** (Redis broker) untuk Python background tasks. Alternatif: **Dramatiq** (lebih simple, lebih modern). Mana yang diprefer?

> [!IMPORTANT]
> 4. **Phase Prioritas** — Apakah semua 7 phase harus dieksekusi berurutan, atau ada phase yang bisa di-skip/parallel?

---

## Arsitektur Baru dengan Python

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                            CLIENT LAYER                                       │
│  ┌─────────────────┐  ┌──────────────────┐  ┌────────────────────────────┐  │
│  │ Flutter Mobile   │  │ Flutter Web(PWA) │  │ Next.js Admin Dashboard    │  │
│  │ + Lazy Loading   │  │ + Lazy Loading   │  │ + Lazy Loading + Pagination│  │
│  │ + Pagination     │  │ + Pagination     │  │                            │  │
│  └────────┬────────┘  └────────┬─────────┘  └─────────────┬──────────────┘  │
└───────────┼─────────────────────┼──────────────────────────┼────────────────┘
            │                     │                          │
            ▼                     ▼                          ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                           GATEWAY (Nginx)                                     │
│  SSL Termination • Rate Limiting • WebSocket Upgrade • Load Balancing        │
└──────────────────────────────┬───────────────────────────────────────────────┘
                               │
         ┌─────────────────────┼──────────────────────┐
         ▼                     ▼                      ▼
┌──────────────────┐ ┌──────────────────┐ ┌──────────────────────────────────┐
│  GO BACKEND      │ │  GO BACKEND      │ │  PYTHON MICROSERVICES            │
│  (Primary)       │ │  (Realtime)      │ │                                  │
│                  │ │                  │ │  ┌────────────────────────────┐  │
│  ├─ REST API     │ │  ├─ WebSocket    │ │  │ py-ai-service (FastAPI)    │  │
│  │  (Fiber v3)   │ │  │  Gateway      │ │  │ ├─ Gemini Chat (SSE)      │  │
│  │               │ │  │  (gorilla/ws) │ │  │ ├─ Itinerary Generator    │  │
│  ├─ Redis Cache  │ │  ├─ NATS Pub/Sub │ │  │ ├─ Cost Estimator         │  │
│  │  Integration  │ │  │               │ │  │ ├─ Route Advisor          │  │
│  │               │ │  ├─ Presence     │ │  │ └─ Redis Cache Layer      │  │
│  ├─ Pagination   │ │  │  Tracking     │ │  └────────────────────────────┘  │
│  │  (Cursor)     │ │  │               │ │  ┌────────────────────────────┐  │
│  │               │ │  └─ Heartbeat    │ │  │ py-analytics-service       │  │
│  └─ gRPC Client  │ │                  │ │  │  (FastAPI)                 │  │
│    (to Python)   │ │                  │ │  │ ├─ Rider Statistics        │  │
│                  │ │                  │ │  │ ├─ Leaderboard Compute     │  │
│                  │ │                  │ │  │ ├─ Badge Evaluation        │  │
│                  │ │                  │ │  │ ├─ Report Generation       │  │
│                  │ │                  │ │  │ │  (PDF/CSV via ReportLab) │  │
│                  │ │                  │ │  │ └─ Redis Cache Layer       │  │
└──────────────────┘ └──────────────────┘ │  └────────────────────────────┘  │
                                          │  ┌────────────────────────────┐  │
                                          │  │ py-data-pipeline           │  │
                                          │  │  (Celery Workers)          │  │
                                          │  │ ├─ Weather Sync (cron)     │  │
                                          │  │ ├─ POI Cache Refresh       │  │
                                          │  │ ├─ Smart Detection ML      │  │
                                          │  │ ├─ Data Retention Cleanup  │  │
                                          │  │ ├─ Touring Feed Aggregator │  │
                                          │  │ └─ Service Reminder Check  │  │
                                          │  └────────────────────────────┘  │
                                          └──────────────────────────────────┘
         │                     │                      │
         ▼                     ▼                      ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                        DATA & MESSAGING LAYER                                 │
│                                                                              │
│  ┌──────────────┐  ┌────────────────────────┐  ┌───────────┐  ┌──────────┐ │
│  │ PostgreSQL 16 │  │ Redis 7 (Multi-Layer)  │  │   NATS    │  │  MinIO   │ │
│  │ + PostGIS     │  │                        │  │ JetStream │  │          │ │
│  │ + TimescaleDB │  │ ├─ L2 Cache            │  │           │  │          │ │
│  │               │  │ │  ├─ Room Data (5min) │  │           │  │          │ │
│  │               │  │ │  ├─ User Prof (10min)│  │           │  │          │ │
│  │               │  │ │  ├─ POI (24hr)       │  │           │  │          │ │
│  │               │  │ │  ├─ Weather (30min)  │  │           │  │          │ │
│  │               │  │ │  ├─ Leaderboard (1hr)│  │           │  │          │ │
│  │               │  │ │  └─ AI Response (1hr)│  │           │  │          │ │
│  │               │  │ ├─ Rate Limiting       │  │           │  │          │ │
│  │               │  │ ├─ Session Store       │  │           │  │          │ │
│  │               │  │ ├─ Pub/Sub Bridge      │  │           │  │          │ │
│  │               │  │ ├─ Task Queue (Celery) │  │           │  │          │ │
│  │               │  │ └─ Rider Position Cache│  │           │  │          │ │
│  └──────────────┘  └────────────────────────┘  └───────────┘  └──────────┘ │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## Tech Stack Tambahan (Python)

| Komponen | Teknologi | Versi | Alasan |
|----------|-----------|-------|--------|
| Language | **Python** | 3.12+ | AI/ML ecosystem terbaik, Gemini SDK native |
| Web Framework | **FastAPI** | 0.115+ | Async native, auto-docs, type hints, SSE support |
| gRPC | **grpcio** + **grpcio-tools** | latest | Inter-service communication Go ↔ Python |
| Task Queue | **Celery** | 5.4+ | Distributed task processing, Redis broker |
| Beat Scheduler | **celery-beat** | latest | Cron scheduling untuk periodic tasks |
| AI/ML | **google-generativeai** | latest | Official Gemini Python SDK |
| Data Processing | **pandas** + **numpy** | latest | Analytics, statistics computation |
| Report Generation | **ReportLab** + **openpyxl** | latest | PDF/CSV/Excel report generation |
| HTTP Client | **httpx** | latest | Async HTTP untuk external API calls |
| ORM | **SQLAlchemy** + **asyncpg** | 2.0+ | Async PostgreSQL access |
| Redis Client | **redis-py** (async) | latest | Async Redis operations |
| NATS Client | **nats-py** | latest | Subscribe ke NATS events |
| Validation | **Pydantic** | 2.x | Data validation (built into FastAPI) |
| Testing | **pytest** + **pytest-asyncio** | latest | Async test support |
| Monitoring | **prometheus-fastapi-instrumentator** | latest | Metrics export |
| Logging | **structlog** | latest | Structured logging (consistent dengan Go Zap) |

---

## Proposed Changes

### Phase 1: Foundation & Infrastructure (Minggu 1-3)

#### Scope
Setup infrastruktur dasar, database, Redis caching layer, dan project scaffolding untuk Go + Python + Flutter + Next.js.

---

#### [NEW] Docker & Infrastructure

##### [NEW] docker-compose.yml
- Semua services: Go API, Go WS, Go Worker, Python AI, Python Analytics, Python Pipeline, PostgreSQL+PostGIS+TimescaleDB, Redis 7, NATS JetStream, LiveKit, MinIO, OSRM, Nginx, Next.js Admin
- Shared network configuration
- Volume mounts untuk persistent data
- Health check configuration

##### [NEW] docker/Dockerfile.go
- Multi-stage build: builder (Go 1.23) → runtime (alpine)
- Target: `api`, `ws`, `worker`

##### [NEW] docker/Dockerfile.python
- Multi-stage build: builder (Python 3.12) → runtime (slim)
- Target: `ai-service`, `analytics-service`, `data-pipeline`
- Install dependencies via `pip` + `requirements.txt`

##### [NEW] docker/nginx.conf
- Reverse proxy configuration
- SSL termination
- WebSocket upgrade (`/ws/*`)
- Rate limiting (global)
- Routing: `/api/*` → Go API, `/ws/*` → Go WS, `/ai/*` → Python AI Service, `/analytics/*` → Python Analytics Service
- Static file serving (MinIO proxy)

---

#### [NEW] Database Setup

##### [NEW] migrations/000001_init_extensions.up.sql
- Enable PostGIS, TimescaleDB, uuid-ossp, pgcrypto extensions

##### [NEW] migrations/000002_create_users.up.sql
- `users` table dengan semua fields dari System Brief §13.1
- `device_tokens` table
- `refresh_tokens` table
- Indexes sebagaimana didefinisikan di brief

##### [NEW] migrations/000003_create_motors.up.sql
- `motors` table dengan encrypted fields (pgcrypto)

##### [NEW] migrations/000004_create_rooms.up.sql
- `touring_rooms` table dengan PostGIS geometry columns
- `room_members` table
- `room_role_history` table
- Spatial indexes, partial indexes

##### [NEW] migrations/000005_create_routes.up.sql
- `routes` table dengan PostGIS LineString
- `waypoints` table

##### [NEW] migrations/000006_create_tracking.up.sql
- `rider_locations` table → TimescaleDB hypertable
- Compression policy (>7 days)
- Retention policy (>90 days)

##### [NEW] migrations/000007_create_emergency.up.sql
- `emergency_events` table
- `sos_events` table

##### [NEW] migrations/000008_create_chat.up.sql
- `chat_messages` table
- `message_reads` table

##### [NEW] migrations/000009_create_voting.up.sql
- `votings` table
- `voting_answers` table

##### [NEW] migrations/000010_create_management.up.sql
- `fuel_logs` table
- `expenses` table
- `service_reminders` table
- `checklist_templates` table
- `touring_checklists` table

##### [NEW] migrations/000011_create_social.up.sql
- `badges` table + seed data (14 badges)
- `user_badges` table
- `user_follows` table
- `user_blocks` table
- `reports` table
- `notifications` table
- `touring_posts` table
- `post_likes` table
- `post_comments` table

---

#### [NEW] Redis Caching Layer (Shared Go + Python)

##### Redis Key Design & TTL Strategy

```
# ============================================
# CACHING KEYS (L2 — Redis)
# ============================================

# User profiles (lazy loaded, cache-aside)
cache:user:{user_id}                    TTL: 10 min
cache:user:{user_id}:public_profile     TTL: 15 min
cache:user:{user_id}:badges             TTL: 1 hour
cache:user:{user_id}:stats              TTL: 30 min

# Room data (frequently accessed during touring)
cache:room:{room_id}                    TTL: 5 min
cache:room:{room_id}:members            TTL: 2 min
cache:room:{room_id}:settings           TTL: 5 min
cache:room:{room_id}:formation          TTL: 1 min
cache:room:{room_id}:active_route       TTL: 5 min

# Realtime rider positions (ultra-short TTL)
pos:room:{room_id}:rider:{user_id}      TTL: 120 sec (HASH: lat, lng, speed, heading, battery, ts)
pos:room:{room_id}:all                  TTL: 10 sec (aggregated positions — SET of rider IDs)

# POI data (lazy loaded per region)
cache:poi:region:{geohash_6}            TTL: 24 hours
cache:poi:nearby:{lat}:{lng}:{types}    TTL: 1 hour

# Weather data
cache:weather:current:{lat}:{lng}       TTL: 30 min
cache:weather:forecast:{lat}:{lng}      TTL: 1 hour
cache:weather:alerts:{lat}:{lng}        TTL: 15 min

# AI responses (semantic cache)
cache:ai:chat:{hash_of_prompt}          TTL: 1 hour
cache:ai:itinerary:{hash_of_params}     TTL: 6 hours
cache:ai:cost:{hash_of_params}          TTL: 6 hours

# Leaderboard (computed by Python analytics)
cache:leaderboard:monthly:{yyyy_mm}     TTL: 1 hour
cache:leaderboard:alltime               TTL: 1 hour

# Feed (paginated, lazy loaded)
cache:feed:user:{user_id}:page:{cursor} TTL: 5 min
cache:feed:explore:page:{cursor}        TTL: 5 min

# Notification count (badge)
cache:notif:unread:{user_id}            TTL: 5 min

# ============================================
# RATE LIMITING KEYS
# ============================================
rate:ip:{ip}:{endpoint}                 TTL: window duration
rate:user:{user_id}:{endpoint}          TTL: window duration
rate:ai:{user_id}                       TTL: 1 hour

# ============================================
# SESSION & PRESENCE KEYS
# ============================================
session:refresh:{token_hash}            TTL: 30 days
presence:heartbeat:{user_id}            TTL: 2 min
presence:room:{room_id}                 SET of online user_ids
ws:conn:{user_id}                       TTL: 5 min (WebSocket connection info)

# ============================================
# DISTRIBUTED LOCK KEYS (for Python workers)
# ============================================
lock:weather_sync:{region}              TTL: 5 min
lock:poi_refresh:{geohash}              TTL: 10 min
lock:badge_eval:{room_id}               TTL: 5 min
lock:analytics:{user_id}                TTL: 5 min
```

##### Cache Strategies per Data Type

| Data | Strategy | Write | Invalidation |
|------|----------|-------|-------------|
| User Profile | **Cache-Aside** (Lazy Load) | Write-Through | On profile update → delete key |
| Room Data | **Cache-Aside** (Lazy Load) | Write-Through | On room update → delete key |
| Rider Positions | **Write-Behind** | Batch write every 5s | Auto-expire (TTL 120s) |
| POI Data | **Cache-Aside** (Lazy Load) | Background refresh (Python pipeline) | TTL-based (24h) |
| Weather | **Cache-Aside** (Lazy Load) | Background refresh (Python pipeline) | TTL-based (30min) |
| AI Responses | **Cache-Aside** (Lazy Load) | On response generated | TTL-based (1h) |
| Leaderboard | **Refresh-Ahead** | Python analytics computes periodically | TTL-based (1h) + on touring complete |
| Feed | **Cache-Aside** (Lazy Load) | On new post → invalidate first page | TTL-based (5min) |
| Notifications | **Write-Through** | On notification create → increment counter | On read → decrement counter |

##### [NEW] go-road-backend/internal/repository/redis/cache_repo.go
```go
// Generic cache repository interface
type CacheRepository interface {
    Get(ctx context.Context, key string, dest interface{}) error
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, keys ...string) error
    Exists(ctx context.Context, key string) (bool, error)
    
    // Hash operations (for rider positions)
    HSet(ctx context.Context, key string, field string, value interface{}) error
    HGetAll(ctx context.Context, key string, dest interface{}) error
    
    // Set operations (for room presence)
    SAdd(ctx context.Context, key string, members ...interface{}) error
    SRem(ctx context.Context, key string, members ...interface{}) error
    SMembers(ctx context.Context, key string) ([]string, error)
    
    // Sorted Set (for leaderboard)
    ZAdd(ctx context.Context, key string, members ...redis.Z) error
    ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
    
    // Rate limiting
    IncrWithExpiry(ctx context.Context, key string, ttl time.Duration) (int64, error)
    
    // Distributed lock
    AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error)
    ReleaseLock(ctx context.Context, key string) error
}
```

##### [NEW] py-shared/redis_client.py
```python
# Shared Redis client for all Python microservices
# Async redis-py with connection pool
# Same key patterns as Go backend
# Cache decorator for automatic caching:
#   @cached(key_pattern="cache:weather:current:{lat}:{lng}", ttl=1800)
#   async def get_weather(lat: float, lng: float) -> WeatherData: ...
```

---

### Phase 2: Go Backend Core (Minggu 3-6)

#### Scope
Implementasi Go backend: auth, rooms, CRUD, middleware, pagination, dan lazy loading.

---

#### [NEW] Go Backend — Core

##### [NEW] go-road-backend/cmd/api/main.go
- Fiber v3 app bootstrap
- Middleware chain: Logger → CORS → RateLimit → Auth → Handler
- Graceful shutdown

##### [NEW] go-road-backend/internal/config/config.go
- Viper-based configuration loading
- Environment variable binding

##### [NEW] go-road-backend/internal/middleware/auth.go
- JWT validation middleware
- Extract user claims from token
- Role-based access control (global + room role)

##### [NEW] go-road-backend/internal/middleware/rate_limit.go
- Redis-based rate limiting
- Per-endpoint, per-IP, per-user limits
- `X-RateLimit-*` response headers

##### [NEW] go-road-backend/internal/middleware/cors.go
- CORS configuration with whitelisted origins

##### [NEW] go-road-backend/internal/middleware/logger.go
- Structured request/response logging via Zap
- Correlation ID per request

---

#### Pagination System (Cursor-Based)

> [!IMPORTANT]
> **Semua list endpoints** menggunakan cursor-based pagination. Tidak ada offset-based pagination.

##### [NEW] go-road-backend/internal/pkg/pagination/cursor.go

```go
// Cursor pagination implementation
// - Encodes/decodes cursor as base64(json({id, created_at}))
// - Supports forward pagination (after cursor)
// - Supports backward pagination (before cursor) — optional
// - Default limit: 20, max limit: 100

type PaginationParams struct {
    Cursor    string `query:"cursor"`
    Limit     int    `query:"limit" validate:"min=1,max=100"`
    Sort      string `query:"sort"`
    Order     string `query:"order" validate:"oneof=asc desc"`
}

type PaginatedResponse[T any] struct {
    Data    []T    `json:"data"`
    Meta    Meta   `json:"meta"`
}

type Meta struct {
    Cursor  string `json:"cursor,omitempty"`   // next page cursor
    HasMore bool   `json:"has_more"`
    Total   *int64 `json:"total,omitempty"`     // optional, expensive query
}

// Usage in repository:
// SELECT * FROM table WHERE (created_at, id) < (cursor_time, cursor_id)
// ORDER BY created_at DESC, id DESC
// LIMIT :limit + 1   -- +1 to check if there are more
```

##### Pagination per Endpoint

| Endpoint | Cursor Field | Default Limit | Max Limit | Cache |
|----------|-------------|---------------|-----------|-------|
| `GET /v1/rooms` | `created_at, id` | 20 | 50 | ❌ |
| `GET /v1/rooms/discover` | `touring_date, id` | 20 | 50 | ✅ 5min |
| `GET /v1/rooms/{id}/members` | `joined_at, id` | 50 | 100 | ✅ 2min |
| `GET /v1/rooms/{id}/messages` | `sent_at, id` | 30 | 100 | ❌ |
| `GET /v1/rooms/{id}/locations` | `recorded_at, id` | 50 | 200 | ❌ |
| `GET /v1/rooms/{id}/tracking` | `recorded_at` | 100 | 1000 | ✅ 10min |
| `GET /v1/fuel-logs` | `logged_at, id` | 20 | 50 | ❌ |
| `GET /v1/expenses` | `logged_at, id` | 20 | 50 | ❌ |
| `GET /v1/notifications` | `created_at, id` | 20 | 50 | ❌ |
| `GET /v1/feed` | `created_at, id` | 10 | 30 | ✅ 5min |
| `GET /v1/feed/explore` | `created_at, id` | 10 | 30 | ✅ 5min |
| `GET /v1/posts/{id}/comments` | `created_at, id` | 20 | 50 | ✅ 5min |
| `GET /v1/leaderboard` | `total_points` (ZSet) | 20 | 100 | ✅ 1hr |
| `GET /v1/conversations` | `last_message_at` | 20 | 50 | ❌ |
| `GET /v1/admin/users` | `created_at, id` | 50 | 100 | ❌ |
| `GET /v1/admin/rooms` | `created_at, id` | 50 | 100 | ❌ |
| `GET /v1/admin/reports` | `created_at, id` | 50 | 100 | ❌ |

---

#### Lazy Loading Strategy (Backend)

##### [NEW] go-road-backend/internal/pkg/lazy/loader.go

```go
// Lazy loading untuk Go backend:
// 1. GORM Preload hanya saat diminta (sparse fieldsets)
// 2. Query parameter ?include=members,route,weather untuk eager load
// 3. Default: hanya entity utama, relasi di-lazy load via separate endpoints

// Example: GET /v1/rooms/{id}?include=members,active_route
// - Tanpa ?include: hanya room data
// - ?include=members: + room_members (paginated, first 20)
// - ?include=active_route: + active route + waypoints
// - ?include=weather: + weather data dari cache/API
```

##### Lazy Loading Rules per Domain

| Entity | Default Fields | Lazy Loaded (via ?include) | Separate Endpoint |
|--------|---------------|---------------------------|-------------------|
| Room | id, name, status, dates, location_names | `members`, `active_route`, `weather`, `settings` | `/rooms/{id}/members`, `/rooms/{id}/routes` |
| User Profile | id, name, photo, skill, role | `motors`, `badges`, `stats`, `recent_touring` | `/riders/{id}/profile`, `/me/badges`, `/me/stats` |
| Room Member | id, user_id, role, position | `user_profile`, `motor` | Embedded in member response |
| Chat Message | id, content, sender_id, sent_at | `sender_profile`, `read_receipts` | `/rooms/{id}/messages/{mid}` |
| Touring Post | id, caption, photos, stats_snapshot | `comments` (paginated), `author_profile` | `/posts/{id}/comments` |
| Route | id, name, distance, duration | `waypoints`, `elevation_profile`, `weather_per_waypoint` | `/rooms/{id}/routes/{rid}` |

---

#### [NEW] Go Backend — Domain, Service, Repository, Handler (per domain)

Untuk setiap domain berikut, buat:
- `internal/domain/{domain}/entity.go` — Go struct entities
- `internal/domain/{domain}/repository.go` — Repository interface
- `internal/domain/{domain}/service.go` — Service interface
- `internal/service/{domain}_service.go` — Service implementation (business logic + cache integration)
- `internal/repository/postgres/{domain}_repo.go` — GORM repository implementation (with cursor pagination)
- `internal/handler/{domain}_handler.go` — Fiber HTTP handler (with pagination params parsing)
- `internal/dto/request/{domain}_request.go` — Request DTOs with validator tags
- `internal/dto/response/{domain}_response.go` — Response DTOs (filtered fields, lazy loaded relations)

**Domains:**

| # | Domain | Key Features |
|---|--------|-------------|
| 1 | `auth` | Register, login, JWT, refresh, profile CRUD, avatar upload |
| 2 | `room` | CRUD room, join/leave, member management, role assignment, status machine |
| 3 | `convoy` | Formation CRUD, location tracking (REST fallback), speed monitoring |
| 4 | `route` | Route CRUD, waypoints, GPX import/export, OSRM integration |
| 5 | `emergency` | Emergency event CRUD, SOS trigger/dismiss, severity-based handling |
| 6 | `chat` | Messages CRUD, pin, delete, read receipts, private messages |
| 7 | `voting` | Voting CRUD, vote submission, results computation |
| 8 | `weather` | Current/forecast/route weather, alerts (delegates to Python) |
| 9 | `ai` | Chat/itinerary/cost/route/packing (delegates to Python via gRPC) |
| 10 | `motor` | Motor CRUD, primary selection, encrypted fields |
| 11 | `fuel` | Fuel log CRUD, analytics (delegates to Python) |
| 12 | `expense` | Expense CRUD, summary, export (delegates to Python) |
| 13 | `notification` | Notification list, read/unread, preferences, FCM push |
| 14 | `poi` | Nearby POI search, report POI, cache management |
| 15 | `badge` | Badge list, user badges, trigger evaluation (delegates to Python) |
| 16 | `social` | Follow/unfollow, block, report, feed, posts, likes, comments |
| 17 | `checklist` | Template CRUD, touring checklist CRUD |
| 18 | `qr` | QR card generation/scan |
| 19 | `service_reminder` | Reminder CRUD, completion tracking |
| 20 | `upload` | File upload to MinIO, virus scan, type validation |
| 21 | `admin` | Dashboard stats, user management, room monitoring, report moderation |

##### [NEW] go-road-backend/internal/handler/routes.go
- All API route definitions (Fiber router groups)
- Middleware binding per group (auth, role check)
- Swagger annotations per endpoint

---

### Phase 3: Python Microservices (Minggu 5-8)

> [!IMPORTANT]
> Python services berjalan parallel dengan Go backend development. Komunikasi via **gRPC** (synchronous) dan **NATS** (event-driven asynchronous).

---

#### [NEW] Proto Definitions (Shared Go + Python)

##### [NEW] proto/ai_service.proto
```protobuf
service AIService {
    rpc ChatStream(ChatRequest) returns (stream ChatResponse);
    rpc GenerateItinerary(ItineraryRequest) returns (ItineraryResponse);
    rpc EstimateCost(CostEstimationRequest) returns (CostEstimationResponse);
    rpc AdviseRoute(RouteAdviceRequest) returns (RouteAdviceResponse);
    rpc GeneratePackingList(PackingListRequest) returns (PackingListResponse);
    rpc AdviseSafety(SafetyRequest) returns (SafetyResponse);
    rpc RecommendPOI(POIRecommendRequest) returns (POIRecommendResponse);
}
```

##### [NEW] proto/analytics_service.proto
```protobuf
service AnalyticsService {
    rpc ComputeRiderStats(RiderStatsRequest) returns (RiderStatsResponse);
    rpc ComputeLeaderboard(LeaderboardRequest) returns (LeaderboardResponse);
    rpc EvaluateBadges(BadgeEvalRequest) returns (BadgeEvalResponse);
    rpc GenerateReport(ReportRequest) returns (ReportResponse);
    rpc ComputeRoomStats(RoomStatsRequest) returns (RoomStatsResponse);
    rpc ComputeFuelAnalytics(FuelAnalyticsRequest) returns (FuelAnalyticsResponse);
    rpc ComputeExpenseSummary(ExpenseSummaryRequest) returns (ExpenseSummaryResponse);
    rpc GetAdminDashboard(AdminDashRequest) returns (AdminDashResponse);
}
```

---

#### [NEW] py-ai-service/ — AI Assistant Microservice

##### [NEW] py-ai-service/requirements.txt
```
fastapi[standard]>=0.115.0
uvicorn[standard]>=0.30.0
grpcio>=1.66.0
grpcio-tools>=1.66.0
google-generativeai>=0.8.0
redis[hiredis]>=5.0.0
pydantic>=2.0
pydantic-settings>=2.0
structlog>=24.0
httpx>=0.27.0
prometheus-fastapi-instrumentator>=7.0
```

##### [NEW] py-ai-service/app/main.py
- FastAPI app dengan health check, metrics endpoint
- gRPC server (concurrent dengan FastAPI HTTP)
- Redis connection pool initialization
- Structured logging setup

##### [NEW] py-ai-service/app/services/gemini_service.py
- Google Gemini API integration
- **Streaming response** via SSE (Server-Sent Events) untuk chat
- Context building: rute, cuaca, jumlah rider, motor data
- **Redis caching** untuk responses (semantic hash key)
- Rate limiting: 20 req/hr per user (Redis counter)
- Graceful degradation: return "AI tidak tersedia" saat offline

##### [NEW] py-ai-service/app/services/itinerary_service.py
- Input: rute + durasi + jumlah rider + motor data
- Output: jadwal touring detail (waktu, stop, estimasi, cuaca)
- **Redis cache**: hash(route_id + params) → cached itinerary (TTL 6h)

##### [NEW] py-ai-service/app/services/cost_service.py
- Input: rute + motor + jumlah rider
- Output: estimasi biaya (BBM, makan, hotel, tol, parkir)
- Data harga BBM dari cache/API
- **Redis cache**: hash(route_id + motor_ids) → cached estimate (TTL 6h)

##### [NEW] py-ai-service/app/services/route_advisor_service.py
- Input: origin + destination + preferences
- Output: alternatif rute dengan pro/cons
- Integration dengan OSRM untuk route calculation

##### [NEW] py-ai-service/app/grpc_server.py
- gRPC server implementation untuk `AIService`
- Delegates ke service classes

##### [NEW] py-ai-service/app/cache.py
```python
# Redis cache decorator
# @cached(key="cache:ai:chat:{prompt_hash}", ttl=3600)
# Supports cache invalidation, cache warming
# Lazy loading pattern: check cache first → compute if miss → store in cache
```

##### [NEW] py-ai-service/Dockerfile
- Python 3.12 slim
- Multi-stage build

---

#### [NEW] py-analytics-service/ — Analytics & Reporting Microservice

##### [NEW] py-analytics-service/requirements.txt
```
fastapi[standard]>=0.115.0
uvicorn[standard]>=0.30.0
grpcio>=1.66.0
grpcio-tools>=1.66.0
sqlalchemy[asyncio]>=2.0
asyncpg>=0.29.0
redis[hiredis]>=5.0.0
pandas>=2.2.0
numpy>=2.0.0
reportlab>=4.2
openpyxl>=3.1.0
pydantic>=2.0
pydantic-settings>=2.0
structlog>=24.0
prometheus-fastapi-instrumentator>=7.0
```

##### [NEW] py-analytics-service/app/main.py
- FastAPI + gRPC server
- PostgreSQL async connection pool (asyncpg)
- Redis connection pool

##### [NEW] py-analytics-service/app/services/stats_service.py
- **Rider Statistics** computation:
  - Total KM, total touring, total riding time
  - Average speed, most used motor
  - Most active month, favorite route
  - Achievement progress
- **Redis cache**: `cache:user:{id}:stats` (TTL 30min)
- **Lazy loading**: compute on first request, cache result
- Heavy queries run against read replica (if available)

##### [NEW] py-analytics-service/app/services/leaderboard_service.py
- **Leaderboard** computation:
  - Monthly: top riders by points earned this month
  - All-time: top riders by total points
- Uses Redis **Sorted Set** (`ZADD`, `ZREVRANGE`)
- **Refresh-Ahead cache**: recomputed every hour by Celery beat
- **Pagination**: `ZREVRANGEBYSCORE` with cursor

##### [NEW] py-analytics-service/app/services/badge_service.py
- **Badge Evaluation** engine:
  - Input: touring data (distance, duration, elevation, time, weather, etc.)
  - Evaluate all 14 badge criteria
  - Award new badges
- Triggered by NATS event `touring.completed`
- **Redis cache**: `cache:user:{id}:badges` (TTL 1hr)

##### [NEW] py-analytics-service/app/services/report_service.py
- **Report Generation**:
  - Fuel analytics: charts data (line, bar, pie)
  - Expense summary: per category, per touring, monthly
  - Touring report: route, participants, statistics
  - Admin report: user growth, active rooms, emergency events
- PDF generation via ReportLab
- CSV/Excel via openpyxl
- Reports stored in MinIO, URL returned
- **Redis cache**: report URL (TTL 1hr)

##### [NEW] py-analytics-service/app/services/fuel_analytics_service.py
- Konsumsi rata-rata per motor
- Tren konsumsi over time (pandas DataFrame → JSON chart data)
- Perbandingan antar jenis BBM
- Total pengeluaran BBM per bulan
- Best/worst fuel efficiency per touring

##### [NEW] py-analytics-service/app/services/expense_summary_service.py
- Pie chart per kategori (pandas groupby → JSON)
- Trend bulanan (pandas resample → JSON)
- Perbandingan antar touring
- Split bill calculation

##### [NEW] py-analytics-service/app/grpc_server.py
- gRPC server implementation untuk `AnalyticsService`

---

#### [NEW] py-data-pipeline/ — Background Workers & Data Pipeline

##### [NEW] py-data-pipeline/requirements.txt
```
celery[redis]>=5.4.0
redis[hiredis]>=5.0.0
sqlalchemy[asyncio]>=2.0
asyncpg>=0.29.0
httpx>=0.27.0
nats-py>=2.7.0
pandas>=2.2.0
numpy>=2.0.0
scikit-learn>=1.5.0
structlog>=24.0
pydantic>=2.0
pydantic-settings>=2.0
```

##### [NEW] py-data-pipeline/app/celery_app.py
- Celery app configuration
- Redis as broker and result backend
- Beat schedule for periodic tasks:

```python
beat_schedule = {
    # Weather sync: every 30 minutes for active rooms
    'weather-sync': {
        'task': 'tasks.weather_sync',
        'schedule': crontab(minute='*/30'),
    },
    # POI cache refresh: every 24 hours
    'poi-cache-refresh': {
        'task': 'tasks.poi_cache_refresh',
        'schedule': crontab(hour=3, minute=0),  # 3 AM
    },
    # Leaderboard recompute: every 1 hour
    'leaderboard-recompute': {
        'task': 'tasks.leaderboard_recompute',
        'schedule': crontab(minute=0),  # every hour
    },
    # Service reminder check: every day at 8 AM
    'service-reminder-check': {
        'task': 'tasks.service_reminder_check',
        'schedule': crontab(hour=8, minute=0),
    },
    # Data retention cleanup: every day at 2 AM
    'data-retention-cleanup': {
        'task': 'tasks.data_retention_cleanup',
        'schedule': crontab(hour=2, minute=0),
    },
    # Pre-touring weather notification: every 15 min
    'pre-touring-weather': {
        'task': 'tasks.pre_touring_weather_notify',
        'schedule': crontab(minute='*/15'),
    },
}
```

##### [NEW] py-data-pipeline/app/tasks/weather_sync.py
- Fetch weather data from OpenWeather API untuk semua active rooms
- Store in **Redis cache** (TTL 30min)
- Alert cuaca ekstrem → publish NATS event → Go broadcast ke room
- **Distributed lock**: Redis lock per region untuk prevent duplicate fetches

##### [NEW] py-data-pipeline/app/tasks/poi_cache_refresh.py
- Fetch POI data dari Overpass API per region (geohash-6)
- Store in **Redis cache** (TTL 24hr)
- Also store in PostgreSQL for backup
- **Lazy approach**: only refresh regions with active rooms or recent queries

##### [NEW] py-data-pipeline/app/tasks/smart_detection.py
- **ML-enhanced straggler detection**:
  - Input: rider positions (from NATS subscription)
  - Calculate group centroid, distance from centroid per rider
  - Detect anomalies using scikit-learn IsolationForest
  - Apply sliding window to avoid false positives
  - Cooldown period tracking via Redis
- Publish alerts via NATS → Go broadcasts to room
- **Note**: Basic detection still runs in Go workers (threshold-based). Python ML model is for enhanced detection (Phase 5+).

##### [NEW] py-data-pipeline/app/tasks/leaderboard_recompute.py
- Query PostgreSQL for all user points
- Compute monthly and all-time rankings
- Store in Redis Sorted Set (`cache:leaderboard:*`)

##### [NEW] py-data-pipeline/app/tasks/service_reminder_check.py
- Query `service_reminders` table
- Check if any reminder is due (H-7 or H-1)
- Publish notification events via NATS

##### [NEW] py-data-pipeline/app/tasks/data_retention.py
- TimescaleDB handles location data retention automatically
- This task handles: chat messages > 1 year, server logs > 30 days, analytics raw data > 2 years
- Soft-deleted accounts > 30 days → hard delete

##### [NEW] py-data-pipeline/app/tasks/feed_aggregator.py
- When touring completes → generate touring post data
- Compute stats_snapshot (distance, duration, participants, etc.)
- Generate route_snapshot (simplified polyline)

##### [NEW] py-data-pipeline/app/nats_subscriber.py
- Subscribe to NATS events:
  - `touring.completed` → trigger badge evaluation + feed aggregation
  - `room.{id}.location` → feed to smart detection ML model (batch)
  - `weather.alert` → process and notify
  - `user.profile_updated` → invalidate Redis cache

---

### Phase 4: Go Backend — Realtime (Minggu 6-9)

#### Scope
WebSocket gateway, NATS pub/sub, background workers, location streaming, smart detection.

---

##### [NEW] go-road-backend/cmd/ws/main.go
- WebSocket server entry point
- gorilla/websocket upgrade handler
- JWT authentication on connect

##### [NEW] go-road-backend/internal/ws/hub.go
- Connection manager (concurrent map of user_id → *Client)
- Room subscription management
- Message routing to room members
- Heartbeat management

##### [NEW] go-road-backend/internal/ws/client.go
- Single WebSocket client handler
- Read/write goroutines
- Message parsing and routing

##### [NEW] go-road-backend/internal/ws/location_handler.go
- Process incoming location updates
- Publish to NATS: `room.{room_id}.location`
- Update Redis: `pos:room:{room_id}:rider:{user_id}`

##### [NEW] go-road-backend/internal/ws/chat_handler.go
- Realtime chat message delivery
- Typing indicators
- Read receipts

##### [NEW] go-road-backend/internal/ws/voting_handler.go
- Realtime voting updates
- Vote submission + result broadcast

##### [NEW] go-road-backend/internal/ws/presence_handler.go
- Online/offline status tracking
- Heartbeat processing (update Redis TTL)

##### [NEW] go-road-backend/internal/event/publisher.go
- NATS JetStream publisher
- Publish events: location, chat, emergency, voting, presence, convoy, system

##### [NEW] go-road-backend/internal/event/subscriber.go
- NATS JetStream subscriber
- Subscribe to: location fan-out, Python service events

##### [NEW] go-road-backend/cmd/worker/main.go
- Asynq worker entry point
- Register task handlers

##### [NEW] go-road-backend/internal/worker/location_aggregator.go
- Batch insert locations to TimescaleDB (every 5s)
- Collect from NATS, batch, bulk insert

##### [NEW] go-road-backend/internal/worker/smart_detection.go
- Threshold-based detection (running in Go for low latency):
  - Straggler: > 1.5km from centroid
  - Stopped too long: > 5min
  - Off route: > 500m from polyline
  - Offline: > 2min no update
  - Speed limit exceeded
  - Battery low: < 15%
- Sliding window + cooldown via Redis
- Publish alerts via NATS + FCM push

##### [NEW] go-road-backend/internal/worker/notification_sender.go
- FCM push notification sender
- Batch send with retry

---

### Phase 5: Flutter Mobile App (Minggu 8-14)

#### Scope
Full Flutter app implementation with lazy loading, pagination, caching, and offline support.

---

#### Flutter Lazy Loading Strategy

```
┌─────────────────────────────────────────────────────┐
│              FLUTTER LAZY LOADING                     │
│                                                       │
│  1. IMAGE LAZY LOADING                                │
│     ├─ CachedNetworkImage (cached_network_image)     │
│     ├─ Placeholder shimmer while loading              │
│     ├─ Progressive image loading (thumbnail → full)  │
│     └─ Lazy load images in ListView (viewport only)  │
│                                                       │
│  2. LIST LAZY LOADING (Infinite Scroll)               │
│     ├─ Riverpod AsyncNotifier + cursor pagination    │
│     ├─ ScrollController listener (80% threshold)     │
│     ├─ Load next page when near bottom               │
│     ├─ Loading indicator at bottom                   │
│     └─ Pull-to-refresh (reset cursor)                │
│                                                       │
│  3. MAP LAZY LOADING                                  │
│     ├─ Tile lazy loading (flutter_map default)       │
│     ├─ Marker clustering (collapse when zoomed out)  │
│     ├─ POI lazy load per viewport (geohash query)    │
│     └─ Route polyline simplification at zoom levels  │
│                                                       │
│  4. DATA LAZY LOADING                                 │
│     ├─ Riverpod family providers (load on demand)    │
│     ├─ Room detail: basic info first, lazy load      │
│     │   members/routes/weather on tab switch         │
│     ├─ Profile: basic info first, lazy load          │
│     │   badges/stats/motors on scroll                │
│     └─ Chat: load last 30 messages, lazy load        │
│         older messages on scroll up                   │
│                                                       │
│  5. SCREEN LAZY LOADING                               │
│     ├─ GoRouter lazy route building                  │
│     ├─ Deferred imports for heavy screens            │
│     └─ Splash screen while initializing              │
│                                                       │
│  6. OFFLINE-FIRST LAZY LOADING                        │
│     ├─ Show cached data immediately (Drift)          │
│     ├─ Fetch fresh data in background                │
│     ├─ Update UI when fresh data arrives             │
│     └─ Stale indicator for old cached data           │
└─────────────────────────────────────────────────────┘
```

##### [NEW] go_road_app/lib/core/network/api_client.dart
- Dio instance dengan interceptors:
  - Auth interceptor (attach JWT, auto-refresh)
  - **Cache interceptor** (check Drift cache → return if fresh, fetch if stale)
  - Retry interceptor (exponential backoff)
  - Connectivity interceptor (queue if offline)
  - Logging interceptor

##### [NEW] go_road_app/lib/core/pagination/paginated_list.dart
```dart
/// Generic paginated list widget with lazy loading
/// - Uses Riverpod AsyncNotifier for state
/// - ScrollController with 80% threshold for next page
/// - Shimmer loading placeholder
/// - Error state with retry
/// - Empty state
/// - Pull-to-refresh (reset cursor)

class PaginatedListView<T> extends ConsumerWidget {
  final AutoDisposeAsyncNotifierProvider<PaginatedNotifier<T>, PaginatedState<T>> provider;
  final Widget Function(T item) itemBuilder;
  final Widget? emptyWidget;
  final Widget? loadingWidget;
  final Widget? errorWidget;
  // ...
}

class PaginatedState<T> {
  final List<T> items;
  final String? nextCursor;
  final bool hasMore;
  final bool isLoadingMore;
}
```

##### [NEW] go_road_app/lib/core/cache/cache_manager.dart
```dart
/// Multi-layer cache manager:
/// - L1: In-memory (LRU cache, TTL 30s)
/// - L2: Drift SQLite (TTL varies per data type)
/// - L3: Network (API call)
///
/// Flow: L1 → L2 → L3 → store back to L1 + L2
/// Lazy loading: only fetch from L3 when L1 + L2 miss
```

##### [NEW] go_road_app/lib/core/lazy/lazy_image.dart
```dart
/// Lazy-loaded image widget with:
/// - Shimmer placeholder
/// - Progressive loading (thumbnail → full)
/// - Cached via cached_network_image
/// - Offline: show cached version or placeholder
```

##### [NEW] go_road_app/lib/data/datasources/local/database.dart
- Drift database class
- All table definitions
- DAO classes per domain
- Migration strategy
- Offline queue table

##### [NEW] go_road_app/lib/data/datasources/local/offline_queue.dart
- Priority-based offline queue (SOS=1, Emergency=2, Location=3, Chat=4, Others=5)
- Sync manager: process queue on reconnect

##### Flutter Screens (all with lazy loading + pagination where applicable)

| Screen | Lazy Loading | Pagination | Caching |
|--------|-------------|------------|---------|
| Splash | Prefetch critical data | — | — |
| Home (Map) | Map tiles, markers by viewport | — | Drift (map tiles, POI) |
| Room List | Image lazy load | Cursor-based infinite scroll | Drift (room list) |
| Room Detail | Tabs: members/routes/chat lazy loaded per tab | Members list paginated | Drift (room data) |
| Chat | Messages lazy loaded on scroll up | Cursor-based (newer/older) | Drift (messages) |
| Convoy Map | Markers lazy loaded by viewport | — | Redis positions via WS |
| Route Planner | Waypoints lazy load, map tiles | — | Drift (routes) |
| Feed | Posts with lazy images | Cursor-based infinite scroll | Drift (feed) |
| Profile | Tabs: stats/badges/motors lazy loaded | — | Drift (profile) |
| Leaderboard | — | Cursor-based infinite scroll | Drift (leaderboard) |
| Notifications | — | Cursor-based infinite scroll | Drift (notifications) |
| Fuel Logs | — | Cursor-based infinite scroll | Drift (fuel logs) |
| Expenses | — | Cursor-based infinite scroll | Drift (expenses) |
| AI Chat | Streaming response (SSE) | Chat history paginated | Drift (AI conversations) |
| Motor List | Photos lazy load | — | Drift (motors) |
| Settings | — | — | Drift (settings) |

---

### Phase 6: Next.js Admin Dashboard (Minggu 12-15)

#### Scope
Admin panel with lazy loading, server-side pagination, dan TanStack Query caching.

---

##### Admin Lazy Loading + Pagination

```typescript
// TanStack Query with cursor-based pagination
const { data, fetchNextPage, hasNextPage, isFetchingNextPage } =
  useInfiniteQuery({
    queryKey: ['admin', 'users'],
    queryFn: ({ pageParam }) => fetchUsers({ cursor: pageParam, limit: 50 }),
    getNextPageParam: (lastPage) => lastPage.meta.has_more ? lastPage.meta.cursor : undefined,
    staleTime: 5 * 60 * 1000,  // 5 min cache
  });

// Lazy loaded components
const RoomLiveMap = lazy(() => import('./components/RoomLiveMap'));
const AnalyticsCharts = lazy(() => import('./components/AnalyticsCharts'));
```

##### [NEW] Admin Pages (all with TanStack Table pagination)

| Page | Features | Pagination | Caching |
|------|----------|------------|---------|
| Dashboard | Stats cards, charts (lazy loaded) | — | TanStack Query 5min |
| Users | Table with search, sort, filter | Server-side cursor | TanStack Query 2min |
| User Detail | Profile, touring history, badges | Sub-lists paginated | TanStack Query 5min |
| Rooms | Table with status filter | Server-side cursor | TanStack Query 2min |
| Room Live | Live map (Leaflet, lazy loaded) | — | WebSocket realtime |
| Reports | Report queue | Server-side cursor | TanStack Query 2min |
| Analytics | Charts (Recharts, lazy loaded) | — | TanStack Query 30min |
| Emergency Log | Event log | Server-side cursor | TanStack Query 1min |
| Moderation | Report review | Server-side cursor | TanStack Query 2min |
| Settings | Feature flags, app config | — | TanStack Query 5min |

---

### Phase 7: Voice, Testing, DevOps (Minggu 14-18)

---

#### Voice (LiveKit PTT)
- Go: LiveKit token generation endpoint
- Flutter: livekit_client integration, PTT button, channel selector
- Voice priority system (server-side via LiveKit metadata)
- Background audio via foreground service

#### Testing

| Component | Framework | Coverage Target |
|-----------|-----------|----------------|
| Go Backend | testify + mockery | 80%+ |
| Python AI Service | pytest + pytest-asyncio | 75%+ |
| Python Analytics | pytest + pandas testing | 75%+ |
| Python Pipeline | pytest + celery testing | 70%+ |
| Flutter | flutter_test + mocktail | 70%+ |
| Next.js Admin | Jest + Testing Library + Playwright | 60%+ |

#### DevOps
- GitHub Actions CI/CD (as specified in System Brief §15.4)
- Docker multi-stage builds for all services
- Prometheus + Grafana + Loki + Sentry setup
- Backup & DR procedures

---

## Folder Structure Tambahan (Python)

```
go-road/
├── go-road-backend/          # Go backend (as defined in System Brief §12.2)
│   └── ...
├── go_road_app/              # Flutter app (as defined in System Brief §12.3)
│   └── ...
├── go-road-admin/            # Next.js admin (as defined in System Brief §12.4)
│   └── ...
├── proto/                    # Shared gRPC proto definitions
│   ├── ai_service.proto
│   ├── analytics_service.proto
│   └── common.proto
├── py-ai-service/            # Python AI Microservice
│   ├── app/
│   │   ├── __init__.py
│   │   ├── main.py           # FastAPI + gRPC server
│   │   ├── config.py         # Pydantic Settings
│   │   ├── cache.py          # Redis cache decorators
│   │   ├── grpc_server.py    # gRPC service implementation
│   │   ├── services/
│   │   │   ├── gemini_service.py
│   │   │   ├── itinerary_service.py
│   │   │   ├── cost_service.py
│   │   │   ├── route_advisor_service.py
│   │   │   └── safety_service.py
│   │   ├── models/           # Pydantic models
│   │   │   ├── requests.py
│   │   │   └── responses.py
│   │   └── proto/            # Generated gRPC code
│   │       ├── ai_service_pb2.py
│   │       └── ai_service_pb2_grpc.py
│   ├── tests/
│   ├── requirements.txt
│   ├── Dockerfile
│   └── README.md
├── py-analytics-service/     # Python Analytics Microservice
│   ├── app/
│   │   ├── __init__.py
│   │   ├── main.py           # FastAPI + gRPC server
│   │   ├── config.py
│   │   ├── cache.py
│   │   ├── grpc_server.py
│   │   ├── database.py       # SQLAlchemy async engine
│   │   ├── services/
│   │   │   ├── stats_service.py
│   │   │   ├── leaderboard_service.py
│   │   │   ├── badge_service.py
│   │   │   ├── report_service.py
│   │   │   ├── fuel_analytics_service.py
│   │   │   └── expense_summary_service.py
│   │   ├── models/
│   │   └── proto/
│   ├── tests/
│   ├── requirements.txt
│   ├── Dockerfile
│   └── README.md
├── py-data-pipeline/         # Python Data Pipeline (Celery Workers)
│   ├── app/
│   │   ├── __init__.py
│   │   ├── celery_app.py     # Celery config + beat schedule
│   │   ├── config.py
│   │   ├── cache.py
│   │   ├── database.py
│   │   ├── nats_subscriber.py  # NATS event listener
│   │   ├── tasks/
│   │   │   ├── weather_sync.py
│   │   │   ├── poi_cache_refresh.py
│   │   │   ├── smart_detection.py      # ML-enhanced detection
│   │   │   ├── leaderboard_recompute.py
│   │   │   ├── service_reminder_check.py
│   │   │   ├── data_retention.py
│   │   │   ├── feed_aggregator.py
│   │   │   └── pre_touring_weather.py
│   │   └── ml/
│   │       ├── straggler_model.py      # scikit-learn IsolationForest
│   │       └── anomaly_detector.py
│   ├── tests/
│   ├── requirements.txt
│   ├── Dockerfile
│   └── README.md
├── py-shared/                # Shared Python utilities
│   ├── redis_client.py       # Shared async Redis client
│   ├── nats_client.py        # Shared NATS client
│   ├── logging_config.py     # Shared structured logging
│   └── common_models.py      # Shared Pydantic models
├── docker-compose.yml        # All services
├── docker-compose.dev.yml    # Development overrides
├── .env.example
└── README.md
```

---

## Redis Caching — Summary Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                    REDIS 7 — MULTI-PURPOSE                           │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  L2 CACHE (Cache-Aside / Lazy Loading)                       │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────────────┐ │   │
│  │  │ User Profile  │ │ Room Data    │ │ Weather Data          │ │   │
│  │  │ TTL: 10min    │ │ TTL: 5min    │ │ TTL: 30min            │ │   │
│  │  └──────────────┘ └──────────────┘ └──────────────────────┘ │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────────────┐ │   │
│  │  │ POI Data      │ │ AI Responses │ │ Feed Pages            │ │   │
│  │  │ TTL: 24hr     │ │ TTL: 1hr     │ │ TTL: 5min             │ │   │
│  │  └──────────────┘ └──────────────┘ └──────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  REALTIME DATA                                                │   │
│  │  ┌──────────────────┐ ┌────────────────┐ ┌────────────────┐ │   │
│  │  │ Rider Positions   │ │ Room Presence   │ │ Heartbeat      │ │   │
│  │  │ HASH, TTL: 120s   │ │ SET             │ │ TTL: 2min      │ │   │
│  │  └──────────────────┘ └────────────────┘ └────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  LEADERBOARD (Refresh-Ahead by Python Analytics)              │   │
│  │  ┌──────────────────────┐ ┌──────────────────────┐          │   │
│  │  │ Monthly Leaderboard   │ │ All-time Leaderboard  │          │   │
│  │  │ ZSET, TTL: 1hr        │ │ ZSET, TTL: 1hr        │          │   │
│  │  └──────────────────────┘ └──────────────────────┘          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────┐ ┌──────────────┐ ┌──────────────────┐   │
│  │ Rate Limiting        │ │ Session Store │ │ Distributed Locks│   │
│  │ INCR + EXPIRE        │ │ STRING        │ │ SET NX + EXPIRE  │   │
│  └─────────────────────┘ └──────────────┘ └──────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  TASK QUEUE (Celery Broker + Result Backend)                  │   │
│  │  Used by: py-data-pipeline (Celery workers)                   │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  PUB/SUB BRIDGE                                               │   │
│  │  Used by: cache invalidation broadcast across Go instances    │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Pagination — End-to-End Flow

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Flutter     │     │   Go API     │     │   Redis      │     │  PostgreSQL  │
│   Client      │     │   (Fiber)    │     │   Cache      │     │              │
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                    │                    │
       │ GET /v1/feed       │                    │                    │
       │ ?cursor=abc&limit=10                    │                    │
       │───────────────────▶│                    │                    │
       │                    │                    │                    │
       │                    │ Check cache        │                    │
       │                    │ cache:feed:user:   │                    │
       │                    │  {id}:page:{cursor}│                    │
       │                    │───────────────────▶│                    │
       │                    │                    │                    │
       │                    │   CACHE MISS       │                    │
       │                    │◀───────────────────│                    │
       │                    │                    │                    │
       │                    │ SELECT * FROM touring_posts             │
       │                    │ WHERE (created_at, id) <               │
       │                    │   (cursor_time, cursor_id)             │
       │                    │ ORDER BY created_at DESC, id DESC      │
       │                    │ LIMIT 11  (limit+1)│                    │
       │                    │───────────────────────────────────────▶│
       │                    │                    │                    │
       │                    │        11 rows (has_more = true)       │
       │                    │◀───────────────────────────────────────│
       │                    │                    │                    │
       │                    │ Store in cache     │                    │
       │                    │ TTL: 5min          │                    │
       │                    │───────────────────▶│                    │
       │                    │                    │                    │
       │  {                 │                    │                    │
       │   data: [10 posts],│                    │                    │
       │   meta: {          │                    │                    │
       │     cursor: "xyz", │                    │                    │
       │     has_more: true │                    │                    │
       │   }                │                    │                    │
       │  }                 │                    │                    │
       │◀───────────────────│                    │                    │
       │                    │                    │                    │
       │ [User scrolls to bottom]               │                    │
       │ ScrollController detects 80% threshold  │                    │
       │                    │                    │                    │
       │ GET /v1/feed       │                    │                    │
       │ ?cursor=xyz&limit=10                    │                    │
       │───────────────────▶│                    │                    │
       │                    │ ... (same flow)    │                    │
       │                    │                    │                    │
```

---

## Lazy Loading — End-to-End Flow

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Flutter     │     │   Go API     │     │   Redis      │     │  PostgreSQL  │
│   Client      │     │   (Fiber)    │     │   Cache      │     │              │
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                    │                    │
       │ [User opens Room Detail screen]        │                    │
       │                    │                    │                    │
       │ GET /v1/rooms/{id} │ (basic data only)  │                    │
       │───────────────────▶│                    │                    │
       │                    │ cache:room:{id}    │                    │
       │                    │───────────────────▶│                    │
       │                    │   CACHE HIT ✅      │                    │
       │                    │◀───────────────────│                    │
       │  Room basic data   │                    │                    │
       │◀───────────────────│                    │                    │
       │                    │                    │                    │
       │ [Render room header, show skeleton for tabs]                │
       │                    │                    │                    │
       │ [User taps "Members" tab] ← LAZY LOAD  │                    │
       │                    │                    │                    │
       │ GET /v1/rooms/{id}/members?limit=20     │                    │
       │───────────────────▶│                    │                    │
       │                    │ cache:room:{id}:   │                    │
       │                    │  members           │                    │
       │                    │───────────────────▶│                    │
       │                    │   CACHE MISS       │                    │
       │                    │◀───────────────────│                    │
       │                    │                    │ query DB           │
       │                    │───────────────────────────────────────▶│
       │                    │◀───────────────────────────────────────│
       │                    │ Store in cache     │                    │
       │                    │───────────────────▶│                    │
       │  Members list (paginated)              │                    │
       │◀───────────────────│                    │                    │
       │                    │                    │                    │
       │ [User taps "Route" tab] ← LAZY LOAD    │                    │
       │                    │                    │                    │
       │ GET /v1/rooms/{id}/routes               │                    │
       │───────────────────▶│  ... (same flow)   │                    │
       │                    │                    │                    │
       │ [User taps "Weather" tab] ← LAZY LOAD  │                    │
       │                    │                    │                    │
       │ GET /v1/weather/route?route_id={rid}    │                    │
       │───────────────────▶│                    │                    │
       │                    │ Delegates to Python AI/Weather via gRPC │
       │                    │                    │                    │
```

---

## Timeline Summary

| Phase | Minggu | Deliverables |
|-------|--------|-------------|
| **Phase 1**: Foundation | 1-3 | Docker, DB migrations, Redis caching layer, project scaffolding |
| **Phase 2**: Go Backend Core | 3-6 | REST API (21 domains), pagination, lazy loading, middleware |
| **Phase 3**: Python Services | 5-8 | AI service, analytics service, data pipeline, gRPC integration |
| **Phase 4**: Go Realtime | 6-9 | WebSocket, NATS, background workers, smart detection |
| **Phase 5**: Flutter App | 8-14 | Full app with lazy loading, pagination, offline, caching |
| **Phase 6**: Next.js Admin | 12-15 | Admin dashboard with pagination, lazy loading |
| **Phase 7**: Voice + Testing + DevOps | 14-18 | LiveKit PTT, testing (all layers), CI/CD, monitoring |

**Total estimated: ~18 minggu (4.5 bulan)**

---

## Verification Plan

### Automated Tests

```bash
# Go Backend
cd go-road-backend && go test ./... -cover -race

# Python AI Service
cd py-ai-service && pytest tests/ -v --cov=app --cov-report=html

# Python Analytics Service
cd py-analytics-service && pytest tests/ -v --cov=app --cov-report=html

# Python Data Pipeline
cd py-data-pipeline && pytest tests/ -v --cov=app --cov-report=html

# Flutter
cd go_road_app && flutter test --coverage

# Next.js Admin
cd go-road-admin && npm test -- --coverage
cd go-road-admin && npx playwright test

# Integration tests (docker-compose up)
docker compose -f docker-compose.test.yml up --abort-on-container-exit
```

### Manual Verification
- Load testing dengan k6: 10K concurrent WebSocket connections
- Verify Redis caching: cache hit rate > 80% for read-heavy endpoints
- Verify pagination: correct cursor encoding/decoding, no duplicate/missing items
- Verify lazy loading: network waterfall shows sequential tab loads, not all-at-once
- Verify offline mode: disconnect network, verify cached data displays, queue operations, resync on reconnect
- Verify Python services: gRPC health check, AI streaming response, analytics computation accuracy
- Battery test: 4-hour simulated touring session, measure battery impact
