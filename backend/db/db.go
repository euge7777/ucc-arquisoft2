package db

import (
	"fmt"
	"os"
	"proyecto-integrador/model"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db        *gorm.DB
	db_user   string
	db_pass   string
	db_host   string
	db_port   string
	db_schema string

	iniciar_conexion sync.Once
)

func GetInstance() *gorm.DB {
	iniciar_conexion.Do(func() {
		db_user = os.Getenv("DB_USER")
		if db_user == "" {
			db_user = "root"
		}
		db_pass = os.Getenv("DB_PASS")
		db_host = os.Getenv("DB_HOST")
		if db_host == "" {
			db_host = "localhost"
		}
		db_port = os.Getenv("DB_PORT")
		if db_port == "" {
			db_port = "3306"
		}
		db_schema = os.Getenv("DB_SCHEMA")
		if db_schema == "" {
			db_schema = "proyecto_integrador"
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_user, db_pass, db_host, db_port, db_schema)
		log.Info("Conectando a la base de datos con dsn: ", dsn)

		// reintentamos conectarnos a la BDD varias veces
		for i := range 10 {
			time.Sleep(3 * time.Second)
			log.Debugf("Intentando conectar (%d/%d)\n", i+1, 10)

			var err error
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})

			if err != nil {
				log.Errorf("Error al conectar a la base de datos: %v", err)
				log.Error("No se pudo establecer conexion a la BDD")
				continue
			}

			break
		}

		log.Info("Conexion a base de datos establecida")
	})

	return db
}

func StartDbEngine() {
	// iniciamos el pool de conexiónes
	GetInstance()

	db.AutoMigrate(&model.Actividad{})
	db.AutoMigrate(&model.Inscripcion{})
	db.AutoMigrate(&model.Usuario{})

	// creamos una vista de actividades que permita verificar la cantidad de lugares restantes
	db.Exec(`DROP VIEW IF EXISTS actividads_lugares`)
	db.Exec(`
	CREATE VIEW actividads_lugares AS
	SELECT *, cupo - (SELECT COUNT(*) FROM inscripcions ins WHERE ins.id_actividad = ac.id_actividad AND ins.is_activa) AS lugares
	FROM actividads ac
	`)

	log.Info("Terminada la migracion de las tablas")
}
