package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Actividad struct {
	Id            uint      `gorm:"column:id_actividad;primaryKey;autoIncrement"`
	Titulo        string    `gorm:"type:varchar(50);not null"`
	Descripcion   string    `gorm:"type:varchar(255)"`
	Cupo          uint      `gorm:"type:int;not null"`
	Dia           string    `gorm:"type:enum('Lunes','Martes','Miercoles','Jueves','Viernes','Sabado','Domingo');not null"`
	HorarioInicio time.Time `gorm:"column:horario_inicio;type:time;not null"`
	HorarioFinal  time.Time `gorm:"column:horario_final;type:time;not null"`
	FotoUrl       string    `gorm:"column:foto_url;type:varchar(511);not null"`
	Instructor    string    `gorm:"type:varchar(50);not null"`
	Categoria     string    `gorm:"type:varchar(40);not null"`

	Inscripciones Inscripciones `gorm:"foreignKey:IdActividad;constraint:OnDelete:CASCADE"`
}

type Actividades []Actividad

// verificación antes de hacer UPDATE: validamos que el cupo sea >= a la cantidad de inscriptos
func (ac *Actividad) BeforeUpdate(tx *gorm.DB) (err error) {
	var insc_activas int64

	err = tx.Model(&Inscripcion{}).
		Where("id_actividad = ? AND is_activa = ?", ac.Id, true).
		Count(&insc_activas).Error
	if err != nil {
		return err
	}

	if ac.Cupo < uint(insc_activas) {
		return fmt.Errorf("No se puede cambiar el cupo, hay inscripciones activas que superan el nuevo límite.")
	}

	return nil
}

type ActividadVista struct {
	Id            uint      `gorm:"column:id_actividad;primaryKey;autoIncrement"`
	Titulo        string    `gorm:"type:varchar(50);not null"`
	Descripcion   string    `gorm:"type:varchar(255)"`
	Cupo          uint      `gorm:"type:int;not null"`
	Dia           string    `gorm:"type:enum('Lunes','Martes','Miercoles','Jueves','Viernes','Sabado','Domingo');not null"`
	HorarioInicio time.Time `gorm:"column:horario_inicio;type:time;not null"`
	HorarioFinal  time.Time `gorm:"column:horario_final;type:time;not null"`
	FotoUrl       string    `gorm:"column:foto_url;type:varchar(511);not null"`
	Instructor    string    `gorm:"type:varchar(50);not null"`
	Categoria     string    `gorm:"type:varchar(40);not null"`
	Lugares       uint      `gorm:"column:lugares"` // Campo calculado de la vista
}

type ActividadesVista []ActividadVista

func (ActividadVista) TableName() string {
	return "actividads_lugares"
}
