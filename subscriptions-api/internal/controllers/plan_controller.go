package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gym-management/subscriptions-api/internal/domain/dtos"
	"github.com/yourusername/gym-management/subscriptions-api/internal/services"
)

// PlanController - Controlador HTTP para planes
type PlanController struct {
	planService *services.PlanService // DI
}

// NewPlanController - Constructor con DI
func NewPlanController(planService *services.PlanService) *PlanController {
	return &PlanController{
		planService: planService,
	}
}

// CreatePlan - POST /plans
func (c *PlanController) CreatePlan(ctx *gin.Context) {
	var req dtos.CreatePlanRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := c.planService.CreatePlan(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, plan)
}

// GetPlan - GET /plans/:id
func (c *PlanController) GetPlan(ctx *gin.Context) {
	id := ctx.Param("id")

	plan, err := c.planService.GetPlanByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, plan)
}

// ListPlans - GET /plans
func (c *PlanController) ListPlans(ctx *gin.Context) {
	var query dtos.ListPlansQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plans, err := c.planService.ListPlans(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, plans)
}
