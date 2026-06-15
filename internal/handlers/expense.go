package handlers

import (
	"net/http"
	"strconv"
	"time"

	"monthly-journal/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateExpense(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Description string `json:"description" binding:"required"`
			Amount      int    `json:"amount" binding:"required,gt=0"`
			Sender      string `json:"sender"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}

		if req.Amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
			return
		}

		sender := req.Sender
		if sender == "" {
			sender = "Unknown"
		}

		expense := models.Expense{
			Description: req.Description,
			Amount:      req.Amount,
			Sender:      sender,
			MonthYear:   time.Now().Format("2006-01"),
			Date:        time.Now(),
		}

		if err := db.Create(&expense).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, expense)
	}
}

func GetExpenses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		month := c.DefaultQuery("month", time.Now().Format("2006-01"))

		var expenses []models.Expense
		var totalResult struct {
			Total int64
		}

		db.Where("month_year = ?", month).Find(&expenses)
		db.Model(&models.Expense{}).Where("month_year = ?", month).Select("COALESCE(SUM(amount), 0) as total").Scan(&totalResult)

		c.JSON(http.StatusOK, gin.H{
			"month":    month,
			"total":    totalResult.Total,
			"count":    len(expenses),
			"expenses": expenses,
		})
	}
}

func DeleteExpense(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		expenseID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
			return
		}

		if err := db.Delete(&models.Expense{}, expenseID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "deleted",
			"id":     id,
		})
	}
}
