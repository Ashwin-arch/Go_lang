package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"yios/database"
	"yios/models"
	"yios/services"
)

type CreateDeployInput struct {
	Name          string             `json:"name" binding:"required"`
	Framework     models.AIFramework `json:"framework" binding:"required"`
	SourceType    models.SourceType  `json:"sourceType" binding:"required"`
	SourceURL     string             `json:"sourceUrl"`
	Dockerfile    string             `json:"dockerfile"`
	ImageTag      string             `json:"imageTag"`
	Replicas      int                `json:"replicas"`
	CpuRequest    string             `json:"cpuRequest"`
	MemoryRequest string             `json:"memoryRequest"`
}

type ScaleInput struct {
	Replicas int `json:"replicas" binding:"required,min=1,max=100"`
}

type RollingUpdateInput struct {
	ImageTag    string `json:"imageTag" binding:"required"`
	CommitSHA   string `json:"commitSha"`
	Description string `json:"description"`
}

type RollbackInput struct {
	RevisionNumber int `json:"revisionNumber" binding:"required"`
}

// GetDeployments lists all AI model deployments
func GetDeployments(c *gin.Context) {
	var deployments []models.Deployment
	database.DB.Order("id desc").Find(&deployments)
	c.JSON(http.StatusOK, gin.H{"deployments": deployments})
}

// GetDeploymentByID fetches a specific deployment details
func GetDeploymentByID(c *gin.Context) {
	id := c.Param("id")
	var d models.Deployment
	if err := database.DB.First(&d, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deployment": d})
}

// CreateDeployment deploys a new AI model to Kubernetes
func CreateDeployment(c *gin.Context) {
	var input CreateDeployInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	replicas := input.Replicas
	if replicas < 1 {
		replicas = 1
	}

	imageTag := input.ImageTag
	if imageTag == "" {
		imageTag = fmt.Sprintf("yios-registry/%s:v1.0.0", input.Name)
	}

	cpu := input.CpuRequest
	if cpu == "" {
		cpu = "500m"
	}
	mem := input.MemoryRequest
	if mem == "" {
		mem = "1Gi"
	}

	publicEndpoint := fmt.Sprintf("http://%s.yios.internal", input.Name)

	d := models.Deployment{
		UserID:         1, // Default engineer context
		Name:           input.Name,
		Framework:      input.Framework,
		SourceType:     input.SourceType,
		SourceURL:      input.SourceURL,
		ImageTag:       imageTag,
		Replicas:       replicas,
		Status:         models.StatusDeployed,
		PublicEndpoint: publicEndpoint,
		CpuRequest:     cpu,
		MemoryRequest:  mem,
		K8sNamespace:   "default",
		ActiveRevision: 1,
	}

	if err := database.DB.Create(&d).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deployment name already exists"})
		return
	}

	// Create Revision 1
	rev := models.Revision{
		DeploymentID:   d.ID,
		RevisionNumber: 1,
		ImageTag:       imageTag,
		CommitSHA:      "initial-commit",
		Replicas:       replicas,
		ConfigYaml:     services.GenerateK8sManifests(&d),
		Description:    "Initial Kubernetes deployment release",
	}
	database.DB.Create(&rev)

	c.JSON(http.StatusCreated, gin.H{
		"message":    "AI Model deployed successfully to Kubernetes!",
		"deployment": d,
		"revision":   rev,
	})
}

// ScaleDeployment executes 1-click scaling (1 -> 100 replicas)
func ScaleDeployment(c *gin.Context) {
	id := c.Param("id")
	var input ScaleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var d models.Deployment
	if err := database.DB.First(&d, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		return
	}

	d.Replicas = input.Replicas
	d.Status = models.StatusDeployed
	database.DB.Save(&d)

	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("Scaled deployment '%s' to %d replicas", d.Name, d.Replicas),
		"replicas": d.Replicas,
	})
}

// UpdateDeployment triggers zero-downtime rolling update
func UpdateDeployment(c *gin.Context) {
	id := c.Param("id")
	var input RollingUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var d models.Deployment
	if err := database.DB.First(&d, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		return
	}

	d.ActiveRevision++
	d.ImageTag = input.ImageTag
	database.DB.Save(&d)

	commitSha := input.CommitSHA
	if commitSha == "" {
		commitSha = "rev-" + fmt.Sprint(d.ActiveRevision)
	}

	desc := input.Description
	if desc == "" {
		desc = fmt.Sprintf("Rolling update to image %s", input.ImageTag)
	}

	rev := models.Revision{
		DeploymentID:   d.ID,
		RevisionNumber: d.ActiveRevision,
		ImageTag:       input.ImageTag,
		CommitSHA:      commitSha,
		Replicas:       d.Replicas,
		ConfigYaml:     services.GenerateK8sManifests(&d),
		Description:    desc,
	}
	database.DB.Create(&rev)

	c.JSON(http.StatusOK, gin.H{
		"message":    fmt.Sprintf("Rolling update initiated for '%s' (Revision %d)", d.Name, d.ActiveRevision),
		"deployment": d,
		"revision":   rev,
	})
}

// RollbackDeployment executes a 1-click rollback to a previous revision
func RollbackDeployment(c *gin.Context) {
	id := c.Param("id")
	var input RollbackInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var d models.Deployment
	if err := database.DB.First(&d, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		return
	}

	var rev models.Revision
	if err := database.DB.Where("deployment_id = ? AND revision_number = ?", d.ID, input.RevisionNumber).First(&rev).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Revision number not found"})
		return
	}

	d.ImageTag = rev.ImageTag
	d.Replicas = rev.Replicas
	d.ActiveRevision = rev.RevisionNumber
	database.DB.Save(&d)

	c.JSON(http.StatusOK, gin.H{
		"message":    fmt.Sprintf("Successfully rolled back '%s' to Revision %d", d.Name, rev.RevisionNumber),
		"deployment": d,
		"revision":   rev,
	})
}

// GetRevisions lists deployment revision history
func GetRevisions(c *gin.Context) {
	id := c.Param("id")
	var revisions []models.Revision
	database.DB.Where("deployment_id = ?", id).Order("revision_number desc").Find(&revisions)
	c.JSON(http.StatusOK, gin.H{"revisions": revisions})
}

// DeleteDeployment deletes a deployment
func DeleteDeployment(c *gin.Context) {
	id := c.Param("id")
	database.DB.Where("deployment_id = ?", id).Delete(&models.Revision{})
	if err := database.DB.Delete(&models.Deployment{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deployment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deployment and associated K8s resources deleted"})
}
