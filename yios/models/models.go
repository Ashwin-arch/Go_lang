package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents an engineer or AI developer on the Yios platform
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         string         `gorm:"default:'developer'" json:"role"` // 'developer' or 'admin'
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// AIFramework defines the targeted AI inference runtime environment
type AIFramework string

const (
	FrameworkFastAPI AIFramework = "FASTAPI"
	FrameworkOllama  AIFramework = "OLLAMA"
	FrameworkVLLM    AIFramework = "VLLM"
)

// SourceType defines whether the deployment originates from a Dockerfile or GitHub repository
type SourceType string

const (
	SourceDockerfile SourceType = "DOCKERFILE"
	SourceGitHubRepo SourceType = "GITHUB_REPO"
)

// DeploymentStatus represents the current state of a Kubernetes deployment
type DeploymentStatus string

const (
	StatusBuilding DeploymentStatus = "BUILDING"
	StatusDeployed DeploymentStatus = "DEPLOYED"
	StatusScaling  DeploymentStatus = "SCALING"
	StatusFailed   DeploymentStatus = "FAILED"
)

// Deployment represents an AI model service deployed on Kubernetes
type Deployment struct {
	ID             uint             `gorm:"primaryKey" json:"id"`
	UserID         uint             `gorm:"not null;index" json:"userId"`
	User           User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name           string           `gorm:"uniqueIndex;not null" json:"name"`
	Framework      AIFramework      `gorm:"not null" json:"framework"`
	SourceType     SourceType       `gorm:"not null" json:"sourceType"`
	SourceURL      string           `json:"sourceUrl"` // GitHub URL or Dockerfile snippet
	ImageTag       string           `gorm:"not null" json:"imageTag"`
	Replicas       int              `gorm:"default:1" json:"replicas"`
	Status         DeploymentStatus `gorm:"default:'DEPLOYED'" json:"status"`
	PublicEndpoint string           `json:"publicEndpoint"`
	CpuRequest     string           `gorm:"default:'500m'" json:"cpuRequest"`
	MemoryRequest  string           `gorm:"default:'1Gi'" json:"memoryRequest"`
	K8sNamespace   string           `gorm:"default:'default'" json:"k8sNamespace"`
	ActiveRevision int              `gorm:"default:1" json:"activeRevision"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
}

// Revision tracks deployment history and rolling update versions
type Revision struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	DeploymentID   uint      `gorm:"not null;index" json:"deploymentId"`
	RevisionNumber int       `gorm:"not null" json:"revisionNumber"`
	ImageTag       string    `gorm:"not null" json:"imageTag"`
	CommitSHA      string    `json:"commitSha"`
	Replicas       int       `gorm:"not null" json:"replicas"`
	ConfigYaml     string    `json:"configYaml"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"createdAt"`
}

// PodStatus represents real-time Kubernetes Pod metrics
type PodStatus struct {
	Name          string  `json:"name"`
	Status        string  `json:"status"` // e.g. "Running", "Pending", "ContainerCreating"
	RestartCount  int     `json:"restartCount"`
	CpuUsagePct   float64 `json:"cpuUsagePct"`
	MemoryUsageMB float64 `json:"memoryUsageMb"`
	NodeName      string  `json:"nodeName"`
	Age           string  `json:"age"`
}

// DeploymentMetrics aggregates real-time monitoring statistics for a deployment
type DeploymentMetrics struct {
	DeploymentID  uint        `json:"deploymentId"`
	Name          string      `json:"name"`
	ReplicasCount int         `json:"replicasCount"`
	CpuUsagePct   float64     `json:"cpuUsagePct"`
	MemoryUsageMB float64     `json:"memoryUsageMb"`
	RequestsPerSec int        `json:"requestsPerSec"`
	Pods          []PodStatus `json:"pods"`
}
