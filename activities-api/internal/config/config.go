package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	MySQL MySQLConfig
	JWT   JWTConfig
	// TODO: Agregar RabbitMQ config cuando el equipo lo implemente
	// RabbitMQ RabbitMQConfig
}

type MySQLConfig struct {
	User   string
	Pass   string
	Host   string
	Port   string
	Schema string
}

type JWTConfig struct {
	Secret string
}

// TODO: Descomentar cuando el equipo implemente RabbitMQ
// type RabbitMQConfig struct {
// 	Host  string
// 	Port  string
// 	User  string
// 	Pass  string
// 	Queue string
// }

func Load() Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	return Config{
		Port: getEnv("PORT", "8082"),
		MySQL: MySQLConfig{
			User:   getEnv("DB_USER", "root"),
			Pass:   getEnv("DB_PASS", ""),
			Host:   getEnv("DB_HOST", "localhost"),
			Port:   getEnv("DB_PORT", "3306"),
			Schema: getEnv("DB_SCHEMA", "proyecto_integrador"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "my-secret-key"),
		},
		// TODO: Descomentar cuando el equipo implemente RabbitMQ
		// RabbitMQ: RabbitMQConfig{
		// 	Host:  getEnv("RABBITMQ_HOST", "localhost"),
		// 	Port:  getEnv("RABBITMQ_PORT", "5672"),
		// 	User:  getEnv("RABBITMQ_USER", "admin"),
		// 	Pass:  getEnv("RABBITMQ_PASS", "admin"),
		// 	Queue: getEnv("RABBITMQ_QUEUE", "gym-events"),
		// },
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
