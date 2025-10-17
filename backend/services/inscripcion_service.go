package services

import (
	"proyecto-integrador/clients/inscripcion"
	"proyecto-integrador/dto"
)

type inscripcionService struct{}

type IInscripcionService interface {
	GetAllInscripciones(id_usuario uint) (dto.InscripcionesDTO, error)
	InscribirUsuario(id_usuario, id_actividad uint) error
	DesinscribirUsuario(id_usuario, id_actividad uint) error
}

var (
	InscripcionService IInscripcionService
)

func init() {
	InscripcionService = &inscripcionService{}
}

func (is *inscripcionService) GetAllInscripciones(id_usuario uint) (dto.InscripcionesDTO, error) {
	inscripciones, err := inscripcion.GetAllInscripciones(id_usuario)
	if err != nil {
		return nil, err
	}

	var resultado dto.InscripcionesDTO
	for _, v := range inscripciones {
		dto := dto.InscripcionDTO{
			IdUsuario:        v.IdUsuario,
			IdActividad:      v.IdActividad,
			FechaInscripcion: v.FechaInscripcion.GoString(),
			IsActiva:         v.IsActiva,
		}
		resultado = append(resultado, dto)
	}

	return resultado, nil
}

func (is *inscripcionService) InscribirUsuario(id_usuario, id_actividad uint) error {
	return inscripcion.AltaDeInscripcion(id_usuario, id_actividad)
}

func (is *inscripcionService) DesinscribirUsuario(id_usuario, id_actividad uint) error {
	return inscripcion.BajaDeInscripcion(id_usuario, id_actividad)
}
