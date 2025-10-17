package dto

type ActividadMinDTO struct {
	Id          uint   `json:"id_actividad"`
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Cupo        uint   `json:"cupo"`
	Dia         string `json:"dia"`
	HoraInicio  string `json:"hora_inicio"`
	HoraFin     string `json:"hora_fin"`
	FotoUrl     string `json:"foto_url"`
	Instructor  string `json:"instructor"`
	Categoria   string `json:"categoria"`
}

type ActividadesMinDTO []ActividadMinDTO
