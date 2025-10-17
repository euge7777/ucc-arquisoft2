package dto

type UsuarioMinDTO struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Username string `json:"username"`
	Password string `json:"password"`
}
