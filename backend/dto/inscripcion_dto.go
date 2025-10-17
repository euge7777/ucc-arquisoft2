package dto

type InscripcionDTO struct {
	IdUsuario        uint   `json:"id_usuario"`
	IdActividad      uint   `json:"id_actividad"`
	FechaInscripcion string `json:"fecha_inscripcion"`
	IsActiva         bool   `json:"is_activa"`
}

type InscripcionesDTO []InscripcionDTO
