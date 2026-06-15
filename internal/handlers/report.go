package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"monthly-journal/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SendReport(db *gorm.DB, emailService *services.EmailService, reportService *services.ReportService) gin.HandlerFunc {
	return func(c *gin.Context) {
		month := time.Now().Format("2006-01")

		htmlBody, total, count, err := reportService.GenerateHTMLReport(month)
		if err != nil {
			log.Printf("Error generating report: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
			return
		}

		recipients := strings.Split(emailService.EmailFrom[0:len(emailService.EmailFrom)], ",")
		if len(emailService.EmailFrom) > 0 {
			recipients = []string{
				"nurdahliana86@gmail.com",
				"waonex86@gmail.com",
			}
		}

		subject := "Monthly Expense Report - " + month

		if err := emailService.SendEmail(recipients, subject, htmlBody); err != nil {
			log.Printf("Error sending email: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     "sent",
			"month":      month,
			"total":      total,
			"count":      count,
			"recipients": recipients,
			"sent_at":    time.Now().Format("2006-01-02T15:04:05Z"),
		})
	}
}
