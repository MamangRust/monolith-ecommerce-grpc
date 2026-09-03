# Observability Runbooks

Runbook operasional untuk alert yang didefinisikan di `observability/rules/`.
Setiap alert punya prosedur: **triage → diagnosis → recovery → prevent**.

Daftar isi:

- [1. ServiceDown (crash loop)](#1-servicedown-crash-loop)
- [2. PostgresPoolExhaustion](#2-postgrespoolexhaustion)
- [3. RedisUnavailable / RedisHighMissRate](#3-redisunavailable--redishighmissrate)
- [4. GrpcUnavailableHigh / GrpcTimeoutHigh](#4-grpcunavailablehigh--grpctimeouthigh)
- [5. Http5xxSpike](#5-http5xxspike)
- [6. KafkaPublishErrorsHigh / KafkaConsumeErrorsHigh](#6-kafkapublisherrrorshigh--kafkaconsumeerrrorshigh)
- [7. Error rate / latency alert per service](#7-error-rate--latency-alert-per-service)

---

## 1. ServiceDown (crash loop)

**Expr:** `up == 0` selama 2 menit — service tidak bisa di-scrape Prometheus.

### Triage
1. Apakah ini deploy baru? (`git log`, image tag). Rollback jika iya.
2. Berapa lama down? `for: 2m` berarti sudah 2+ menit.

### Diagnosis
- `docker ps` / `kubectl get pods` — status container (restart count naik = crash loop).
- Log service: `docker logs <svc> --tail 200` — cari panic, fatal error, `OOMKilled`.
- Resource check: memory/CPU limit (runtime metrics `go_memstats_*` jika sempat naik).

### Recovery
- Restart service. Jika crash loop, cek startup dependency (DB/Redis/Kafka tidak bisa dijangkau).
- Jika OOM: naikkan limit atau fix memory leak.

### Prevent
- Healthcheck di compose/k8s, resource limit realistis, `ENABLE_REFLECTION` off di prod.

---

## 2. PostgresPoolExhaustion

**Expr:** `postgresql_pool_acquired_conns >= postgresql_pool_max_conns`.

### Triage
- Apakah ini setelah traffic spike atau setelah deploy?
- Berapa service yang kena? (pool exhaustion sering menular ke downstream).

### Diagnosis
- Query lambat: `pg_stat_activity` — cari `state = 'idle in transaction'` (leak) atau `long_running`.
- Cek pool config: `DB_MAX_OPEN_CONNS` (default 100), `DB_MIN_IDLE_CONNS` (50).
- Cek apakah transaksi outbox / query besar menggantung.

### Recovery
- Kill idle-in-transaction: `SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE state='idle in transaction' AND now()-state_change > interval '5 minutes';`
- Sementara: naikkan `DB_MAX_OPEN_CONNS`.
- Restart service yang bocor koneksi.

### Prevent
- Pastikan `defer` close/rollback semua transaksi & rows.
- Batasi concurrency (request limiter sudah ada di resilience interceptor).

---

## 3. RedisUnavailable / RedisHighMissRate

**Expr:** `redis_pool_total_conns == 0` (critical) / miss rate > 50% (warning).

### Triage
- Redis down atau hanya satu service yang tidak konek?
- Cache miss tinggi = cache berguna tapi tidak efektif, atau Redis baru restart (cache kosong).

### Diagnosis
- `redis-cli ping` / container health.
- Cek env `REDIS_HOST_<SERVICE>`/`REDIS_PORT_<SERVICE>`.
- Miss rate: apakah TTL terlalu pendek? Apakah eviction aktif (`maxmemory`)?

### Recovery
- Restart Redis atau service yang gagal konek.
- Jika miss tinggi karena Redis restart: biarkan warm-up, atau kurangi TTL supaya tidak thundering herd.

### Prevent
- Redis AOF/persistence sesuai kebutuhan, monitor `redis_pool_hits_total` vs `misses`.

---

## 4. GrpcUnavailableHigh / GrpcTimeoutHigh

**Expr:** rate `grpc_requests_total{grpc_status="Unavailable"|"DeadlineExceeded"}`.

### Triage
- Unavailable = downstream down. Timeout = downstream lambat.
- Service mana yang jadi `method` target? (label `method` = `/svc.Method`)

### Diagnosis
- Cek service target: up? crash loop? resource?
- Latency histogram target: `grpc_request_duration_seconds` di service penerima.
- Cek timeout config: `ContextMiddleware(30s)` di server, keepalive client di gateway.

### Recovery
- Restart service target. Naikkan timeout jika query memang berat (bukan bug).

### Prevent
- Circuit breaker (sudah ada) — pastikan threshold tidak terlalu longgar.
- Alerts per-service latency sudah menandakan lebih awal.

---

## 5. Http5xxSpike

**Expr:** >10% request ke route tertentu 5xx.

### Triage
- Route mana? (label `route`)
- Apakah 5xx konsisten (bug) atau intermittent (downstream flaky)?

### Diagnosis
- Cek upstream gRPC service untuk route itu (Unavailable? Internal?).
- `http_errors_total` + log gateway (`createLoggerMiddleware` mencatat status & error).
- Error mapping: pastikan bukan bug `Internal` yang seharusnya `NotFound/BadRequest`.

### Recovery
- Fix di upstream service; jika downstream down, restore service tsb.

### Prevent
- Error mapping konsisten (`ParseGrpcError` di gateway), test negatif per domain.

---

## 6. KafkaPublishErrorsHigh / KafkaConsumeErrorsHigh

**Expr:** rate `kafka_publish_errors_total` / `kafka_consume_errors_total`.

### Triage
- Topic/group mana? (label `topic` / `group`)
- Apakah broker reachable? `kafka-exporter` up?

### Diagnosis
- Cek `kafka_topic_partition_under_replicated_partition` (alert KafkaUnderReplicatedPartitions).
- Publish: error di service yang memanggil `SendMessage` (outbox relay untuk order/transaction).
- Consume: error log consumer (email service).

### Recovery
- Restart broker/service. Untuk publish: outbox akan retry (transaksional) — pastikan relay jalan.
- Cek offset: consumer lag (`kafka_topic_partition_current_offset` vs high watermark).

### Prevent
- Outbox pattern sudah ada untuk publish reliability. Monitor consumer group lag.

---

## 7. Error rate / latency alert per service

**Expr:** `rate(requests_total{status="error"}) > 0.1` atau p95 latency > 1s.

### Triage
- Service mana? (`requests_total` direkam oleh `NewObservability` per service)
- Command vs query? (beberapa service punya alert terpisah)

### Diagnosis
- Error log dengan `trace_id` (OTLP → Loki) — cari trace yang gagal end-to-end.
- Latency: histogram `request_duration_seconds` — apakah p95 naik karena query DB, cache miss, atau downstream.

### Recovery
- Sesuai akar masalah (DB/Redis/Kafka/downstream). Panik → rollback deploy terakhir.

### Prevent
- Exit criteria Fase 5: trace end-to-end (gateway → gRPC → DB) via traceparent propagation; dependency failure vs business error bisa dibedakan dari status code.
