# SATD - Server Agent Threat Detection

<img width="1908" height="1166" alt="system_design_diagram" src="https://github.com/user-attachments/assets/6e87ad43-2f4c-4958-afea-173955139746" />


Demo:

https://github.com/user-attachments/assets/1cf696a2-9753-4020-bdc4-b0919a8fbd1a

SATD is a distributed, concurrent threat detection system written primarily in Go, designed to monitor and analyze activity across multiple systems in real-time. It leverages lightweight agents, a central server, a Node.js-based dashboard, and Elasticsearch for deep analytics.


## Feature List

- **Go-based Server**: A concurrent, gRPC-powered central node that ingests telemetry data from agents and performs preliminary threat detection.
- **Go-based Agent**: A system-level daemon that observes host behavior and transmits metadata to the server.
- **Node.js Dashboard**: A TypeScript-based UI server for displaying summaries and real-time system statuses across all agents.
- **Elasticsearch (Dockerized)**: Stores logs and metadata for deep inspection and manual or automated analysis.

---

## Technologies Used

| Component      | Tech Stack                            |
|----------------|----------------------------------------|
| Server         | Go, gRPC, Concurrency (goroutines)     |
| Agent          | Go, gRPC                               |
| Dashboard      | Node.js, TypeScript, Express           |
| Data Backend   | Elasticsearch (via Docker)             |
| Containerization | Docker (for Elasticsearch and possibly others) |
| Transport      | gRPC (TLS encrypted)                    |
| Observability  | Log messages, Elastic logs, Heartbeats |

---

## How It Works

### 1. Agent Behavior
- Collects system or network metadata (e.g., process behavior, network usage).
- Sends data to the server over a secure TLS-encrypted gRPC channel.
- Periodically emits heartbeat signals for liveness detection.

### 2. Server Behavior
- Accepts concurrent streams of telemetry data using Go’s native concurrency primitives.
- Detects anomalies (e.g., unexpected patterns, missing heartbeats). ( WIP )
- Optionally pushes data to Elasticsearch for indexing.

### 3. Elasticsearch
- Stores logs and metrics for longer-term retention and advanced querying.
- Can be paired with Kibana for visualization.

### 4. Node.js / React.js Dashboard
- Provides a web UI for summarizing system status.
- Queries the Go server or Elasticsearch for aggregated data.

---

## Dependencies

```
libpcap dev library (sudo dnf install libpcap-devel)
Go 1.24.5+
Node v22.18.0
npm 10.9.3+
```

## Architecture Diagram

```plaintext
[Agent (Go)] ---> [gRPC] ---> [Server (Go)] ---> [ElasticSearch (Docker)]
                                 |
                                 └--> [Dashboard (Node.js/TS)]
```

## Environment Configuration

### `./server/.server_env` Configuration

```
ELASTIC_API_KEY=YOUR_ELASTIC_KEY
IPDB_API_KEY=YOUR_IPDB_KEY
DASHBOARD_SERVER_AUTH_ADDR=https://localhost:3000/login
DASHBOARD_SERVER_PROT_ADDR=https://localhost:3000/add-dashboard-info
NODEJS_USER=admin
NODEJS_PASS=YOUR_ADMIN_PASSWORD
```

### `./gateway/.env` Configuration

```SECRET_JWT_KEY=YOUR_JWT_KEY
DB_USER=sleepy
DB_HOST=localhost # change as you see fit
DB_NAME=satd
DB_PASSWORD=groovy
DB_PORT=5432
```

remember not to push these env files to the repo :)
