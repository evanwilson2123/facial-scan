package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type AzureResponse struct {
	FaceId string `json:"faceId"`
}

func UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file"})
	}

	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to open file"})
	}
	defer fileHeader.Close()

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to create form file"})
	}
	_, err = io.Copy(part, fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to copy file"})
	}
	writer.Close()

	azureEndpoint := "https://<your-region>.api.cognitive.microsoft.com/face/v1.0/detect?returnFaceId=true&returnFaceAttributes=age,gender"
	req, err := http.NewRequest("POST", azureEndpoint, buf)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to create request"})
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("AZURE_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to send request"})
	}
	defer resp.Body.Close()

	var azureResponse []AzureResponse
	if err := json.NewDecoder(resp.Body).Decode(&azureResponse); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to parse response"})
	}

	if len(azureResponse) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No face detected"})
	}

	return c.JSON(fiber.Map{"faceId": azureResponse[0].FaceId})
}
