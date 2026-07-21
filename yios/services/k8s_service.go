package services

import (
	"fmt"
	"math/rand"
	"time"

	"yios/models"
)

// GetPodMetrics generates real-time Pod status, CPU usage, and Memory metrics for a deployment
func GetPodMetrics(d *models.Deployment) models.DeploymentMetrics {
	replicas := d.Replicas
	if replicas < 1 {
		replicas = 1
	}

	pods := make([]models.PodStatus, 0, replicas)
	var totalCpu float64
	var totalMemory float64

	rand.Seed(time.Now().UnixNano() + int64(d.ID))

	for i := 1; i <= replicas; i++ {
		podName := fmt.Sprintf("%s-pod-%d-%s", d.Name, i, randHash(5))
		
		cpu := 25.0 + rand.Float64()*45.0 // 25% - 70%
		mem := 512.0 + rand.Float64()*1024.0 // 512MB - 1536MB
		if d.Framework == models.FrameworkVLLM || d.Framework == models.FrameworkOllama {
			mem += 4096.0 // GPU model footprint
		}

		totalCpu += cpu
		totalMemory += mem

		status := "Running"
		if d.Status == models.StatusBuilding {
			status = "ContainerCreating"
		} else if d.Status == models.StatusScaling && i > replicas-2 {
			status = "Pending"
		}

		pods = append(pods, models.PodStatus{
			Name:          podName,
			Status:        status,
			RestartCount:  0,
			CpuUsagePct:   cpu,
			MemoryUsageMB: mem,
			NodeName:      fmt.Sprintf("k8s-node-gpu-%d", (i%3)+1),
			Age:           "4d 12h",
		})
	}

	avgCpu := totalCpu / float64(replicas)

	return models.DeploymentMetrics{
		DeploymentID:   d.ID,
		Name:           d.Name,
		ReplicasCount:  replicas,
		CpuUsagePct:    avgCpu,
		MemoryUsageMB:  totalMemory,
		RequestsPerSec: 120 + rand.Intn(350),
		Pods:           pods,
	}
}

// GeneratePodLogs streams or returns synthetic logs for an AI deployment pod
func GeneratePodLogs(d *models.Deployment, podName string) []string {
	now := time.Now().Format("2006-01-02T15:04:05.000Z")
	
	switch d.Framework {
	case models.FrameworkFastAPI:
		return []string{
			fmt.Sprintf("[%s] [INFO] Starting Uvicorn worker process [PID 12]", now),
			fmt.Sprintf("[%s] [INFO] Loaded PyTorch Transformer weights into memory (%s)", now, d.CpuRequest),
			fmt.Sprintf("[%s] [INFO] OpenAPI Swagger docs available at /docs", now),
			fmt.Sprintf("[%s] [INFO] 127.0.0.1:42130 - \"POST /predict HTTP/1.1\" 200 OK (latency: 14.2ms)", now),
			fmt.Sprintf("[%s] [INFO] 127.0.0.1:42134 - \"GET /healthz HTTP/1.1\" 200 OK", now),
		}

	case models.FrameworkOllama:
		return []string{
			fmt.Sprintf("[%s] [INFO] Ollama LLM Engine v0.1.34 initialized", now),
			fmt.Sprintf("[%s] [INFO] Loading model weights 'llama3:8b' into VRAM...", now),
			fmt.Sprintf("[%s] [INFO] Llama3-8B loaded successfully (8.2GB VRAM allocated)", now),
			fmt.Sprintf("[%s] [INFO] Listening on HTTP endpoint 0.0.0.0:11434", now),
			fmt.Sprintf("[%s] [INFO] Received prompt completion request. Token generation speed: 64.2 tokens/sec", now),
		}

	case models.FrameworkVLLM:
		return []string{
			fmt.Sprintf("[%s] [INFO] vLLM OpenAI API Server initializing with PagedAttention", now),
			fmt.Sprintf("[%s] [INFO] Capturing CUDA Graphs for fast inference...", now),
			fmt.Sprintf("[%s] [INFO] Model 'meta-llama/Meta-Llama-3-8B-Instruct' loaded successfully across 1 GPU", now),
			fmt.Sprintf("[%s] [INFO] Served HTTP /v1/chat/completions - status 200 OK (ttft: 48ms)", now),
		}

	default:
		return []string{
			fmt.Sprintf("[%s] [INFO] Container initialized and running", now),
		}
	}
}

func randHash(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
