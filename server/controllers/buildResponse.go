package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"facial-scan/models"
	"facial-scan/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
)

type ImageData struct {
	UserID string `json:"user_id"`
	URL    string `json:"url"`
}

type ImageRequest struct {
	ImageUrl string `json:"imageUrl"`
	UserId   string `json:"userId"`
}

type ImageResponse struct {
	Symmetry              float64 `json:"symmetry"`
	FacialDefinition      float64 `json:"facial_definition"`
	Jawline               float64 `json:"jawline"`
	Cheekbones            float64 `json:"cheekbones"`
	JawlineToCheekbones   float64 `json:"jawline_to_cheekbones"`
	CanthalTilt           float64 `json:"canthal_tilt"`
	ProportionAndRatios   float64 `json:"proportion_and_ratios"`
	SkinQuality           float64 `json:"skin_quality"`
	LipFullness           float64 `json:"lip_fullness"`
	FacialFat             float64 `json:"facial_fat"`
	CompleteFacialHarmony float64 `json:"complete_facial_harmony"`
	TotalScore            float32 `json:"total_score"`
	ImageURL 			  string  `json:"image_url"`

}

type ImageDataStore struct {
	ImageURL 			  string  `json:"image_url"`
	UserId				  string  `json:"user_id"`
	TotalScore			  float32 `json:"total_score"`	
	Symmetry              float64 `json:"symmetry"`
	FacialDefinition      float64 `json:"facial_definition"`
	Jawline               float64 `json:"jawline"`
	Cheekbones            float64 `json:"cheekbones"`
	JawlineToCheekbones   float64 `json:"jawline_to_cheekbones"`
	CanthalTilt           float64 `json:"canthal_tilt"`
	ProportionAndRatios   float64 `json:"proportion_and_ratios"`
	SkinQuality           float64 `json:"skin_quality"`
	LipFullness           float64 `json:"lip_fullness"`
	FacialFat             float64 `json:"facial_fat"`
	CompleteFacialHarmony float64 `json:"complete_facial_harmony"`			
}

func UploadSaveAndRespond(c *fiber.Ctx) error {
	var bucketName = os.Getenv("BUCKET_NAME")
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file"})
	}

	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "unable to open file"})
	}
	defer fileHeader.Close()

	userID := c.Locals("user_id").(string)
	bucket := utils.StorageClient.Bucket(bucketName)
	fileName := fmt.Sprintf("images/%d_%s", time.Now().Unix(), file.Filename)
	object := bucket.Object(fileName)
	writer := object.NewWriter(context.Background())

	if _, err := io.Copy(writer, fileHeader); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to write to bucket"})
	}

	if err := writer.Close(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to close bucket writer"})
	}

	if err := object.ACL().Set(context.Background(), storage.AllUsers, storage.RoleReader); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to set bucket ACL"})
	}
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName)

	imageRequest := ImageRequest{
		ImageUrl: publicURL,
	}

	jsonData, err := json.Marshal(imageRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error marshalling JSON"})
	}

	response, err := http.Post("http://localhost:3000/process-image", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error making request to image processing server"})
	}
	defer response.Body.Close()

	var imageResponse ImageResponse
	if err := json.NewDecoder(response.Body).Decode(&imageResponse); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error decoding response from image processing server"})
	}


	imageResponse.ImageURL = publicURL

	if imageResponse.TotalScore == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"unable to process image"})
	}

	if !setHighScore(userID, float64(imageResponse.TotalScore)) {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"error updating high score"})
	}

	imageData := ImageDataStore{
		ImageURL: publicURL,
		UserId: userID,
		TotalScore: float32(imageResponse.TotalScore),
		Symmetry: imageResponse.Symmetry,
		FacialDefinition: imageResponse.FacialDefinition,      
		Jawline: imageResponse.Jawline,               
		Cheekbones: imageResponse.Cheekbones,          
		JawlineToCheekbones: imageResponse.JawlineToCheekbones,
		CanthalTilt: imageResponse.CanthalTilt,           
		ProportionAndRatios: imageResponse.ProportionAndRatios,
		SkinQuality: imageResponse.SkinQuality,
		LipFullness: imageResponse.LipFullness,
		FacialFat: imageResponse.FacialFat,
		CompleteFacialHarmony: imageResponse.CompleteFacialHarmony,
	}

	ctx := context.Background()
	_, _, err = utils.FirestoreClient.Collection("images").Add(ctx, imageData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"error saving image to database"})
	}

	return c.JSON(imageResponse)
}


func getUserFromDatabase(userId string) models.User {
	docRef := utils.FirestoreClient.Collection("users").Doc(userId)
	doc, err := docRef.Get(context.Background())
	if err != nil {
		log.Println("error getting user document from database:", err)
		return models.User{}
	}
	var existingUser models.User
	if err := doc.DataTo(&existingUser); err != nil {
		log.Println("Error unmarshalling user data from Firestore:", err)
		return models.User{}
	}
	return existingUser
}

func setHighScore(uid string, totalScore float64) (bool) {
	user := getUserFromDatabase(uid)
	currentHighScore := user.HighScore
	if totalScore > currentHighScore {
		user.HighScore = totalScore
		_, err := utils.FirestoreClient.Collection("users").Doc(user.UID).Set(context.Background(), user)
		if err != nil {
			log.Printf("Error saving user to database: %v", err)
			return false
		}
		return true
	}
	return true
} 
