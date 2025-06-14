package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mukund/mediaconvert/internal/auth"
	"github.com/mukund/mediaconvert/internal/models"
	"github.com/mukund/mediaconvert/internal/pipeline"
	"gorm.io/gorm"
)

type PipelineHandler struct {
	db *gorm.DB
}

func NewPipelineHandler(db *gorm.DB) *PipelineHandler {
	return &PipelineHandler{db: db}
}

type CreatePipelineRequest struct {
	Name    string `json:"name" binding:"required"`
	Format  string `json:"format" binding:"required,oneof=yaml json"`
	Content string `json:"content" binding:"required"`
}

type PipelineResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Format  string `json:"format"`
	Content string `json:"content"`
}

// CreatePipeline creates a new pipeline
func (h *PipelineHandler) CreatePipeline(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse and validate pipeline
	var p *pipeline.Pipeline
	var err error
	if req.Format == "yaml" {
		p, err = pipeline.ParseYAML([]byte(req.Content))
	} else {
		p, err = pipeline.ParseJSON([]byte(req.Content))
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline format: " + err.Error()})
		return
	}

	if err := p.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pipeline validation failed: " + err.Error()})
		return
	}

	// Check for duplicate name
	var existing models.Pipeline
	if err := h.db.Where("user_id = ? AND name = ?", userID, req.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Pipeline with this name already exists"})
		return
	}

	// Create pipeline
	pipelineModel := models.Pipeline{
		UserID:  userID,
		Name:    req.Name,
		Format:  models.PipelineFormat(req.Format),
		Content: req.Content,
	}

	if err := h.db.Create(&pipelineModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pipeline"})
		return
	}

	c.JSON(http.StatusCreated, PipelineResponse{
		ID:      pipelineModel.ID,
		Name:    pipelineModel.Name,
		Format:  string(pipelineModel.Format),
		Content: pipelineModel.Content,
	})
}

// ListPipelines returns all pipelines for the user
func (h *PipelineHandler) ListPipelines(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var pipelines []models.Pipeline
	if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&pipelines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pipelines"})
		return
	}

	response := make([]PipelineResponse, len(pipelines))
	for i, p := range pipelines {
		response[i] = PipelineResponse{
			ID:      p.ID,
			Name:    p.Name,
			Format:  string(p.Format),
			Content: p.Content,
		}
	}

	c.JSON(http.StatusOK, gin.H{"pipelines": response})
}

// GetPipeline returns a single pipeline
func (h *PipelineHandler) GetPipeline(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pipelineID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	var p models.Pipeline
	if err := h.db.Where("id = ? AND user_id = ?", pipelineID, userID).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pipeline"})
		}
		return
	}

	c.JSON(http.StatusOK, PipelineResponse{
		ID:      p.ID,
		Name:    p.Name,
		Format:  string(p.Format),
		Content: p.Content,
	})
}

// UpdatePipeline updates an existing pipeline
func (h *PipelineHandler) UpdatePipeline(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pipelineID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	var req CreatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse and validate pipeline
	var p *pipeline.Pipeline
	if req.Format == "yaml" {
		p, err = pipeline.ParseYAML([]byte(req.Content))
	} else {
		p, err = pipeline.ParseJSON([]byte(req.Content))
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline format: " + err.Error()})
		return
	}

	if err := p.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pipeline validation failed: " + err.Error()})
		return
	}

	// Find and update pipeline
	var pipelineModel models.Pipeline
	if err := h.db.Where("id = ? AND user_id = ?", pipelineID, userID).First(&pipelineModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pipeline"})
		}
		return
	}

	pipelineModel.Name = req.Name
	pipelineModel.Format = models.PipelineFormat(req.Format)
	pipelineModel.Content = req.Content

	if err := h.db.Save(&pipelineModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pipeline"})
		return
	}

	c.JSON(http.StatusOK, PipelineResponse{
		ID:      pipelineModel.ID,
		Name:    pipelineModel.Name,
		Format:  string(pipelineModel.Format),
		Content: pipelineModel.Content,
	})
}

// DeletePipeline deletes a pipeline
func (h *PipelineHandler) DeletePipeline(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pipelineID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", pipelineID, userID).Delete(&models.Pipeline{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pipeline"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pipeline deleted successfully"})
}
