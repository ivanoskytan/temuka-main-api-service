package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/util/database"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	postgres, err := database.NewPostgreSQL(
		os.Getenv("PG_HOST"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DATABASE"),
	)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := postgres.DB.AutoMigrate(
		&model.User{},
		&model.Community{},
		&model.Post{},
		&model.Conversation{},
		&model.Comment{},
		&model.CommunityMember{},
		&model.CommunityPost{},
		&model.Moderator{},
		&model.Participant{},
		&model.UserFollow{},
		&model.Notification{},
		&model.Report{},
		&model.Location{},
		&model.University{},
		&model.Review{},
		&model.Major{},
		&model.MajorReview{},
		&model.CommunityRule{},
		&model.PostLike{},
		&model.PostSave{},
	); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	log.Println("Database migration completed successfully.")
}
