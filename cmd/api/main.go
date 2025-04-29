package main

import (
	"log"
	"net/http"
	"time"
	"github.com/joaoramos09/go-ai/config"
	"github.com/joaoramos09/go-ai/ia/openai"
	"github.com/joaoramos09/go-ai/ia"
	"github.com/joaoramos09/go-ai/internal/database"
	"github.com/joaoramos09/go-ai/internal/http/chi"
	"github.com/joaoramos09/go-ai/internal/user"
	user_postgres "github.com/joaoramos09/go-ai/internal/user/postgres"
	"github.com/joho/godotenv"
	"github.com/joaoramos09/go-ai/internal/user/auth"
	auth_jwt "github.com/joaoramos09/go-ai/internal/user/auth/jwt"
	"github.com/joaoramos09/go-ai/internal/http/middlewares"
	"github.com/joaoramos09/go-ai/ia/pinecone"
	"context"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Get()

	p, err := database.NewPostgres(
		cfg.DBConfig.Addr,
		cfg.DBConfig.MaxIdleConns,
		cfg.DBConfig.MaxOpenConns,
		cfg.DBConfig.MaxIdleTime,
		cfg.DBConfig.Port,
	)

	pc := database.NewPinecone(context.Background(), cfg.VectorStoreConfig.IndexName, cfg.VectorStoreConfig.APIKey, cfg.VectorStoreConfig.Namespace)
	if err != nil {
		log.Fatalf("Error creating pinecone: %v", err)
	}
	userRepo := user_postgres.NewRepository(p)
	pineconeRepo := pinecone.NewRepository(pc)

	openaiService := openai.NewService(cfg.AIConfig.APIKey)
	aiService := ia.NewService(openaiService, pineconeRepo)
	userService := user.NewService(userRepo)
	tokenService := auth_jwt.NewService(cfg.JWTConfig.SecretKey, cfg.JWTConfig.TokenExpiry, cfg.JWTConfig.TokenIssuer, cfg.JWTConfig.TokenAudience)
	authService := auth.NewService(tokenService)
	middlewares := middlewares.NewService(authService, userService)

	r := chi.Handlers(aiService, userService, authService, middlewares)

	http.Handle("/", r)
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      http.DefaultServeMux,
	}

	log.Println("Starting server on port 8080")
	srv.ListenAndServe()

}
