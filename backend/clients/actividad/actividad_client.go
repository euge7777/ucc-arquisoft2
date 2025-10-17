package actividad

import (
	"fmt"
	"proyecto-integrador/db"
	"proyecto-integrador/model"

	log "github.com/sirupsen/logrus"
)

func GetActividadesByParams(params map[string]any) model.ActividadesVista {
	var actividades model.ActividadesVista
	query := db.GetInstance().Model(&model.ActividadVista{})

	if params["id"] != "" {
		query = query.Where("id_actividad = ?", params["id"])
	}
	if params["titulo"] != "" {
		query = query.Where("titulo LIKE ?", fmt.Sprintf("%%%s%%", params["titulo"]))
	}
	if params["horario"] != "" {
		query = query.Where("TIME(?) BETWEEN TIME(horario_inicio) AND TIME(horario_final)", params["horario"])
	}
	if params["categoria"] != "" {
		query = query.Where("categoria LIKE ?", fmt.Sprintf("%%%s%%", params["categoria"]))
	}

	var err error = query.Find(&actividades).Error
	if err != nil {
		log.Error("Error al buscar actividades: ", err)
	}

	log.Debug("Actividades: ", actividades)

	return actividades
}

func GetActividadById(id int) model.ActividadVista {
	var actividad model.ActividadVista
	db.GetInstance().Where("id_actividad = ?", id).First(&actividad)

	log.Debug("Actividad: ", actividad)

	return actividad
}

func GetAllActividades() model.ActividadesVista {
	var actividades model.ActividadesVista
	db.GetInstance().Find(&actividades)

	log.Debug("Actividades: ", actividades)

	return actividades
}

func CreateActividad(actividad model.Actividad) error {
	result := db.GetInstance().Create(&actividad)
	if result.Error != nil {
		log.Error("Error al crear actividad:", result.Error)
		return result.Error
	}
	return nil
}

func UpdateActividad(actividad model.Actividad) error {
	result := db.GetInstance().Model(&model.Actividad{Id: actividad.Id, Cupo: actividad.Cupo}).
		Updates(&actividad)

	if result.Error != nil {
		log.Error("Error al actualizar actividad:", result.Error)
		return result.Error
	}

	return nil
}

func DeleteActividad(id uint) error {
	result := db.GetInstance().Delete(&model.Actividad{}, id)
	if result.Error != nil {
		log.Error("Error al eliminar actividad:", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Error("Error al eliminar actividad: no se encontro actividad con ID ", id)
		return fmt.Errorf("no se encontró la actividad con ID %d", id)
	}
	return nil
}
