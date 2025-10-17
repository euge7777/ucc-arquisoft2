package dto

type InscripcionMinDTO struct {
	IdUsuario   uint   `json:"id_usuario"`
	IdActividad uint   `json:"id_actividad"`
	IsActiva    string `json:"is_activa"`
}

type InscripcionesMinDTO []InscripcionMinDTO
