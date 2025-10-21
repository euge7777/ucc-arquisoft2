package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/yourusername/gym-management/search-api/internal/models"
)

// SearchService maneja las operaciones de búsqueda
// NOTA: Esta es una implementación en memoria para desarrollo
// En producción, reemplazar con integración real a Apache Solr
type SearchService struct {
	documents map[string]models.SearchDocument
	mu        sync.RWMutex
}

func NewSearchService() *SearchService {
	return &SearchService{
		documents: make(map[string]models.SearchDocument),
	}
}

// IndexDocument indexa un documento
func (s *SearchService) IndexDocument(doc models.SearchDocument) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.documents[doc.ID] = doc
	return nil
}

// DeleteDocument elimina un documento del índice
func (s *SearchService) DeleteDocument(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.documents, id)
	return nil
}

// Search realiza una búsqueda en los documentos
func (s *SearchService) Search(req models.SearchRequest) (*models.SearchResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Valores por defecto
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	var results []models.SearchDocument

	// Filtrar documentos
	for _, doc := range s.documents {
		if s.matchesSearch(doc, req) {
			results = append(results, doc)
		}
	}

	// Calcular paginación
	totalCount := len(results)
	totalPages := (totalCount + req.PageSize - 1) / req.PageSize

	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	if start >= totalCount {
		results = []models.SearchDocument{}
	} else {
		if end > totalCount {
			end = totalCount
		}
		results = results[start:end]
	}

	return &models.SearchResponse{
		Results:    results,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// matchesSearch verifica si un documento coincide con los criterios de búsqueda
func (s *SearchService) matchesSearch(doc models.SearchDocument, req models.SearchRequest) bool {
	// Filtrar por tipo si se especifica
	if req.Type != "" && doc.Type != req.Type {
		return false
	}

	// Búsqueda por query (búsqueda de texto)
	if req.Query != "" {
		query := strings.ToLower(req.Query)
		text := strings.ToLower(fmt.Sprintf("%s %s %s %s %s",
			doc.Titulo, doc.Descripcion, doc.Categoria,
			doc.Instructor, doc.PlanNombre))

		if !strings.Contains(text, query) {
			return false
		}
	}

	// Aplicar filtros adicionales
	for key, value := range req.Filters {
		switch key {
		case "categoria":
			if doc.Categoria != value {
				return false
			}
		case "dia":
			if doc.Dia != value {
				return false
			}
		case "instructor":
			if doc.Instructor != value {
				return false
			}
		case "sucursal_id":
			if doc.SucursalID != value {
				return false
			}
		case "requiere_premium":
			reqPremium := value == "true"
			if doc.RequierePremium != reqPremium {
				return false
			}
		case "estado":
			if doc.Estado != value {
				return false
			}
		}
	}

	return true
}

// GetDocumentByID obtiene un documento por ID
func (s *SearchService) GetDocumentByID(id string) (*models.SearchDocument, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, exists := s.documents[id]
	if !exists {
		return nil, fmt.Errorf("documento no encontrado")
	}

	return &doc, nil
}

// IndexFromEvent indexa un documento desde un evento de RabbitMQ
func (s *SearchService) IndexFromEvent(event models.RabbitMQEvent) error {
	// Convertir data a SearchDocument
	docBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var doc models.SearchDocument
	if err := json.Unmarshal(docBytes, &doc); err != nil {
		return err
	}

	doc.ID = event.Type + "_" + event.ID
	doc.Type = event.Type

	return s.IndexDocument(doc)
}

// GetStats retorna estadísticas del índice
func (s *SearchService) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := map[string]interface{}{
		"total_documents": len(s.documents),
		"types":           make(map[string]int),
	}

	typeCounts := make(map[string]int)
	for _, doc := range s.documents {
		typeCounts[doc.Type]++
	}
	stats["types"] = typeCounts

	return stats
}
