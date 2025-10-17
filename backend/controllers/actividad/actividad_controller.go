package actividad

import (
	"net/http"
	"proyecto-integrador/dto"
	"proyecto-integrador/services"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func GetActividadesByParams(ctx *gin.Context) {
	actividades, err := services.ActividadService.GetActividadesByParams(map[string]any{
		"id":        ctx.Query("id"),
		"titulo":    ctx.Query("titulo"),
		"horario":   ctx.Query("horario"),
		"categoria": ctx.Query("categoria")},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar actividades"})
	}

	ctx.JSON(http.StatusOK, actividades)
}

func GetAllActividades(ctx *gin.Context) {
	actividades, err := services.ActividadService.GetAllActividades()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar actividades"})
		return
	}

	ctx.JSON(http.StatusOK, actividades)
}

func GetActividadById(ctx *gin.Context) {
	id_actividad, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El id debe ser un numero"})
		return
	}

	actividad, err := services.ActividadService.GetActividadByID(id_actividad)
	if err != nil {
		log.Error("Error al buscar actividad:", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "La actividad no existe"})
		return
	}

	ctx.JSON(http.StatusOK, actividad)
}

func CreateActividad(ctx *gin.Context) {
	var actividadDTO dto.ActividadDTO
	if err := ctx.BindJSON(&actividadDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto"})
		return
	}

	err := services.ActividadService.CreateActividad(actividadDTO)
	if err != nil {
		log.Error("Error al crear actividad:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la actividad"})
		return
	}

	ctx.JSON(http.StatusCreated, actividadDTO)
}

func UpdateActividad(ctx *gin.Context) {
	idActividad, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El id debe ser un número"})
		return
	}

	var actividadDTO dto.ActividadDTO
	if err := ctx.BindJSON(&actividadDTO); err != nil {
		log.Error("Error al parsear JSON:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto"})
		return
	}

	// Validar que el ID en la URL coincida con el ID en el body
	if actividadDTO.Id != 0 && actividadDTO.Id != uint(idActividad) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El ID en la URL no coincide con el ID en el body"})
		return
	}

	actividadDTO.Id = uint(idActividad)
	err = services.ActividadService.UpdateActividad(actividadDTO)
	if err != nil {
		errString := err.Error()

		if strings.Contains(errString, "inscripciones activas que superan el nuevo límite") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, actividadDTO)
}

func DeleteActividad(ctx *gin.Context) {
	idActividad, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El id debe ser un número"})
		return
	}

	err = services.ActividadService.DeleteActividad(uint(idActividad))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
