package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"w3st/usecase"
)

type ProjectController struct {
	projectUsecase usecase.ProjectUsecase
}

func NewProjectController(projectUsecase usecase.ProjectUsecase) *ProjectController {
	return &ProjectController{
		projectUsecase: projectUsecase,
	}
}

func (c *ProjectController) GetAllProjects(ctx *gin.Context) {
	projects, err := c.projectUsecase.GetAllProjects(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var request struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		RateLimit   int    `json:"rate_limit"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// デフォルトのレート制限を設定
	if request.RateLimit == 0 {
		request.RateLimit = 1000
	}

	project, err := c.projectUsecase.CreateProject(ctx.Request.Context(), request.Name, request.Description, request.RateLimit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"project": project})
}

func (c *ProjectController) GetProjectByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := c.projectUsecase.GetProjectByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"project": project})
}