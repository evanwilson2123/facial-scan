package controllers

import (
	"context"
	"facial-scan/models"
	"facial-scan/utils"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AccountCreation(c *fiber.Ctx) error {

	uid := c.Locals("user_id").(string)

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		log.Println("Error parsing body of patch request")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "unable to parse body"})
	}

	// Check if the username is being updated and if it already exists
	if username, ok := updateData["username"].(string); ok {
		iter := utils.FirestoreClient.Collection("users").Where("Username", "==", username).Documents(context.Background())
		for {
			doc, err := iter.Next()
			if err != nil {
				break
			}
			if doc.Ref.ID != uid {
				log.Println("Username already exists:", username)
				return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
			}
		}
	}

	docRef := utils.FirestoreClient.Collection("users").Doc(uid)
	doc, err := docRef.Get(context.Background())
	if err != nil {
		log.Println("error getting user document from database:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error finding user in database"})
	}
	var existingUser models.User
	if err := doc.DataTo(&existingUser); err != nil {
		log.Println("Error unmarshalling user data from Firestore:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get user from database"})
	}

	// Update only the fields that are present in the request body
	if age, ok := updateData["age"].(int); ok {
		existingUser.Age = age
	}
	if username, ok := updateData["username"].(string); ok {
		existingUser.Username = username
	}
	if gender, ok := updateData["gender"].(string); ok {
		existingUser.Gender = gender
	}

	if _, err := docRef.Set(context.Background(), existingUser); err != nil {
		log.Println("Error saving user account details to Firestore")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error updating user"})
	}

	log.Println("User updated successfully")
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User account updated successfully"})
}

func HasUsername(c *fiber.Ctx) error {
	uid := c.Locals("user_id").(string)

	docRef := utils.FirestoreClient.Collection("users").Doc(uid)
	doc, err := docRef.Get(context.Background())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"error fetching user from database"})
	}
	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"error unmarshalling data"})
	}
	if user.Username == "" {
		return c.Status(http.StatusOK).JSON(fiber.Map{"message":false})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message":true})
}


func DashBoardLoad(c *fiber.Ctx) error {
	uid := c.Locals("user_id").(string)
	user := getUserFromDatabase(uid)
	if user.UID == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"unable to find user from database"})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"username": user.Username,
	})
}

func GetScore(c *fiber.Ctx) error {
	uid := c.Locals("user_id").(string)
	if uid == "" {
		log.Printf("User ID is blank")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "user ID is missing"})
	}

	log.Printf("Fetching user for ID: %s", uid)
	user := getUserFromDatabase(uid)
	if user.UID == "" {
		log.Printf("No user found for ID: %s", uid)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "no user found"})
	}

	log.Printf("Querying for document with UserId: %s", uid)
	query := utils.FirestoreClient.Collection("images").Where("UserId", "==", uid).Limit(1)
	iter := query.Documents(context.Background())
	doc, err := iter.Next()
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Printf("Document not found for user ID: %s", uid)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "image not found in database"})
		}
		log.Printf("Error querying documents: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "unable to get image from database"})
	}

	var imageData ImageDataStore
	if err := doc.DataTo(&imageData); err != nil {
		log.Printf("Error unmarshalling image data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error unmarshalling image data"})
	}

	log.Printf("Successfully fetched and unmarshalled image data for user ID: %s", uid)
	return c.Status(http.StatusOK).JSON(imageData)
}