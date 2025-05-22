package handlers

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/khekrn/apprunner-fiber/services"
)

type FileHandler struct {
	s3Service *services.S3Service
}

func NewFileHandler(s3Service *services.S3Service) *FileHandler {
	return &FileHandler{
		s3Service: s3Service,
	}
}

func (h *FileHandler) getFileKey(userID, filename string) string {
	return fmt.Sprintf("users/%s/files/%s", userID, filename)
}

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Get file from multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer src.Close()

	// Read file data
	data, err := io.ReadAll(src)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read file",
		})
	}

	// Determine content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Create file key
	fileKey := h.getFileKey(userID, file.Filename)

	// Upload to S3
	fileInfo, err := h.s3Service.PutObject(c.Context(), fileKey, data, contentType, map[string]string{
		"user-id":       userID,
		"original-name": file.Filename,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Enhance file info
	fileInfo.UserID = userID
	fileInfo.Filename = file.Filename

	return c.Status(fiber.StatusCreated).JSON(fileInfo)
}

func (h *FileHandler) GetFile(c *fiber.Ctx) error {
	userID := c.Params("userId")
	filename := c.Params("filename")

	if userID == "" || filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID and filename are required",
		})
	}

	fileKey := h.getFileKey(userID, filename)

	data, fileInfo, err := h.s3Service.GetObject(c.Context(), fileKey)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Set appropriate headers
	c.Set("Content-Type", fileInfo.ContentType)
	c.Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size))
	c.Set("ETag", fileInfo.ETag)

	// Return file data
	return c.Send(data)
}

func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
	userID := c.Params("userId")
	filename := c.Params("filename")

	if userID == "" || filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID and filename are required",
		})
	}

	fileKey := h.getFileKey(userID, filename)

	if err := h.s3Service.DeleteObject(c.Context(), fileKey); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *FileHandler) ListUserFiles(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	prefix := fmt.Sprintf("users/%s/files/", userID)
	files, err := h.s3Service.ListObjects(c.Context(), prefix)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Enhance file info with user-friendly data
	for _, file := range files {
		file.UserID = userID
		file.Filename = path.Base(file.Key)
		// Remove prefix from key for cleaner display
		file.Key = strings.TrimPrefix(file.Key, prefix)
	}

	return c.JSON(fiber.Map{
		"files": files,
		"count": len(files),
	})
}
