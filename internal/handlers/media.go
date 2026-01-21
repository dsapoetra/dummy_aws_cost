package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"dummy_aws_cost/internal/database"
	"dummy_aws_cost/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUploadDir() string {
	dir := os.Getenv("UPLOAD_DIR")
	if dir == "" {
		dir = "./uploads"
	}
	return dir
}

func GetMedia(c *gin.Context) {
	rows, err := database.DB.Query(
		"SELECT id, filename, original_name, mime_type, size, created_at FROM media ORDER BY created_at DESC",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch media"})
		return
	}
	defer rows.Close()

	var media []models.Media
	for rows.Next() {
		var m models.Media
		if err := rows.Scan(&m.ID, &m.Filename, &m.OriginalName, &m.MimeType, &m.Size, &m.CreatedAt); err != nil {
			continue
		}
		media = append(media, m)
	}

	if media == nil {
		media = []models.Media{}
	}

	c.JSON(http.StatusOK, media)
}

func UploadMedia(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	uploadDir := GetUploadDir()
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	mimeType := file.Header.Get("Content-Type")
	now := time.Now()

	result, err := database.DB.Exec(
		"INSERT INTO media (filename, original_name, mime_type, size, created_at) VALUES (?, ?, ?, ?, ?)",
		filename, file.Filename, mimeType, file.Size, now,
	)
	if err != nil {
		os.Remove(dst)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save media record"})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, models.Media{
		ID:           id,
		Filename:     filename,
		OriginalName: file.Filename,
		MimeType:     mimeType,
		Size:         file.Size,
		CreatedAt:    now,
	})
}

func DeleteMedia(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	var filename string
	err = database.DB.QueryRow("SELECT filename FROM media WHERE id = ?", id).Scan(&filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM media WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	filePath := filepath.Join(GetUploadDir(), filename)
	os.Remove(filePath)

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted"})
}
