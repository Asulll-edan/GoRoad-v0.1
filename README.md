<div align="center">
  <br/>
  <h1>GoRoad</h1>
  <p><strong>platform komunitas touring motor</strong></p>
  <br/>
  <p>
    <img src="https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go"/>
    <img src="https://img.shields.io/badge/Python-3.12-3776AB?style=for-the-badge&logo=python&logoColor=white" alt="Python"/>
    <img src="https://img.shields.io/badge/Flutter-3.38-02569B?style=for-the-badge&logo=flutter&logoColor=white" alt="Flutter"/>
    <img src="https://img.shields.io/badge/Dart-3.10-0175C2?style=for-the-badge&logo=dart&logoColor=white" alt="Dart"/>
    <img src="https://img.shields.io/badge/Next.js-15-000000?style=for-the-badge&logo=next.js&logoColor=white" alt="Next.js"/>
  </p>
  <p>
    <img src="https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL"/>
    <img src="https://img.shields.io/badge/Redis-7-DC382D?style=for-the-badge&logo=redis&logoColor=white" alt="Redis"/>
    <img src="https://img.shields.io/badge/NATS-2.10-27AAE1?style=for-the-badge&logo=nats&logoColor=white" alt="NATS"/>
    <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker"/>
    <img src="https://img.shields.io/badge/MinIO-71C53B?style=for-the-badge&logo=minio&logoColor=white" alt="MinIO"/>
    <img src="https://img.shields.io/badge/LiveKit-00C853?style=for-the-badge&logo=livekit&logoColor=white" alt="LiveKit"/>
  </p>
  <p>
    <img src="https://img.shields.io/badge/gRPC-2446E6?style=for-the-badge&logo=grpc&logoColor=white" alt="gRPC"/>
    <img src="https://img.shields.io/badge/Celery-5.4-37814A?style=for-the-badge&logo=celery&logoColor=white" alt="Celery"/>
    <img src="https://img.shields.io/badge/OSRM-5.27-FF6C37?style=for-the-badge&logo=openstreetmap&logoColor=white" alt="OSRM"/>
    <img src="https://img.shields.io/badge/PostGIS-3.4-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostGIS"/>
    <img src="https://img.shields.io/badge/TimescaleDB-2.17-FBB040?style=for-the-badge&logo=timescale&logoColor=white" alt="TimescaleDB"/>
  </p>
  <p>
    <img src="https://img.shields.io/badge/status-active-brightgreen?style=for-the-badge"/>
    <img src="https://img.shields.io/badge/license-MIT-yellow?style=for-the-badge"/>
  </p>
  <br/>
</div>

---

iseng iseng bikin pas lagi bosen, eh malah gede sendiri sampe segini. go road adalah platform buat komunitas touring motor -- biar kalo touring bareng ga ilang ilangan, urusan bensin & expense ke track, dan yang pasti biar makin seru.

## fitur

| fitur | apa fungsinya |
|-------|--------------|
| room touring | bikin room, ngajak temen, atur formasi berangkat |
| chat realtime | ngobrol pas touring -- pake websocket, realtime |
| tracking live | tau posisi rombongan biar ga ada yang ketinggalan |
| rute | bikin rute + waypoint + export gpx + osrm integration |
| sos & emergency | panic button kalo ada apa apa di jalan |
| voting | mutusin tempat makan, istirahat, atau tujuan bareng bareng |
| fuel & expense | catat bensin & pengeluaran, split bill otomatis |
| checklist | biar ga lupa bawa jas hujan, tools, atau dokumen |
| social feed | pamer foto touring, likes & comments |
| leaderboard & badge | kompetisi jarak tempuh & pencapaian |
| ai chat | pake gemini -- nanya rute, packing list, saran safety |
| ptt voice | push to talk pake livekit, biar ga ribet ngetik |
| admin dashboard | pantau semua aktivitas dari satu tempat |

## tech stack

```
CLIENT
  Flutter Mobile (Android + iOS)
  Next.js Admin Dashboard

GATEWAY
  Nginx (ssl, ws proxy, rate limit)

BACKEND
  Go API (REST)
  Go WS (WebSocket)
  Go Worker (background tasks)
  Python AI Service (Gemini gRPC)
  Python Analytics (stats, leaderboard)
  Python Pipeline (Celery workers)

DATA & MESSAGING
  PostgreSQL + PostGIS + TimescaleDB
  Redis 7 (cache, session, rate limit)
  NATS JetStream (event driven)
  MinIO (file storage)
  LiveKit (voice chat)
  OSRM (route engine)
```

<details>
<summary><b>detail versi</b></summary>

| komponen | versi |
|----------|-------|
| Go | 1.23 |
| Python | 3.12 |
| Flutter | 3.38.10 |
| Dart | 3.10.9 |
| Next.js | 15.x |
| PostgreSQL | 16 + PostGIS 3.4 + TimescaleDB 2.17 |
| Redis | 7 (alpine) |
| NATS | 2.10 (JetStream) |
| LiveKit | 1.6 |
| OSRM | 5.27 |
| Celery | 5.4 (Redis broker) |
| gRPC | protobuf (Go to Python) |

</details>

## cara jalanin

### docker (semua service)

paling gampang tinggal jalanin ini:

```bash
cp .env.example .env   # lalu isi api key
docker compose up --build -d
```

### manual (development)

prerequisite: postgres, redis, nats, minio harus jalan duluan.

#### go backend
```bash
cd go-road-backend

# rest API -- port 8080
go run ./cmd/api/main.go

# websocket -- port 8081
go run ./cmd/ws/main.go

# background worker
go run ./cmd/worker/main.go
```

#### python services
```bash
# ai service (gemini) -- port 8001
cd py-ai-service
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8001

# analytics service -- port 8002
cd py-analytics-service
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8002

# data pipeline (celery workers)
cd py-data-pipeline
pip install -r requirements.txt
celery -A app.tasks.celery_app worker --loglevel=info
celery -A app.tasks.celery_app beat --loglevel=info
```

#### flutter
```bash
cd go_road_app_flutter
flutter pub get
flutter run
```

#### admin dashboard
```bash
cd go-road-admin
npm install
npm run dev
```

## environment

copy `.env.example` jadi `.env` terus isi:

| variable | buat apa |
|----------|----------|
| JWT_SECRET | secret jwt token |
| GEMINI_API_KEY | ai chat pake gemini |
| OPENWEATHER_API_KEY | data cuaca realtime |
| LIVEKIT_API_KEY | voice chat |
| LIVEKIT_API_SECRET | secret voice chat |
| SENTRY_DSN | error tracking |

## struktur folder

```
go-road/
  go-road-backend/          go -- api, ws, worker
    cmd/
      api/main.go
      ws/main.go
      worker/main.go
    internal/
      config/
      domain/               21 domain entities
      dto/                  request & response
      handler/              fiber handlers
      middleware/           auth, cors, logger, rate limit
      repository/           postgres + redis
      service/              business logic
      ws/                   websocket hub & handlers
      event/                nats pub/sub
      worker/               background tasks
  go_road_app_flutter/      flutter mobile app
  go-road-admin/            next.js admin dashboard
  py-ai-service/            python -- ai (gemini)
  py-analytics-service/     python -- analytics
  py-data-pipeline/         python -- celery workers
  py-shared/                shared python utils
  proto/                    grpc protobuf definitions
  migrations/               sql migrations (11 files)
  docker/                   dockerfiles + nginx.conf
  monitoring/               prometheus, grafana, loki, sentry
  loadtest/                 k6 load test
  docker-compose.yml
  docker-compose.dev.yml
  .env.example
```

## testing

```bash
# go
cd go-road-backend && go test ./... -cover

# python ai
cd py-ai-service && pytest -v --cov=app

# python analytics
cd py-analytics-service && pytest -v --cov=app

# flutter
cd go_road_app_flutter && flutter test

# admin
cd go-road-admin && npm test
```

## kontribusi

kalo nemu bug atau punya ide, tinggal buka issue atau pull request. siapa tau bisa jadi tandingan waze buat motoran.
---

dibikin pake hati, seneng tipis tipis, lelah tebal tebal.
