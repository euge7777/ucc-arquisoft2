package services

import (
	"fmt"
	"proyecto-integrador/clients/actividad"
	"proyecto-integrador/dto"
	"proyecto-integrador/model"
	"time"

	log "github.com/sirupsen/logrus"
)

type actividadService struct{}

type IactividadService interface {
	GetAllActividades() (dto.ActividadesDTO, error)
	GetActividadesByParams(params map[string]any) (dto.ActividadesDTO, error)
	GetActividadByID(id int) (dto.ActividadDTO, error)
	DeleteActividad(id uint) error
	CreateActividad(actividadDTO dto.ActividadDTO) error
	UpdateActividad(actividadDTO dto.ActividadDTO) error
}

var (
	ActividadService IactividadService
)

func init() {
	ActividadService = &actividadService{}
}

func validarCamposBasicos(actividadDTO dto.ActividadDTO) error {
	if actividadDTO.Cupo == 0 {
		return fmt.Errorf("el cupo debe ser mayor a 0")
	}

	if actividadDTO.Titulo == "" {
		return fmt.Errorf("el título no puede estar vacío")
	}

	if actividadDTO.Dia == "" {
		return fmt.Errorf("el día no puede estar vacío")
	}

	return nil
}

// Función auxiliar para parsear horas
func parsearHoras(horaInicio, horaFin string) (time.Time, time.Time, error) {
	// Obtener la zona horaria local
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		loc = time.Local // Si no se puede cargar, usar la zona horaria local del sistema
	}

	// Usar una fecha base (2024-01-01) para parsear las horas
	fechaBase := "2024-01-01"
	inicio, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", fechaBase, horaInicio), loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("formato de hora inicio inválido (debe ser HH:MM): %v", err)
	}

	fin, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", fechaBase, horaFin), loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("formato de hora fin inválido (debe ser HH:MM): %v", err)
	}

	// Validar que hora fin sea después de hora inicio
	if fin.Before(inicio) {
		return time.Time{}, time.Time{}, fmt.Errorf("la hora de fin debe ser posterior a la hora de inicio")
	}

	return inicio, fin, nil
}

func convertirADTO(v model.ActividadVista) dto.ActividadDTO {
	return dto.ActividadDTO{
		Id:          v.Id,
		Titulo:      v.Titulo,
		Descripcion: v.Descripcion,
		Cupo:        v.Cupo,
		Dia:         v.Dia,
		HoraInicio:  v.HorarioInicio.Format("15:04"),
		HoraFin:     v.HorarioFinal.Format("15:04"),
		FotoUrl:     v.FotoUrl,
		Instructor:  v.Instructor,
		Categoria:   v.Categoria,
		Lugares:     v.Lugares,
	}
}

func (s *actividadService) GetAllActividades() (dto.ActividadesDTO, error) {
	var actividades model.ActividadesVista = actividad.GetAllActividades()
	var actividadesDTO dto.ActividadesDTO = make(dto.ActividadesDTO, len(actividades))

	for i, v := range actividades {
		actividadesDTO[i] = convertirADTO(v)
	}

	return actividadesDTO, nil
}

func (s *actividadService) GetActividadesByParams(params map[string]any) (dto.ActividadesDTO, error) {
	var actividades model.ActividadesVista = actividad.GetActividadesByParams(params)
	var actividadesDTO dto.ActividadesDTO = make(dto.ActividadesDTO, len(actividades))

	for i, v := range actividades {
		actividadesDTO[i] = convertirADTO(v)
	}

	return actividadesDTO, nil
}

func (s *actividadService) GetActividadByID(id int) (dto.ActividadDTO, error) {
	var actividad model.ActividadVista = actividad.GetActividadById(id)
	if actividad.Id == 0 {
		return dto.ActividadDTO{}, fmt.Errorf("actividad con ID %d no encontrada", id)
	}

	return convertirADTO(actividad), nil
}

func (s *actividadService) CreateActividad(actividadDTO dto.ActividadDTO) error {
	log.Printf("Recibiendo DTO para crear actividad: %+v\n", actividadDTO)

	if err := validarCamposBasicos(actividadDTO); err != nil {
		return err
	}

	horaInicio, horaFin, err := parsearHoras(actividadDTO.HoraInicio, actividadDTO.HoraFin)
	if err != nil {
		return err
	}

	nuevaActividad := model.Actividad{
		Titulo:        actividadDTO.Titulo,
		Descripcion:   actividadDTO.Descripcion,
		Cupo:          actividadDTO.Cupo,
		Dia:           actividadDTO.Dia,
		HorarioInicio: horaInicio,
		HorarioFinal:  horaFin,
		FotoUrl:       actividadDTO.FotoUrl,
		Instructor:    actividadDTO.Instructor,
		Categoria:     actividadDTO.Categoria,
	}

	log.Printf("Creando actividad: %+v\n", nuevaActividad)
	return actividad.CreateActividad(nuevaActividad)
}

func (s *actividadService) UpdateActividad(actividadDTO dto.ActividadDTO) error {
	log.Printf("Recibiendo DTO para actualizar actividad: %+v\n", actividadDTO)

	if err := validarCamposBasicos(actividadDTO); err != nil {
		return err
	}

	horaInicio, horaFin, err := parsearHoras(actividadDTO.HoraInicio, actividadDTO.HoraFin)
	if err != nil {
		return err
	}

	actividadActualizada := model.Actividad{
		Id:            actividadDTO.Id,
		Titulo:        actividadDTO.Titulo,
		Descripcion:   actividadDTO.Descripcion,
		Cupo:          actividadDTO.Cupo,
		Dia:           actividadDTO.Dia,
		HorarioInicio: horaInicio,
		HorarioFinal:  horaFin,
		FotoUrl:       actividadDTO.FotoUrl,
		Instructor:    actividadDTO.Instructor,
		Categoria:     actividadDTO.Categoria,
	}

	log.Infof("Actualizando actividad: %+v\n", actividadActualizada)
	if err := actividad.UpdateActividad(actividadActualizada); err != nil {
		return fmt.Errorf("error al actualizar actividad: %v", err)
	}

	return nil
}

func (s *actividadService) DeleteActividad(id uint) error {
	return actividad.DeleteActividad(id)
}
