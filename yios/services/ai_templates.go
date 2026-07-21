package services

import (
	"fmt"

	"yios/models"
)

// GenerateDockerfile generates a optimized Dockerfile for the selected AI Framework
func GenerateDockerfile(framework models.AIFramework, customDockerfile string) string {
	if customDockerfile != "" {
		return customDockerfile
	}

	switch framework {
	case models.FrameworkFastAPI:
		return `# Optimized FastAPI AI Inference Engine Dockerfile
FROM python:3.10-slim

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt uvicorn fastapi torch transformers

COPY . .
EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "4"]`

	case models.FrameworkOllama:
		return `# Ollama LLM Container Environment
FROM ollama/ollama:latest

EXPOSE 11434
ENV OLLAMA_HOST=0.0.0.0:11434

ENTRYPOINT ["/bin/ollama"]
CMD ["serve"]`

	case models.FrameworkVLLM:
		return `# High-Performance GPU-Accelerated vLLM Engine
FROM vllm/vllm-openai:latest

EXPOSE 8000
ENTRYPOINT ["python3", "-m", "vllm.entrypoints.openai.api_server"]
CMD ["--model", "meta-llama/Meta-Llama-3-8B-Instruct", "--port", "8000", "--tensor-parallel-size", "1"]`

	default:
		return "# Default App Dockerfile\nFROM alpine:latest\nCMD [\"echo\", \"Running AI Service\"]"
	}
}

// GenerateK8sManifests generates K8s Deployment, Service, Ingress, and HPA YAML specifications
func GenerateK8sManifests(d *models.Deployment) string {
	port := 8000
	if d.Framework == models.FrameworkOllama {
		port = 11434
	}

	return fmt.Sprintf(`---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
  namespace: %s
  labels:
    app: %s
    platform: yios
spec:
  replicas: %d
  selector:
    matchLabels:
      app: %s
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%%
      maxUnavailable: 25%%
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - name: ai-model
        image: %s
        ports:
        - containerPort: %d
        resources:
          requests:
            cpu: "%s"
            memory: "%s"
          limits:
            cpu: "4000m"
            memory: "16Gi"
---
apiVersion: v1
kind: Service
metadata:
  name: %s-svc
  namespace: %s
spec:
  type: ClusterIP
  selector:
    app: %s
  ports:
  - port: 80
    targetPort: %d
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: %s-ingress
  namespace: %s
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  rules:
  - host: %s.yios.internal
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: %s-svc
            port:
              number: 80
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: %s-hpa
  namespace: %s
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: %s
  minReplicas: 1
  maxReplicas: 100
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 80
`, d.Name, d.K8sNamespace, d.Name, d.Replicas, d.Name, d.Name, d.ImageTag, port, d.CpuRequest, d.MemoryRequest,
		d.Name, d.K8sNamespace, d.Name, port,
		d.Name, d.K8sNamespace, d.Name, d.Name,
		d.Name, d.K8sNamespace, d.Name)
}
