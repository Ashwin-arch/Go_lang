package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"yios/database"
	"yios/models"
	"yios/routes"
)

func setupTestApp() http.Handler {
	database.InitDB(":memory:")
	return routes.SetupRouter()
}

func TestYiosAIKubernetesPlatform(t *testing.T) {
	app := setupTestApp()

	// 1. Test Auth: User Registration & Login
	t.Run("Register & Login Developer Account", func(t *testing.T) {
		regBody, _ := json.Marshal(map[string]string{
			"name":     "Test AI Dev",
			"email":    "aidev@yios.ai",
			"password": "password123",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(regBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected 201 Created, got %d (body: %s)", w.Code, w.Body.String())
		}

		loginBody, _ := json.Marshal(map[string]string{
			"email":    "aidev@yios.ai",
			"password": "password123",
		})
		reqLogin := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
		reqLogin.Header.Set("Content-Type", "application/json")
		wLogin := httptest.NewRecorder()
		app.ServeHTTP(wLogin, reqLogin)

		if wLogin.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK for login, got %d", wLogin.Code)
		}
	})

	// 2. Test Listing AI Model Deployments (Seeded: FastAPI sentiment analyzer & Ollama Llama-3)
	t.Run("Get Active Deployments", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/deployments", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}
	})

	// 3. Test Deploying New vLLM AI Model (GPU PagedAttention Engine)
	t.Run("Deploy vLLM AI Model Service", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"name":          "vllm-deepseek-r1",
			"framework":     models.FrameworkVLLM,
			"sourceType":    models.SourceGitHubRepo,
			"sourceUrl":     "https://github.com/vllm-project/vllm.git",
			"imageTag":      "vllm/vllm-openai:v0.4.0",
			"replicas":      4,
			"memoryRequest": "16Gi",
		})
		req := httptest.NewRequest("POST", "/api/v1/deployments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected 201 Created for vLLM deployment, got %d (body: %s)", w.Code, w.Body.String())
		}
	})

	// 4. Test 1-Click Elastic Scaling (Scale Deployment 1 to 50 Replicas)
	t.Run("Elastic Scaling 1 -> 50 Replicas", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"replicas": 50,
		})
		req := httptest.NewRequest("POST", "/api/v1/deployments/1/scale", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for replica scaling, got %d", w.Code)
		}
	})

	// 5. Test Zero-Downtime Rolling Update & Revision History Creation
	t.Run("Rolling Update & Revision History", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"imageTag":    "yios-registry/fastapi-sentiment:v2.0.0-gpu",
			"commitSha":   "987654321abcdef",
			"description": "Upgraded PyTorch runtime & enabled TensorRT acceleration",
		})
		req := httptest.NewRequest("POST", "/api/v1/deployments/1/update", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for rolling update, got %d", w.Code)
		}

		// Verify Revision list
		reqRevs := httptest.NewRequest("GET", "/api/v1/deployments/1/revisions", nil)
		wRevs := httptest.NewRecorder()
		app.ServeHTTP(wRevs, reqRevs)

		if wRevs.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for revisions, got %d", wRevs.Code)
		}
	})

	// 6. Test 1-Click Rollback to Revision 1
	t.Run("1-Click Rollback to Previous Revision", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"revisionNumber": 1,
		})
		req := httptest.NewRequest("POST", "/api/v1/deployments/1/rollback", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for rollback, got %d", w.Code)
		}
	})

	// 7. Test Real-time Pod Metrics & K8s Manifest Synthesis
	t.Run("Pod Metrics & Manifest YAML Synthesis", func(t *testing.T) {
		reqMetrics := httptest.NewRequest("GET", "/api/v1/deployments/1/metrics", nil)
		wMetrics := httptest.NewRecorder()
		app.ServeHTTP(wMetrics, reqMetrics)

		if wMetrics.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK for pod metrics, got %d", wMetrics.Code)
		}

		reqManifest := httptest.NewRequest("GET", "/api/v1/deployments/1/manifest", nil)
		wManifest := httptest.NewRecorder()
		app.ServeHTTP(wManifest, reqManifest)

		if wManifest.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK for manifest synthesis, got %d", wManifest.Code)
		}
	})
}
