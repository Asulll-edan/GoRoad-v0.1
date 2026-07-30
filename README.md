# GoRoad

aplikasi buat komunitas touring motor/hiking ke gunung bareng besties. awalnya iseng bikin pas lagi bosen dikantor waktu pkl.

## apa aja fiturnya

- bikin room touring, ngajak temen, atur formasi berangkat
- chat realtime pas touring biar ga kehilangan rombongan
- tracking lokasi biar tau siapa yang ketinggalan
- rute touring + waypoint + export gpx
- sos & emergency button kalo ada apa apa
- voting buat mutusin tempat makan atau istirahat
- catat bensin & expense bareng bareng
- checklist barang biar ga lupa bawa jas hujan
- feed & social biar bisa pamer touring
- leaderboard & badge buat kompetisi jarak tempuh
- ai chat pake gemini buat saran rute atau packing
- admin dashboard buat mantau semua activity
- ptt voice pake livekit kalo males ngetik

## tech stack

`go` buat backend api, websocket, worker.
`python` buat ai service, analytics, data pipeline (celery).
`flutter` buat mobile app (android & ios).
`next.js` buat admin dashboard.
`postgresql` + `postgis` + `timescaledb` buat database.
`redis` buat cache & message broker.
`nats` buat event driven pub/sub.
`minio` buat nyimpen file upload.
`livekit` buat voice chat.
`osrm` buat route engine.
`docker` buat deployment.

## cara jalanin

paling gampang pake docker compose:

```bash
docker compose up --build -d
```

kalo mau jalanin manual, ada petunjuk dikit di bawah.

### prerequisite

yang harus jalan sebelum ngapa ngapain:
- postgresql 16 + postgis + timescaledb
- redis 7
- nats jetstream
- minio

abis itu tinggal jalanin service yang mau dipake.

### go backend

```bash
cd go-road-backend

# api
go run ./cmd/api/main.go

# websocket
go run ./cmd/ws/main.go

# worker
go run ./cmd/worker/main.go
```

### python services

```bash
cd py-ai-service
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8001
```

```bash
cd py-analytics-service
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8002
```

```bash
cd py-data-pipeline
pip install -r requirements.txt
celery -A app.tasks.celery_app worker --loglevel=info
celery -A app.tasks.celery_app beat --loglevel=info
```

### flutter

```bash
cd go_road_app_flutter
flutter pub get
flutter run
```

### admin dashboard

```bash
cd go-road-admin
npm install
npm run dev
```

## environment

copy `.env.example` ke `.env` terus isi api key yang diperlukan:

| variable | kegunaan |
|----------|----------|
| `JWT_SECRET` | secret buat jwt token |
| `GEMINI_API_KEY` | api key gemini buat ai chat |
| `OPENWEATHER_API_KEY` | api key openweather buat cuaca |
| `LIVEKIT_API_KEY` | api key buat voice chat |
| `LIVEKIT_API_SECRET` | secret buat voice chat |
| `SENTRY_DSN` | dsn sentry buat error tracking |

## struktur folder

```
├── go-road-backend/        # go backend (api, ws, worker)
├── go_road_app_flutter/     # flutter mobile app
├── go-road-admin/           # next.js admin dashboard
├── py-ai-service/           # python ai service (gemini)
├── py-analytics-service/    # python analytics service
├── py-data-pipeline/        # python celery workers
├── py-shared/               # shared python utilities
├── proto/                   # grpc proto definitions
├── migrations/              # database migrations
├── docker/                  # dockerfile & nginx config
├── monitoring/              # prometheus, grafana, loki
├── loadtest/                # k6 load test script
└── docker-compose.yml       # all services orchestration
```

## catatan

proyek ini masih jalan, fitur nambah dikit demi dikit. kalo ada yang mau dibenerin atau ditambah, silahkan kontribusi.

dibikin pake kopi dan rokok tipis tipis.
