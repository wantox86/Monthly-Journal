package routes

import (
	"monthly-journal/internal/config"
	"monthly-journal/internal/handlers"
	"monthly-journal/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Services
	emailSvc := services.NewEmailService(cfg)
	reportSvc := services.NewReportService(db)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Expense routes
	r.POST("/api/expenses", handlers.CreateExpense(db))
	r.GET("/api/expenses", handlers.GetExpenses(db))
	r.DELETE("/api/expenses/:id", handlers.DeleteExpense(db))

	// Report routes
	r.POST("/api/reports/send", handlers.SendReport(db, emailSvc, reportSvc))

	return r
}
