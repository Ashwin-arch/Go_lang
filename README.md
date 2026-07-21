# Go_lang

A repository for Go language projects and applications.

---

## 📁 Projects Overview

1. **[🧮 1. Calculator (GUI)](#-1-calculator-gui)**: A desktop GUI calculator built using Go and [Fyne v2](https://fyne.io/).
2. **[🌐 2. REST API](#-2-rest-api)**: A lightweight, thread-safe RESTful API built in Go using standard library `net/http` and Go 1.22+ routing features.
3. **[📊 3. GraphQL API](#-3-graphql-api)**: A thread-safe GraphQL API built in Go with nested relational schema, CRUD mutations, unit tests, and an embedded **GraphiQL IDE**.
4. **[🔒 4. Authentication System](#-4-authentication-system)**: A secure authentication system in Go featuring **Bcrypt password hashing**, **HTTP-Only session cookies**, and a modern **Glassmorphic web interface** (HTML/CSS/JS).
5. **[🛒 5. E-commerce REST API (Gin)](#-5-e-commerce-rest-api-gin)**: A full-featured e-commerce backend built with the **Gin Framework**, **GORM**, and **SQLite**, featuring atomic database transactions for checkout.
6. **[🏦 6. Banking REST API (Echo)](#-6-banking-rest-api-echo)**: A financial banking backend built with the **Echo Framework**, **GORM**, and **SQLite**, featuring ACID transactional transfers with pessimistic row locking (`FOR UPDATE`) and audit logging.
7. **[🤖 7. Yios - AI Kubernetes Platform](#-7-yios---ai-kubernetes-platform)**: An AI deployment and orchestration control plane platform built with **Go (Gin)**, **GORM**, and a **Dual-Mode Kubernetes Controller**, featuring 1-click elastic scaling (1 → 100 replicas), AI model runtime provisioning (**FastAPI**, **Ollama**, **vLLM**), rolling updates, 1-click rollbacks, and an interactive glassmorphic dashboard.

---

## 🧮 1. Calculator (GUI)

A cross-platform desktop calculator application built using Go and Fyne v2.

---

## 🌐 2. REST API

A thread-safe RESTful API built in Go using **zero external dependencies**, leveraging Go's standard library `net/http` package and Go 1.22+ routing features.

---

## 📊 3. GraphQL API

A lightweight, thread-safe GraphQL API built in Go using `github.com/graphql-go/graphql` with zero code-generation dependencies. Features full CRUD operations for Books and Authors, nested relationship resolvers, and an embedded interactive **GraphiQL IDE**.

---

## 🔒 4. Authentication System

A secure authentication system in Go featuring **Bcrypt password hashing**, **HTTP-Only session cookies**, and a modern **Glassmorphic web interface** (HTML/CSS/JS).

---

## 🛒 5. E-commerce REST API (Gin)

A production-grade E-commerce RESTful API built using the **Gin Web Framework** (`github.com/gin-gonic/gin`), **GORM ORM** (`gorm.io/gorm`), and **SQLite** (`gorm.io/driver/sqlite`).

---

## 🏦 6. Banking REST API (Echo)

A financial Banking RESTful API built using the **Echo Web Framework** (`github.com/labstack/echo/v4`), **GORM ORM** (`gorm.io/gorm`), and **SQLite** (`gorm.io/driver/sqlite`).

---

## 🤖 7. Yios - AI Kubernetes Platform

An AI Kubernetes deployment and orchestration platform designed to automate container builds, AI inference engine provisioning (**FastAPI**, **Ollama**, **vLLM**), Kubernetes manifest synthesis, elastic auto-scaling (1 → 100 replicas), zero-downtime rolling updates, 1-click rollbacks, and real-time Pod metrics/log streaming.

### 🏛️ Architecture Overview

```
                  +-----------------------------------+
                  |      Yios React/TS Web UI         |
                  |  (Live Dashboard, Scaling Sliders) |
                  +-----------------+-----------------+
                                    |
                                    v
+-----------------------------------+-----------------------------------+
|               Go Gin Control Plane REST Engine                       |
|  (Auth, AI Runtime Synthesizer, Dual-Mode K8s Controller, Proxy Router) |
+------------------+---------------------------------+------------------+
                   |                                 |
                   v                                 v
+------------------+------------------+   +----------+-------------------+
|    GORM SQLite / PostgreSQL DB      |   |  Kubernetes Cluster / Simulator  |
|  (Users, Deployments, Revisions)    |   |  (Deployments, Services, Ingress)|
+-------------------------------------+   +----------------------------------+
```

### ✨ Key Features
- **🔐 User Authentication**: Developer accounts, login, and token session management.
- **📦 Source Payload Ingestion**: Upload custom Dockerfiles or link GitHub Repositories.
- **🤖 AI Runtimes**: Support for **FastAPI**, **Ollama** (`Llama3:8b`), and **vLLM** (PagedAttention GPU engine).
- **⚙️ K8s Resource Auto-Synthesis**: Automatically generates `Deployment`, `Service`, `Ingress`, and `HPA` specs.
- **📊 Real-time Monitoring**: Pod status tracking (`Running`, `Pending`), CPU & VRAM/RAM metrics, and live log stream viewer.
- **📈 1-Click Elastic Scaling**: Scale replicas dynamically from **1 to 100** using UI sliders or API endpoint.
- **🔄 Rolling Updates & Rollbacks**: Revisions history (`v1`, `v2`, `v3`) with 1-click rollbacks.
- **🌐 Public Ingress Route Generation**: Public endpoint URLs (`http://model-name.yios.internal`) with inference proxy router.

### 🚀 How to Run & Test
```bash
cd yios
go run main.go
```
Server runs at `http://localhost:8083`.

#### Run Unit Tests
```bash
cd yios
go test -v ./...
```

### 📡 cURL Examples
```bash
# 1. Deploy New vLLM AI Model Service
curl -i -X POST http://localhost:8083/api/v1/deployments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "vllm-deepseek-7b",
    "framework": "VLLM",
    "sourceType": "GITHUB_REPO",
    "sourceUrl": "https://github.com/vllm-project/vllm.git",
    "imageTag": "vllm/vllm-openai:latest",
    "replicas": 4,
    "memoryRequest": "16Gi"
  }'

# 2. Scale Replicas to 50
curl -i -X POST http://localhost:8083/api/v1/deployments/1/scale \
  -H "Content-Type: application/json" \
  -d '{"replicas": 50}'

# 3. Trigger Rolling Update
curl -i -X POST http://localhost:8083/api/v1/deployments/1/update \
  -H "Content-Type: application/json" \
  -d '{
    "imageTag": "yios-registry/fastapi-sentiment:v2.0.0-gpu",
    "commitSha": "f9e8d7c6",
    "description": "Upgraded PyTorch runtime & enabled TensorRT acceleration"
  }'

# 4. Execute 1-Click Rollback to Revision 1
curl -i -X POST http://localhost:8083/api/v1/deployments/1/rollback \
  -H "Content-Type: application/json" \
  -d '{"revisionNumber": 1}'
```
