package inscripcion

import (
	"proyecto-integrador/db"
	"proyecto-integrador/model"

	"errors"

	log "github.com/sirupsen/logrus"
)

func GetAllInscripciones(id_usuario uint) (model.Inscripciones, error) {
	var inscripciones model.Inscripciones
	query := db.GetInstance().Model(&model.Inscripcion{})

	var err error
	if err = query.Where("id_usuario = ?", id_usuario).Find(&inscripciones).Error; err != nil {
		log.Error("Error al buscar inscripciones: ", err)
		return nil, err
	}

	log.Debug("Inscripciones: ", inscripciones)

	return inscripciones, nil
}

func AltaDeInscripcion(id_usuario, id_actividad uint) error {
	insc := model.Inscripcion{
		IdUsuario:   id_usuario,
		IdActividad: id_actividad,
	}

	result := db.GetInstance().Where(&insc).FirstOrCreate(&insc)
	err := result.Error
	if err != nil {
		return err
	}

	// si el registro ya existe, actualizamos is_activa = 1
	if result.RowsAffected == 0 {
		if insc.IsActiva {
			return errors.New("El usuario ya esta inscripto")
		}

		insc.IsActiva = true
		return db.GetInstance().Model(&insc).
			Update("is_activa", true).Error
	}

	return nil
}

func BajaDeInscripcion(id_usuario, id_actividad uint) error {
	// hacemos una consulta tipo UPDATE a la base de datos
	return db.GetInstance().Model(&model.Inscripcion{
		IdUsuario:   id_usuario,
		IdActividad: id_actividad}).
		Update("is_activa", false).
		Where("id_usuario = ? AND id_actividad = ?", id_usuario, id_actividad).Error
}
