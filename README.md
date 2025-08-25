# E-commerce Platform

The **E-commerce Platform** is a robust, scalable, and secure integrated system designed to provide a comprehensive online shopping experience. Built using a monolithic architecture, all core functionalities—such as user management, product catalog, shopping cart, order processing, and merchant interactions—are consolidated within a single application. This approach simplifies development,  and deployment, ensuring a consistent and efficient delivery of e-commerce services within a unified environment.


## 🛠️ Technologies Used
- 🚀 **gRPC** — Provides high-performance, strongly-typed APIs.
- 📡 **Kafka** — Used to publish balance-related events (e.g., after card creation).
- 📈 **Prometheus** — Collects metrics like request count and latency for each RPC method.
- 🛰️ **OpenTelemetry (OTel)** — Enables distributed tracing for observability.
- 🦫 **Go (Golang)** — Implementation language.
- 🌐 **Echo** — HTTP framework for Go.
- 🪵 **Zap Logger** — Structured logging for debugging and operations.
- 📦 **Sqlc** — SQL code generator for Go.
- 🧳 **Goose** — Database migration tool.
- 🐳 **Docker** — Containerization tool.
- 🧱 **Docker Compose** — Simplifies containerization for development and production environments.
- 🐘 **PostgreSQL** — Relational database for storing user data.
- 📃 **Swago** — API documentation generator.
- 🧭 **Zookeeper** — Distributed configuration management.
- 🔀 **Nginx** — Reverse proxy for HTTP traffic.
- 🔍 **Jaeger** — Distributed tracing for observability.
- 📊 **Grafana** — Monitoring and visualization tool.
- 🧪 **Postman** — API client for testing and debugging endpoints.
- ☸️ **Kubernetes** — Container orchestration platform for deployment, scaling, and management.
- 🧰 **Redis** — In-memory key-value store used for caching and fast data access.
- 📥 **Loki** — Log aggregation system for collecting and querying logs.
- 📤 **Promtail** — Log shipping agent that sends logs to Loki.
- 🔧 **OTel Collector** — Vendor-agnostic collector for receiving, processing, and exporting telemetry data (metrics, traces, logs).
- 🖥️ **Node Exporter** — Exposes system-level (host) metrics such as CPU, memory, disk, and network stats for Prometheus.




## Architecture Ecommerce Platform


### Docker

<img src="./images/architecture_ecommerce_docker.png" alt="docker-architecture">

### Kubernetes

<img src="./images/architecture_ecommerce_kubernetes.png" alt="kubernetes-architecture">


### Screenshoot

### Sql
<img src="./images/ecommerce.png" alt="ecommerce "/>


### Loki
<img src="./images/loki.png" alt="loki" />


### Jaeger
<img src="./images/jaeger.png" alt="jaeger" />


### Prometheus 

#### Alert

<img src="./images/prometheus-alert.png" alt="prometheus" />


<img src="./images/prometheus.png" />



### Grafana Prometheus

<img src="./images/grafana-promethues.png" alt="grafana-prometheus" />



### Node Exporter

<img src="./images/node-exporter.png" alt="node-exporter" />
