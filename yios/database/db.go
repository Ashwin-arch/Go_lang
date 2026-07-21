package database

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"yios/models"
)

var DB *gorm.DB

// InitDB connects to the database and performs auto-migrations
func InitDB(dbPath string) *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to Yios database: %v", err)
	}

	fmt.Println("📦 Yios Database connected successfully.")

	// Auto-migrate schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Deployment{},
		&models.Revision{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	SeedData(DB)
	return DB
}

// SeedData inserts initial demo users, AI model deployments, and revision history
func SeedData(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Name:         "AI Engineer",
		Email:        "developer@yios.ai",
		PasswordHash: string(hash),
		Role:         "developer",
	}
	db.Create(&user)

	// Seed Sample FastAPI Deployment
	d1 := models.Deployment{
		UserID:         user.ID,
		Name:           "fastapi-sentiment-analyzer",
		Framework:      models.FrameworkFastAPI,
		SourceType:     models.SourceGitHubRepo,
		SourceURL:      "https://github.com/yios-ai/fastapi-nlp-demo.git",
		ImageTag:       "yios-registry/fastapi-sentiment:v1.2.0",
		Replicas:       3,
		Status:         models.StatusDeployed,
		PublicEndpoint: "http://fastapi-sentiment.yios.internal",
		CpuRequest:     "500m",
		MemoryRequest:  "1Gi",
		K8sNamespace:   "default",
		ActiveRevision: 2,
	}
	db.Create(&d1)

	db.Create(&models.Revision{
		DeploymentID:   d1.ID,
		RevisionNumber: 1,
		ImageTag:       "yios-registry/fastapi-sentiment:v1.0.0",
		CommitSHA:      "a1b2c3d4",
		Replicas:       1,
		Description:    "Initial release deployment",
	})
	db.Create(&models.Revision{
		DeploymentID:   d1.ID,
		RevisionNumber: 2,
		ImageTag:       "yios-registry/fastapi-sentiment:v1.2.0",
		CommitSHA:      "f9e8d7c6",
		Replicas:       3,
		Description:    "Scaled replicas & added batch inference route",
	})

	// Seed Sample Ollama Model Deployment
	d2 := models.Deployment{
		UserID:         user.ID,
		Name:           "ollama-llama3-8b",
		Framework:      models.FrameworkOllama,
		SourceType:     models.SourceDockerfile,
		SourceURL:      "FROM ollama/ollama:latest\nRUN ollama pull llama3:8b",
		ImageTag:       "ollama/ollama:latest",
		Replicas:       2,
		Status:         models.StatusDeployed,
		PublicEndpoint: "http://ollama-llama3.yios.internal",
		CpuRequest:     "2000m",
		MemoryRequest:  "8Gi",
		K8sNamespace:   "default",
		ActiveRevision: 1,
	}
	db.Create(&d2)

	db.Create(&models.Revision{
		DeploymentID:   d2.ID,
		RevisionNumber: 1,
		ImageTag:       "ollama/ollama:latest",
		CommitSHA:      "e5f4g3h2",
		Replicas:       2,
		Description:    "Initial Ollama Llama-3 deployment",
	})

	fmt.Println("🌱 Initial demo AI deployments & revisions seeded.")
}
