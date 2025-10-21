package services

import (
	"context"
	"time"

	"github.com/yourusername/gym-management/subscriptions-api/internal/domain/dtos"
	"github.com/yourusername/gym-management/subscriptions-api/internal/domain/entities"
	"github.com/yourusername/gym-management/subscriptions-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PlanService - Servicio de lógica de negocio para planes
type PlanService struct {
	planRepo repository.PlanRepository // Inyección de Dependencias (Interface)
}

// NewPlanService - Constructor con DI
func NewPlanService(planRepo repository.PlanRepository) *PlanService {
	return &PlanService{
		planRepo: planRepo,
	}
}

// CreatePlan - Crea un nuevo plan
func (s *PlanService) CreatePlan(ctx context.Context, req dtos.CreatePlanRequest) (*dtos.PlanResponse, error) {
	// Mapear DTO a entidad
	plan := &entities.Plan{
		ID:                    primitive.NewObjectID(),
		Nombre:                req.Nombre,
		Descripcion:           req.Descripcion,
		PrecioMensual:         req.PrecioMensual,
		TipoAcceso:            req.TipoAcceso,
		DuracionDias:          req.DuracionDias,
		Activo:                req.Activo,
		ActividadesPermitidas: req.ActividadesPermitidas,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	// Guardar en repositorio
	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}

	// Mapear entidad a DTO de respuesta
	return s.mapPlanToResponse(plan), nil
}

// GetPlanByID - Obtiene un plan por ID
func (s *PlanService) GetPlanByID(ctx context.Context, id string) (*dtos.PlanResponse, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	plan, err := s.planRepo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	return s.mapPlanToResponse(plan), nil
}

// ListPlans - Lista planes con filtros
func (s *PlanService) ListPlans(ctx context.Context, query dtos.ListPlansQuery) (*dtos.PaginatedPlansResponse, error) {
	// Construir filtros
	filters := make(map[string]interface{})
	if query.Activo != nil {
		filters["activo"] = *query.Activo
	}

	// Obtener planes
	plansList, err := s.planRepo.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Obtener total
	total, err := s.planRepo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Mapear a DTOs
	var plans []dtos.PlanResponse
	for _, plan := range plansList {
		plans = append(plans, *s.mapPlanToResponse(plan))
	}

	// Calcular paginación
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dtos.PaginatedPlansResponse{
		Plans:      plans,
		Total:      int(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// mapPlanToResponse - Helper para mapear entidad a DTO
func (s *PlanService) mapPlanToResponse(plan *entities.Plan) *dtos.PlanResponse {
	return &dtos.PlanResponse{
		ID:                    plan.ID.Hex(),
		Nombre:                plan.Nombre,
		Descripcion:           plan.Descripcion,
		PrecioMensual:         plan.PrecioMensual,
		TipoAcceso:            plan.TipoAcceso,
		DuracionDias:          plan.DuracionDias,
		Activo:                plan.Activo,
		ActividadesPermitidas: plan.ActividadesPermitidas,
		CreatedAt:             plan.CreatedAt,
		UpdatedAt:             plan.UpdatedAt,
	}
}
