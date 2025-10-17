package usuario

import (
	"net/http"
	"proyecto-integrador/dto"
	"proyecto-integrador/services"
	"strings"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func Login(ctx *gin.Context) {
	var loginJSON dto.UsuarioLoginDTO
	if err := ctx.BindJSON(&loginJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto"})
		log.Debug("LoginDTO:", loginJSON)
		return
	}

	token, err := services.UsuarioService.GenerateToken(loginJSON.Username, loginJSON.Password)
	if err != nil {
		log.Debug(err)
		if err == services.IncorrectCredentialsError {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ocurrio un error en el servidor"})
			return
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   1800, // en segundos
	})
}

func Register(ctx *gin.Context) {
	var datos dto.UsuarioMinDTO
	if err := ctx.BindJSON(&datos); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto"})
		log.Debug("LoginDTO:", datos)
		return
	}

	err := services.UsuarioService.RegisterUser(datos)
	if err != nil {
		log.Errorf("Error al registrar un usuario: %s\nDTO: %v", err.Error(), datos)

		errString := strings.ToLower(err.Error())
		if strings.Contains(errString, "error 1062") {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "El usuario ya está registrado"})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al registrarse"})
		}

		return
	}

	// una vez registrado el usuario le generamos un token
	token, err := services.UsuarioService.GenerateToken(datos.Username, datos.Password)
	if err != nil {
		log.Debug(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ocurrio un error en el servidor"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   1800, // en segundos
	})
}
