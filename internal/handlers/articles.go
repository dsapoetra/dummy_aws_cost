package handlers

import (
	"net/http"
	"strconv"
	"time"

	"dummy_aws_cost/internal/database"
	"dummy_aws_cost/internal/models"

	"github.com/gin-gonic/gin"
)

func GetArticles(c *gin.Context) {
	rows, err := database.DB.Query(
		"SELECT id, title, content, author, status, created_at, updated_at FROM articles ORDER BY created_at DESC",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var a models.Article
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Author, &a.Status, &a.CreatedAt, &a.UpdatedAt); err != nil {
			continue
		}
		articles = append(articles, a)
	}

	if articles == nil {
		articles = []models.Article{}
	}

	c.JSON(http.StatusOK, articles)
}

func GetArticle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var a models.Article
	err = database.DB.QueryRow(
		"SELECT id, title, content, author, status, created_at, updated_at FROM articles WHERE id = ?",
		id,
	).Scan(&a.ID, &a.Title, &a.Content, &a.Author, &a.Status, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, a)
}

func CreateArticle(c *gin.Context) {
	var a models.Article
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if a.Status == "" {
		a.Status = "draft"
	}

	now := time.Now()
	result, err := database.DB.Exec(
		"INSERT INTO articles (title, content, author, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		a.Title, a.Content, a.Author, a.Status, now, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	id, _ := result.LastInsertId()
	a.ID = id
	a.CreatedAt = now
	a.UpdatedAt = now

	c.JSON(http.StatusCreated, a)
}

func UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var a models.Article
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	now := time.Now()
	result, err := database.DB.Exec(
		"UPDATE articles SET title = ?, content = ?, author = ?, status = ?, updated_at = ? WHERE id = ?",
		a.Title, a.Content, a.Author, a.Status, now, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	a.ID = id
	a.UpdatedAt = now
	c.JSON(http.StatusOK, a)
}

func DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM articles WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
}
