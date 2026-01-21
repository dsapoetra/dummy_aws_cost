package handlers

import (
	"net/http"
	"strconv"
	"time"

	"dummy_aws_cost/internal/database"
	"dummy_aws_cost/internal/models"

	"github.com/gin-gonic/gin"
)

func GetPages(c *gin.Context) {
	rows, err := database.DB.Query(
		"SELECT id, title, slug, content, created_at, updated_at FROM pages ORDER BY created_at DESC",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages"})
		return
	}
	defer rows.Close()

	var pages []models.Page
	for rows.Next() {
		var p models.Page
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			continue
		}
		pages = append(pages, p)
	}

	if pages == nil {
		pages = []models.Page{}
	}

	c.JSON(http.StatusOK, pages)
}

func GetPage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var p models.Page
	err = database.DB.QueryRow(
		"SELECT id, title, slug, content, created_at, updated_at FROM pages WHERE id = ?",
		id,
	).Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

func CreatePage(c *gin.Context) {
	var p models.Page
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	now := time.Now()
	result, err := database.DB.Exec(
		"INSERT INTO pages (title, slug, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		p.Title, p.Slug, p.Content, now, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create page. Slug may already exist."})
		return
	}

	id, _ := result.LastInsertId()
	p.ID = id
	p.CreatedAt = now
	p.UpdatedAt = now

	c.JSON(http.StatusCreated, p)
}

func UpdatePage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var p models.Page
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	now := time.Now()
	result, err := database.DB.Exec(
		"UPDATE pages SET title = ?, slug = ?, content = ?, updated_at = ? WHERE id = ?",
		p.Title, p.Slug, p.Content, now, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	p.ID = id
	p.UpdatedAt = now
	c.JSON(http.StatusOK, p)
}

func DeletePage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM pages WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted"})
}
