package services

import (
	"activities-api/internal/domain"
	"activities-api/internal/repository"
	"context"
	"fmt"
)

// InscripcionesService define la interfaz del servicio de inscripciones
type InscripcionesService interface {
	ListByUser(ctx context.Context, usuarioID uint) ([]domain.InscripcionResponse, error)
	Create(ctx context.Context, usuarioID, actividadID uint) (domain.InscripcionResponse, error)
	Deactivate(ctx context.Context, usuarioID, actividadID uint) error
}

// InscripcionesServiceImpl implementa InscripcionesService
// Migrado de backend/services/inscripcion_service.go con dependency injection
type InscripcionesServiceImpl struct {
	inscripcionesRepo repository.InscripcionesRepository
	actividadesRepo   repository.ActividadesRepository
}

// NewInscripcionesService crea una nueva instancia del servicio
func NewInscripcionesService(inscripcionesRepo repository.InscripcionesRepository, actividadesRepo repository.ActividadesRepository) *InscripcionesServiceImpl {
	return &InscripcionesServiceImpl{
		inscripcionesRepo: inscripcionesRepo,
		actividadesRepo:   actividadesRepo,
	}
}

// ListByUser obtiene todas las inscripciones de un usuario
// Migrado de backend/services/inscripcion_service.go:24
func (s *InscripcionesServiceImpl) ListByUser(ctx context.Context, usuarioID uint) ([]domain.InscripcionResponse, error) {
	inscripciones, err := s.inscripcionesRepo.ListByUser(ctx, usuarioID)
	if err != nil {
		return nil, fmt.Errorf("error listing inscripciones: %w", err)
	}

	// Convertir a Response DTO
	responses := make([]domain.InscripcionResponse, len(inscripciones))
	for i, insc := range inscripciones {
		responses[i] = insc.ToResponse()
	}

	return responses, nil
}

// Create inscribe a un usuario en una actividad
// Migrado de backend/services/inscripcion_service.go:44
func (s *InscripcionesServiceImpl) Create(ctx context.Context, usuarioID, actividadID uint) (domain.InscripcionResponse, error) {
	// TODO: Validación 1 - Validar que el usuario existe (HTTP call a users-api)
	// if err := s.validateUserExists(ctx, usuarioID); err != nil {
	//     return domain.InscripcionResponse{}, fmt.Errorf("usuario inválido: %w", err)
	// }

	// TODO: Validación 2 - Validar que tiene suscripción activa (HTTP call a subscriptions-api)
	// activeSub, err := s.getActiveSubscription(ctx, usuarioID)
	// if err != nil {
	//     return domain.InscripcionResponse{}, fmt.Errorf("no tiene suscripción activa: %w", err)
	// }

	// TODO: Validación 3 - Validar que el plan cubra la actividad
	// actividad, err := s.actividadesRepo.GetByID(ctx, actividadID)
	// if err != nil {
	//     return domain.InscripcionResponse{}, fmt.Errorf("actividad no encontrada: %w", err)
	// }
	// if actividad.RequierePlanPremium && activeSub.Plan.TipoAcceso != "completo" {
	//     return domain.InscripcionResponse{}, fmt.Errorf("esta actividad requiere plan premium")
	// }

	// Validar que la actividad existe
	_, err := s.actividadesRepo.GetByID(ctx, actividadID)
	if err != nil {
		return domain.InscripcionResponse{}, fmt.Errorf("actividad no encontrada: %w", err)
	}

	// Crear inscripción
	inscripcion := domain.Inscripcion{
		UsuarioID:   usuarioID,
		ActividadID: actividadID,
		IsActiva:    true,
		// TODO: SuscripcionID: activeSub.ID (cuando subscriptions-api esté listo)
	}

	createdInscripcion, err := s.inscripcionesRepo.Create(ctx, inscripcion)
	if err != nil {
		return domain.InscripcionResponse{}, fmt.Errorf("error creating inscripcion: %w", err)
	}

	// TODO: Publicar evento a RabbitMQ
	// if err := s.publisher.Publish(ctx, "inscription.created", createdInscripcion.ID); err != nil {
	//     log.Printf("Error publishing event: %v", err)
	// }

	return createdInscripcion.ToResponse(), nil
}

// Deactivate desinscribe a un usuario de una actividad
// Migrado de backend/services/inscripcion_service.go:48
func (s *InscripcionesServiceImpl) Deactivate(ctx context.Context, usuarioID, actividadID uint) error {
	if err := s.inscripcionesRepo.Deactivate(ctx, usuarioID, actividadID); err != nil {
		return fmt.Errorf("error deactivating inscripcion: %w", err)
	}

	// TODO: Publicar evento a RabbitMQ
	// if err := s.publisher.Publish(ctx, "inscription.deleted", inscripcionID); err != nil {
	//     log.Printf("Error publishing event: %v", err)
	// }

	return nil
}

// TODO: Los compañeros deben implementar estas funciones:
//
// func (s *InscripcionesServiceImpl) validateUserExists(ctx context.Context, userID uint) error {
//     resp, err := http.Get(fmt.Sprintf("http://users-api:8080/users/%d", userID))
//     if err != nil {
//         return fmt.Errorf("error validating user: %w", err)
//     }
//     defer resp.Body.Close()
//
//     if resp.StatusCode == 404 {
//         return errors.New("user not found")
//     }
//     if resp.StatusCode != 200 {
//         return errors.New("error validating user")
//     }
//
//     return nil
// }
//
// func (s *InscripcionesServiceImpl) getActiveSubscription(ctx context.Context, userID uint) (Subscription, error) {
//     resp, err := http.Get(fmt.Sprintf("http://subscriptions-api:8081/subscriptions/active/%d", userID))
//     if err != nil {
//         return Subscription{}, fmt.Errorf("error getting active subscription: %w", err)
//     }
//     defer resp.Body.Close()
//
//     if resp.StatusCode == 404 {
//         return Subscription{}, errors.New("no active subscription found")
//     }
//     if resp.StatusCode != 200 {
//         return Subscription{}, errors.New("error getting active subscription")
//     }
//
//     var subscription Subscription
//     if err := json.NewDecoder(resp.Body).Decode(&subscription); err != nil {
//         return Subscription{}, fmt.Errorf("error decoding subscription: %w", err)
//     }
//
//     return subscription, nil
// }
