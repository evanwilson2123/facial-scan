package controllers

import (
	"context"
	"facial-scan/models"
	"facial-scan/utils"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {

	uid := c.Locals("user_id").(string)

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Error parsing user object: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"bad request"})
	}

	user.UID = uid
	user.CreatedAt = time.Now().Unix()

	_, err := utils.FirestoreClient.Collection("users").Doc(user.UID).Set(context.Background(), user)
	if err != nil {
		log.Printf("Error saving user to Firestore: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"failed to save user"})
	}
	log.Printf("User created successfully!: %v\n", user.Email)
	return c.Status(http.StatusCreated).JSON(fiber.Map{"message":"User created successfully"})
}

func HealthCheck(c *fiber.Ctx) error {
	return  c.Status(http.StatusAccepted).JSON(fiber.Map{"health":"good"})
}