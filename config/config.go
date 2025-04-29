package config

import (
	"os"
	"strconv"
	"log"
	"time"
)

type Config struct {
	DBConfig
	VectorStoreConfig
	JWTConfig
	ServerConfig
	AIConfig
}

type VectorStoreConfig struct {
	APIKey string
	IndexName string
	Namespace string
}

type DBConfig struct {
	Addr string
	Port string
	User string
	Password string
	DBName string
	MaxIdleConns int
	MaxOpenConns int
	MaxIdleTime string
}

type AIConfig struct {
	APIKey string
}

type JWTConfig struct {
	SecretKey string
	TokenExpiry time.Duration
	TokenIssuer string
	TokenAudience string
}

type ServerConfig struct {
	Port string
	Adress string
}

func Get() *Config {

	maxIdleConns, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if err != nil {
		log.Fatal("Error converting DB_MAX_IDLE_CONNS to int")
	}
	maxOpenConns, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if err != nil {
		log.Fatal("Error converting DB_MAX_OPEN_CONNS to int")
	}

	return &Config{
		DBConfig: DBConfig{
			Addr: os.Getenv("DB_ADDR"),
			User: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName: os.Getenv("DB_NAME"),
			MaxIdleConns: maxIdleConns,
			MaxOpenConns: maxOpenConns,
			MaxIdleTime: os.Getenv("DB_MAX_IDLE_TIME"),
			Port: os.Getenv("DB_PORT"),
		},
		VectorStoreConfig: VectorStoreConfig{
			APIKey: os.Getenv("PINECONE_API_KEY"),
			IndexName: os.Getenv("PINECONE_INDEX_NAME"),
			Namespace: os.Getenv("PINECONE_NAMESPACE"),
		},
		ServerConfig: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		JWTConfig: JWTConfig{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
			TokenExpiry: time.Hour * 1,
			TokenIssuer: os.Getenv("JWT_TOKEN_ISSUER"),
			TokenAudience: os.Getenv("JWT_TOKEN_AUDIENCE"),
		},
		AIConfig: AIConfig{
			APIKey: os.Getenv("OPENAI_API_KEY"),
		},
	}
}


