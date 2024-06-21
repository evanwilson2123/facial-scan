package controllers

import (
	"context"
	"facial-scan/models"
	"facial-scan/utils"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
)

type LeaderBoard struct {
	Username      string         `json:"username"`
	ImageResponse ImageDataStore `json:"image_response"`
}

// GetLeaderboardMale retrieves the top 5 male users by high score and their highest scoring images
func GetLeaderboardMale(c *fiber.Ctx) error {
	ctx := context.Background()
	query := utils.FirestoreClient.Collection("users").Where("Gender", "==", "Male").OrderBy("HighScore", firestore.Desc).Limit(5)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var users []models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating user documents: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating user documents"})
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			log.Printf("Error unmarshalling user data: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error unmarshalling user data"})
		}
		users = append(users, user)
		log.Printf("Fetched user: %s, High Score: %f, UserID: %s", user.Username, user.HighScore, user.UID)
	}

	if len(users) == 0 {
		log.Println("No documents found for male users with high scores.")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "No documents found"})
	}

	var leaderboard []LeaderBoard
	for _, user := range users {
		log.Printf("Processing user: %s, High Score: %f, UserID: %s", user.Username, user.HighScore, user.UID)
		processUserImage(ctx, &leaderboard, user)
	}

	log.Printf("Final leaderboard: %+v", leaderboard)
	return c.Status(http.StatusOK).JSON(leaderboard)
}

// processUserImage retrieves the highest scoring image for a given user
func processUserImage(ctx context.Context, leaderboard *[]LeaderBoard, user models.User) {
	imageQuery := utils.FirestoreClient.Collection("images").Where("UserId", "==", user.UID).OrderBy("TotalScore", firestore.Desc).Limit(1)
	imageIter := imageQuery.Documents(ctx)
	defer imageIter.Stop()

	for {
		imageDoc, err := imageIter.Next()
		if err == iterator.Done {
			log.Printf("No image found for user ID: %s", user.UID)
			break
		}
		if err != nil {
			log.Printf("Error iterating image documents: %v", err)
			return
		}

		var imageDataStore ImageDataStore
		if err := imageDoc.DataTo(&imageDataStore); err != nil {
			log.Printf("Error unmarshalling image data: %v", err)
			return
		}

		log.Printf("Successfully fetched image data for user ID: %s, Total Score: %f", user.UID, imageDataStore.TotalScore)
		log.Printf("Total score for user %s: %f", user.Username, imageDataStore.TotalScore)

		*leaderboard = append(*leaderboard, LeaderBoard{
			Username:      user.Username,
			ImageResponse: imageDataStore,
		})
	}
}

func GetLeaderboardFemale(c *fiber.Ctx) error {
	ctx := context.Background()
	query := utils.FirestoreClient.Collection("users").Where("Gender", "==", "Female").OrderBy("HighScore", firestore.Desc).Limit(5)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var users []models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating user documents: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating user documents"})
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			log.Printf("Error unmarshalling user data: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error unmarshalling user data"})
		}
		users = append(users, user)
		log.Printf("Fetched user: %s, High Score: %f, UserID: %s", user.Username, user.HighScore, user.UID)
	}

	if len(users) == 0 {
		log.Println("No documents found for male users with high scores.")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "No documents found"})
	}

	var leaderboard []LeaderBoard
	for _, user := range users {
		processUserImage(ctx, &leaderboard, user)
	}

	log.Printf("Final leaderboard: %+v", leaderboard)
	return c.Status(http.StatusOK).JSON(leaderboard)
}

