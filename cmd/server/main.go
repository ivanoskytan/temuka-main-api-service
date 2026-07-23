package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	router "github.com/temuka-api-service/api"
	"github.com/temuka-api-service/internal/constant"
	"github.com/temuka-api-service/util/database"
	"github.com/temuka-api-service/util/file_storage"
	"github.com/temuka-api-service/util/key_value_store"
	"github.com/temuka-api-service/util/queue"
	ws "github.com/temuka-api-service/util/websocket"
)

func EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	postgres, err := database.NewPostgreSQL(
		os.Getenv(constant.EnvPgHost),
		os.Getenv(constant.EnvPgUser),
		os.Getenv(constant.EnvPgPass),
		os.Getenv(constant.EnvPgPort),
		os.Getenv(constant.EnvPgDB),
	)
	if err != nil {
		log.Fatalf("Error initiating relational database: %v", err)
	}

	redis, err := key_value_store.NewRedisConnection(
		os.Getenv(constant.EnvRedisHost),
		os.Getenv(constant.EnvRedisUser),
		os.Getenv(constant.EnvRedisPass),
	)
	if err != nil {
		log.Fatalf("Error initiating key value store: %v", err)
	}

	storage, err := file_storage.NewS3(
		os.Getenv(constant.EnvAWSRegion),
		os.Getenv(constant.EnvAWSAccessKeyID),
		os.Getenv(constant.EnvAWSSecretAccessKey),
		os.Getenv(constant.EnvS3Bucket),
		os.Getenv(constant.EnvS3Endpoint),
	)
	if err != nil {
		log.Fatalf("Error initiating file storage: %v", err)
	}

	rmq, err := queue.NewRabbitMQConnection(
		os.Getenv(constant.EnvRabbitMQURL),
	)
	if err != nil {
		log.Fatalf("Error initiating message queue: %v", err)
	}

	mqChannel, err := rmq.Channel()
	if err != nil {
		log.Fatalf("Error creating message queue channel: %v", err)
	}

	hub := ws.NewHub()
	go hub.Run()

	router := router.Routes(*postgres, *redis, *storage, *mqChannel, hub)
	protectedRoutes := EnableCors(router)

	http.Handle("/", protectedRoutes)
	log.Println("Server is listening on port 3200")
	log.Fatal(http.ListenAndServe("0.0.0.0:3200", nil))
}
