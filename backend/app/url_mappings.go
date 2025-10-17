package app

import (
	"proyecto-integrador/controllers/actividad"
	"proyecto-integrador/controllers/inscripcion"
	"proyecto-integrador/controllers/usuario"
)

func MapURLs() {
	// actividades
	router.GET("/actividades", actividad.GetAllActividades)
	router.GET("/actividades/:id", actividad.GetActividadById)
	router.GET("/actividades/buscar", actividad.GetActividadesByParams)
	router.POST("/actividades", JWTValidationMiddle, IsAdminMiddle, actividad.CreateActividad)
	router.PUT("/actividades/:id", JWTValidationMiddle, IsAdminMiddle, actividad.UpdateActividad)
	router.DELETE("/actividades/:id", JWTValidationMiddle, IsAdminMiddle, actividad.DeleteActividad)

	// usuarios
	router.POST("/login", usuario.Login)
	router.POST("/register", usuario.Register)

	// inscripciones
	router.GET("/inscripciones", JWTValidationMiddle, inscripcion.GetAllInscripciones)
	router.POST("/inscripciones", JWTValidationMiddle, inscripcion.InscribirUsuario)
	router.DELETE("/inscripciones", JWTValidationMiddle, inscripcion.DesinscribirUsuario)
}
