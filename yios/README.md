# 🤖 Yios - AI Kubernetes Control Plane Platform

**Yios** is an AI Kubernetes deployment and orchestration platform designed to automate container builds, AI inference engine provisioning (**FastAPI**, **Ollama**, **vLLM**), Kubernetes manifest synthesis, elastic auto-scaling (1 → 100 replicas), zero-downtime rolling updates, 1-click rollbacks, and real-time Pod metrics/log streaming.

---

## 🏛️ Architecture Overview

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

---

## ✨ Key Platform Features

- 🔐 **User Authentication**: Secure developer signup, login, and token session management (`/api/v1/auth/register`, `/api/v1/auth/login`).
- 📦 **Payload Ingestion**: Deploy AI models via custom **Dockerfile** snippets or **GitHub Repositories**.
- 🤖 **AI Model Runtime Support**:
  - **FastAPI**: Custom Python REST API inference server.
  - **Ollama**: Local LLM server deployment (`Llama3:8b`, `Mistral`, `Phi-3`).
  - **vLLM**: GPU-accelerated PagedAttention high-throughput LLM engine.
- ⚙️ **K8s Manifest Auto-Synthesis**: Automatically generates and applies `Deployment`, `Service` (ClusterIP), `Ingress` (NGINX), and `HorizontalPodAutoscaler` (HPA) manifests.
- 📊 **Real-time Pod Monitoring**: Displays live Pod status (`Running`, `Pending`, `ContainerCreating`), CPU usage %, VRAM/RAM memory footprints, and live container log streaming.
- 📈 **1-Click Elastic Scaling**: Scale replicas dynamically from **1 to 100** using UI sliders or API endpoint.
- 🔄 **Rolling Updates & 1-Click Rollbacks**: Track revision history (`v1`, `v2`, `v3`) with commit SHAs and trigger instant rollbacks.
- 🌐 **Public Ingress Route Generation**: Generates public URLs (e.g. `http://model-name.yios.internal`) with an integrated proxy inference router.

---

## ⚙️ How to Run & Test

### 1. Run Control Plane Server
```bash
cd yios
go run main.go
```
The server will start listening at `http://localhost:8083/`.

### 2. Run Automated Test Suite
```bash
cd yios
go test -v ./...
```

---

## 📡 API Endpoints & cURL Testing

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Register developer account |
| `POST` | `/api/v1/auth/login` | Authenticate developer credentials |
| `GET` | `/api/v1/deployments` | List all AI model deployments |
| `POST` | `/api/v1/deployments` | Deploy new AI model service |
| `POST` | `/api/v1/deployments/:id/scale` | 1-Click scale replicas (1 -> 100) |
| `POST` | `/api/v1/deployments/:id/update` | Zero-downtime rolling update |
| `POST` | `/api/v1/deployments/:id/rollback` | 1-Click rollback to previous revision |
| `GET` | `/api/v1/deployments/:id/revisions` | View deployment revision history |
| `GET` | `/api/v1/deployments/:id/metrics` | Fetch real-time Pod status & CPU/Mem metrics |
| `GET` | `/api/v1/deployments/:id/logs` | Stream live pod logs |
| `GET` | `/api/v1/deployments/:id/manifest` | Inspect synthesized Kubernetes YAML |
| `POST` | `/api/v1/deployments/:id/proxy` | Test public inference endpoint |

### cURL Examples

#### 1. Deploy New vLLM AI Model Service
```bash
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
```

#### 2. Scale Replicas to 50
```bash
curl -i -X POST http://localhost:8083/api/v1/deployments/1/scale \
  -H "Content-Type: application/json" \
  -d '{"replicas": 50}'
```

#### 3. Trigger Zero-Downtime Rolling Update
```bash
curl -i -X POST http://localhost:8083/api/v1/deployments/1/update \
  -H "Content-Type: application/json" \
  -d '{
    "imageTag": "yios-registry/fastapi-sentiment:v2.0.0-gpu",
    "commitSha": "f9e8d7c6",
    "description": "Upgraded PyTorch runtime & enabled TensorRT acceleration"
  }'
```

#### 4. Execute 1-Click Rollback to Revision 1
```bash
curl -i -X POST http://localhost:8083/api/v1/deployments/1/rollback \
  -H "Content-Type: application/json" \
  -d '{"revisionNumber": 1}'
```
