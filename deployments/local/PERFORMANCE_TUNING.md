# Performance Tuning & Release Notes (Milestone E)

Dokumen ini merangkum hasil **Milestone E — Performance & Production Readiness**:
baseline load test, panduan tuning, dan release notes. Berlaku untuk stack lokal
(`deployments/local/`) yang berjalan via docker compose.

---

## 1. Baseline (diukur 2026-08-07)

Semua pengukuran diambil dari mesin dev (5.5 GiB RAM, 8 CPU, Docker compose,
stack lengkap 20 service + infra). Alat: `deployments/local/scripts/load-test.sh`
(bash + curl, tanpa dependency tambahan).

| Endpoint | Req | Worker | Throughput | p50 | p95 | p99 | Fail |
|---|---|---|---|---|---|---|---|
| `GET /health` | 100 | 10 | 341 req/s | 23.1 ms | 34.6 ms | 39.7 ms | 0 |

### Catatan baseline
- `/health` hanya memvalidasi gateway (tanpa dependency chain) — ini *upper bound*
  throughput gateway.
- Endpoint dengan chain gRPC (mis. `/api/auth/login` → auth → user/role) akan jauh
  lebih lambat dan lebih peka terhadap pool/koneksi; ukur terpisah bila perlu.
- Untuk benchmark realistis di CI/hardware dedicated, gunakan k6/hey/wrk (tidak
  ditambahkan ke repo agar bebas dependency).

### Cara mengulang
```bash
just up            # atau docker compose -f deployments/local/docker-compose.yml up -d
just smoke         # pastikan 6/6 PASS dulu
bash deployments/local/scripts/load-test.sh http://localhost:5000 200 10
ENDPOINT=/api/product-query bash deployments/local/scripts/load-test.sh
```

---

## 2. Temuan & perbaikan selama validasi (Fase 7–8)

Validasi nyata (minikube → dibatalkan karena hardware, lanjut docker compose)
menemukan beberapa bug yang **tidak terlihat oleh render/config static**:

### 2.1 PostgreSQL connection exhaustion
- Gejala: semua service `exit 2` dengan `FATAL: sorry, too many clients already`.
- Akar: `pkg/database` default `DB_MIN_IDLE_CONNS=50` per service × 20 service
  = 1000 koneksi idle, sedangkan postgres default `max_connections=100`.
- Fix di `deployments/local/docker.env`:
  - `DB_MAX_OPEN_CONNS=10`, `DB_MAX_IDLE_CONNS=2`
- Fix di `deployments/local/docker-compose.yml` (postgres):
  - `command: ["postgres", "-c", "max_connections=400", "-c", "shared_buffers=256MB"]`
- **Tuning**: untuk produksi, sesuaikan `max_connections` dengan jumlah service
  × pool; jangan biarkan `DB_MIN_IDLE_CONNS` default 50 per service.

### 2.2 Migrate service tidak jalan
- Gejala: tabel tidak pernah dibuat; `migrate-ecommerce` exit 0 tanpa migrasi.
- Akar 1: compose tidak meneruskan argumen → `command: ["./migrate", "up"]`.
- Akar 2: `service/migrate/Dockerfile` tidak meng-Copy `migrations/` ke image
  runtime → `COPY --from=builder /app/service/migrate/migrations ./migrations`.

### 2.3 Kafka image `bitnami/*:latest` tidak tersedia
- `bitnami/kafka:latest` dan `bitnami/zookeeper:latest` tidak ada manifest di
  registry (validasi nyata: pull gagal). Kafka diganti ke `apache/kafka:latest`
  dengan env `KAFKA_*` (KRaft), zookeeper dihapus (tidak dibutuhkan KRaft).
- Volume kafka perlu `chown 1000:1000` (image berjalan sebagai `appuser`):
  ```bash
  docker run --rm -v local_kafka_ecommerce_data:/data alpine chown -R 1000:1000 /data
  ```

### 2.4 Apigateway: HTTP server tidak pernah start
- Gejala: `RunClient()` di `service/apigateway/apps/client.go` hanya memanggil
  `NewClient()` — `client.Run()` (yang menjalankan `c.Echo.Start`) tidak pernah
  dipanggil, sehingga port 5000 tidak pernah bind meski log sukses.
- Fix: `RunClient()` kini memulai `client.Run()` di goroutine.

### 2.5 Apigateway: `Cache: nil` → panic di `/api/auth/login`
- Gejala: login panic `nil pointer dereference` di handler/auth handler.go:150
  (`h.cache.GetCachedLogin`).
- Akar: `handler.NewHandler` meng-pass `Cache: nil` dengan TODO.
- Fix: `Cache: auth_cache.NewMencache(deps.Cache)` di `service/apigateway/handler/handler.go`.

### 2.6 Apigateway: GRPC service address selalu default (salah)
- Gejala: `/api/product-query` 503; gateway dial `product:50064` padahal service
  listen `50058`.
- Akar: `loadServiceAddresses()` memakai `v.SetEnvPrefix("grpc")` + replacer,
  tetapi `dotenv.Viper()` sudah memanggil `AutomaticEnv()` tanpa prefix lebih
  dulu → `SetEnvPrefix` tidak berpengaruh → env `GRPC_*_ADDR` tidak pernah
  terbaca → selalu fallback ke default hardcoded yang salah.
- Fix: baca `GRPC_<NAME>_ADDR` secara eksplisit via `v.GetString("GRPC_"+name+"_ADDR")`
  dengan fallback default yang sudah dikoreksi (product:50058, user:50053, dst).

---

## 3. Panduan tuning

### 3.1 PostgreSQL
| Parameter | Dev | Produksi (estimasi) |
|---|---|---|
| `max_connections` | 400 | `service_count × max_pool × 2` |
| `shared_buffers` | 256 MB | 25% RAM |
| `DB_MAX_OPEN_CONNS` (app) | 10 | 20–50 |
| `DB_MIN_IDLE_CONNS` (app) | 2 | 5–10 |

### 3.2 Redis
- Per-service DB (`REDIS_DB_*`) sudah dipisah — jangan digabung di produksi
  tanpa Redis Cluster.
- Healthcheck pakai `redis-cli -a dragon_knight ping` (requirepass).

### 3.3 Kafka
- KRaft (tanpa zookeeper), `apache/kafka:latest`.
- Volume data harus dimiliki `1000:1000`.
- Partisi/topic default OK untuk dev; produksi perlu `replication.factor=3`.

### 3.4 gRPC gateway
- Keepalive client sudah dikonfigurasi (`PermitWithoutStream: true`).
- Trace context dipropagasi lintas gRPC (OTel).
- Rate limit `/api/auth/*` = 10 rps/burst (brute-force guard).

---

## 4. Release notes (ringkasan Fase 7–8)

### Fase 7 — Kubernetes Readiness
- `kustomization.yaml` (baru), render 116 resource, DNS-1123 bersih.
- Fix DNS underscore → dash (7 service).
- Probes 34/34 deployment; Ingress; 5 NetworkPolicy; migrate Job + TTL.
- Email-hpa.yaml diperbaiki (salinan category-hpa).
- **Residual**: validasi server-side (apply/rollout/smoke di cluster nyata)
  tidak dieksekusi penuh — minikube dibatalkan karena hardware (5.5 GiB RAM).
  Sebagai gantinya, validasi nyata dilakukan di docker compose (lihat §2).

### Fase 8 — Backlog Error Handling + Validasi Nyata + Milestone E
- Backlog: `InvalidAccessToken` → `ErrUnauthorized`; round-trip test 429/503;
  test duplicate register (auth service).
- Compose stack: down stack asing, fix kafka/migrate/pool/env, smoke 6/6 PASS.
- Milestone E: secret audit (`secrets.yaml` di-untrack + template example),
  baseline load test (341 req/s), tuning guide ini.

---

## 5. Status & residual jujur
- [x] Smoke test 6/6 PASS (health, ready, register, login, me, product-query).
- [x] Baseline load test `/health` (341 req/s, p99 39.7 ms).
- [x] Secret audit: tidak ada creds hardcoded di kode; `secrets.yaml` k8s di-untrack.
- [ ] Load test chain gRPC (login) — belum diukur, perlu endpoint + setup khusus.
- [ ] Kubernetes server-side validation — butuh cluster ≥ 8 GiB RAM.
- [ ] Audit CVE image (trivy/grype) — belum dijalankan.
