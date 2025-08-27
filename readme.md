# SATD - Server Agent Threat Detection


SATD is a distributed, concurrent threat detection system written primarily in Go, designed to monitor and analyze activity across multiple systems in real-time. It leverages lightweight agents, a central server, a Node.js-based dashboard, and Elasticsearch for deep analytics.

## 🌐 Overview

SATD consists of:

- **Go-based Server**: A concurrent, gRPC-powered central node that ingests telemetry data from agents and performs preliminary threat detection.
- **Go-based Agent**: A system-level daemon that observes host behavior and transmits metadata to the server.
- **Node.js Dashboard**: A TypeScript-based UI server intended to display summaries and real-time system statuses.
- **Elasticsearch (Dockerized)**: Stores logs and metadata for deep inspection and manual or automated analysis.

This architecture enables scalable, real-time monitoring of networked systems with extensible data pipelines for both security analysis and operational visibility.

---

## ⚙️ Technologies Used

| Component      | Tech Stack                            |
|----------------|----------------------------------------|
| Server         | Go, gRPC, Concurrency (goroutines)     |
| Agent          | Go, gRPC                               |
| Dashboard      | Node.js, TypeScript, Express           |
| Data Backend   | Elasticsearch (via Docker)             |
| Containerization | Docker (for Elasticsearch and possibly others) |
| Transport      | gRPC (TLS optional)                    |
| Observability  | Log messages, Elastic logs, Heartbeats |

---

## 🧠 How It Works

### 1. Agent Behavior
- Collects system or network metadata (e.g., process behavior, network usage).
- Sends data to the server over a secure gRPC channel.
- Periodically emits heartbeat signals for liveness detection.

### 2. Server Behavior
- Accepts concurrent streams of telemetry data using Go’s native concurrency primitives.
- Detects anomalies (e.g., unexpected patterns, missing heartbeats). ( WIP )
- Optionally pushes data to Elasticsearch for indexing.

### 3. Elasticsearch
- Stores logs and metrics for longer-term retention and advanced querying.
- Can be paired with Kibana for visualization.

### 4. Node.js Dashboard (WIP/Optional)
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

## 🗂️ Architecture Diagram

```plaintext
[Agent (Go)] ---> [gRPC] ---> [Server (Go)] ---> [ElasticSearch (Docker)]
                                 |
                                 └--> [Dashboard (Node.js/TS)]

