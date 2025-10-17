package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Inscripcion struct {
	IdUsuario        uint      `gorm:"column:id_usuario;primaryKey"`
	IdActividad      uint      `gorm:"column:id_actividad;primaryKey"`
	FechaInscripcion time.Time `gorm:"column:fecha_inscripcion;type:timestamp;default:CURRENT_TIMESTAMP;not null"`
	IsActiva         bool      `gorm:"column:is_activa;default:true;not null"`

	Usuario   Usuario   `gorm:"foreignKey:IdUsuario;constraint:OnDelete:CASCADE"`
	Actividad Actividad `gorm:"foreignKey:IdActividad;constraint:OnDelete:CASCADE"`
}

type Inscripciones []Inscripcion

// verificación antes de hacer INSERT: validamos que haya cupo en la actividad
func (ins *Inscripcion) BeforeCreate(tx *gorm.DB) (err error) {
	var lugares int64

	err = tx.Model(&ActividadVista{}).
		Select("lugares").
		Where("id_actividad = ?", ins.IdActividad).
		First(&lugares).Error
	if err != nil {
		return err
	}

	if lugares <= 0 {
		return fmt.Errorf("No se puede inscribir, el cupo de la actividad ha sido alcanzado.")
	}

	return nil
}

// verificación antes de hacer UPDATE: antes de activar una ins. verificamos que haya cupo
func (ins *Inscripcion) BeforeUpdate(tx *gorm.DB) (err error) {
	if ins.IsActiva {
		var lugares int64

		err = tx.Model(&ActividadVista{}).
			Select("lugares").
			Where("id_actividad = ?", ins.IdActividad).
			First(&lugares).Error
		if err != nil {
			return err
		}

		if lugares <= 0 {
			return fmt.Errorf("No se puede inscribir, el cupo de la actividad ha sido alcanzado.")
		}
	}

	return nil
}
