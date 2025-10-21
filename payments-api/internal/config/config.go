package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	MongoURI      string
	MongoDatabase string
	StripeKey     string
	MercadoPagoKey string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		Port:           getEnv("PORT", "8083"),
		MongoURI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:  getEnv("MONGO_DATABASE", "payments"),
		StripeKey:      getEnv("STRIPE_SECRET_KEY", ""),
		MercadoPagoKey: getEnv("MERCADOPAGO_ACCESS_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
