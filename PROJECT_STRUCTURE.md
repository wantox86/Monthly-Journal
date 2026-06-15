# Project Structure & Implementation Guide

## Final Directory Structure

```
monthly-journal/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Load & parse environment variables
│   ├── models/
│   │   ├── models.go               # Register GORM models
│   │   └── expense.go              # Expense model struct
│   ├── database/
│   │   └── db.go                   # Database connection setup
│   ├── handlers/
│   │   ├── expense.go              # Expense CRUD handlers
│   │   └── report.go               # Report send handler
│   ├── services/
│   │   ├── email.go                # SMTP service
│   │   ├── report.go               # Report generation service
│   │   └── repository.go           # Database queries (optional)
│   └── routes/
│       └── routes.go               # Route registration
├── .env                            # Environment variables (LOCAL ONLY - gitignore)
├── .env.example                    # Template for .env
├── .gitignore                      # Exclude .env, .exe, etc
├── go.mod                          # Go module file
├── go.sum                          # Dependencies lock file
├── README.md                       # Setup & usage guide
├── SPRINTS.md                      # This file - sprint breakdown
├── claude.md                       # Project specs
└── PROJECT_STRUCTURE.md            # This file - structure guide
```

---

## Key Files to Create

### 1. cmd/server/main.go
**Purpose:** Application entry point

```go
package main

import (
    "log"
    "monthly-journal/internal/config"
    "monthly-journal/internal/database"
    "monthly-journal/internal/routes"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    db, err := database.Connect(cfg)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    router := routes.SetupRoutes(db, cfg)
    router.Run(":" + cfg.ServerPort)
}
```

### 2. internal/config/config.go
**Purpose:** Load environment variables

```go
package config

import "github.com/joho/godotenv"

type Config struct {
    // Database
    DBHost string
    DBPort string
    DBName string
    DBUser string
    DBPass string
    
    // SMTP
    SMTPHost string
    SMTPPort int
    SMTPUser string
    SMTPPass string
    
    // Email
    EmailFrom      string
    EmailRecipients string
    
    // Server
    ServerPort string
    ServerEnv  string
}

func LoadConfig() (*Config, error) {
    godotenv.Load()
    // Parse os.Getenv() untuk setiap variable
    // Return *Config struct
}
```

### 3. internal/models/expense.go
**Purpose:** Expense data model

```go
package models

import "time"

type Expense struct {
    ID          int       `gorm:"primaryKey" json:"id"`
    Date        time.Time `json:"date"`
    Description string    `json:"description"`
    Amount      int       `json:"amount"`
    Sender      string    `json:"sender"`
    MonthYear   string    `json:"month_year"`
    CreatedAt   time.Time `json:"created_at"`
}

func (Expense) TableName() string {
    return "expenses"
}
```

### 4. internal/database/db.go
**Purpose:** Database connection

```go
package database

import (
    "fmt"
    "monthly-journal/internal/config"
    "monthly-journal/internal/models"
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
    )
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Auto migration
    db.AutoMigrate(&models.Expense{})
    
    return db, nil
}
```

### 5. internal/handlers/expense.go
**Purpose:** Expense CRUD handlers

```go
package handlers

import (
    "monthly-journal/internal/models"
    "time"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func CreateExpense(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            Description string `json:"description" binding:"required"`
            Amount      int    `json:"amount" binding:"required"`
            Sender      string `json:"sender"`
        }
        
        if err := c.BindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
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
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(201, expense)
    }
}

func GetExpenses(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        month := c.DefaultQuery("month", time.Now().Format("2006-01"))
        
        var expenses []models.Expense
        var total int64
        
        db.Where("month_year = ?", month).Find(&expenses)
        db.Model(&models.Expense{}).Where("month_year = ?", month).Sum("amount", &total)
        
        c.JSON(200, gin.H{
            "month":    month,
            "total":    total,
            "count":    len(expenses),
            "expenses": expenses,
        })
    }
}

func DeleteExpense(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        
        if err := db.Delete(&models.Expense{}, id).Error; err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(200, gin.H{
            "status": "deleted",
            "id":     id,
        })
    }
}
```

### 6. internal/services/email.go
**Purpose:** SMTP email service

```go
package services

import (
    "fmt"
    "net/smtp"
    "monthly-journal/internal/config"
)

type EmailService struct {
    SMTPHost   string
    SMTPPort   int
    SMTPUser   string
    SMTPPass   string
    EmailFrom  string
}

func NewEmailService(cfg *config.Config) *EmailService {
    return &EmailService{
        SMTPHost:  cfg.SMTPHost,
        SMTPPort:  cfg.SMTPPort,
        SMTPUser:  cfg.SMTPUser,
        SMTPPass:  cfg.SMTPPass,
        EmailFrom: cfg.EmailFrom,
    }
}

func (es *EmailService) SendEmail(to []string, subject string, body string) error {
    auth := smtp.PlainAuth("", es.SMTPUser, es.SMTPPass, es.SMTPHost)
    addr := fmt.Sprintf("%s:%d", es.SMTPHost, es.SMTPPort)
    
    msg := fmt.Sprintf("From: %s\r\n", es.EmailFrom)
    msg += fmt.Sprintf("To: %s\r\n", to[0])
    msg += fmt.Sprintf("Subject: %s\r\n", subject)
    msg += "MIME-Version: 1.0\r\n"
    msg += "Content-Type: text/html; charset=UTF-8\r\n\r\n"
    msg += body
    
    return smtp.SendMail(addr, auth, es.EmailFrom, to, []byte(msg))
}
```

### 7. internal/services/report.go
**Purpose:** Report generation service

```go
package services

import (
    "fmt"
    "monthly-journal/internal/models"
    "time"
    "gorm.io/gorm"
)

type ReportService struct {
    db *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
    return &ReportService{db: db}
}

func (rs *ReportService) GenerateHTMLReport(month string) (string, int, int, error) {
    var expenses []models.Expense
    var total int64
    
    rs.db.Where("month_year = ?", month).Find(&expenses)
    rs.db.Model(&models.Expense{}).Where("month_year = ?", month).Sum("amount", &total)
    
    html := "<table border='1' cellpadding='10'>"
    html += "<tr><th>No</th><th>Tanggal</th><th>Deskripsi</th><th>Amount (Rp)</th><th>Pengirim</th></tr>"
    
    for i, e := range expenses {
        sender := e.Sender
        if sender == "" {
            sender = "NULL"
        }
        html += fmt.Sprintf(
            "<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
            i+1,
            e.Date.Format("02-01-2006 15:04"),
            e.Description,
            formatAmount(e.Amount),
            sender,
        )
    }
    
    html += "</table>"
    return html, len(expenses), int(total), nil
}

func formatAmount(amount int) string {
    // Format dengan comma: 82350 → 82,350
}
```

### 8. internal/routes/routes.go
**Purpose:** Route registration

```go
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
    
    // Routes
    r.POST("/api/expenses", handlers.CreateExpense(db))
    r.GET("/api/expenses", handlers.GetExpenses(db))
    r.DELETE("/api/expenses/:id", handlers.DeleteExpense(db))
    
    r.POST("/api/reports/send", handlers.SendReport(db, emailSvc, reportSvc))
    
    return r
}
```

---

## Dependencies to Install

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/joho/godotenv
```

---

## Environment Variables (.env)

```
DB_HOST=192.168.50.131
DB_PORT=3306
DB_NAME=monthly_bill
DB_USER=copilot
DB_PASS=copilot123

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password

EMAIL_FROM=waonex86@gmail.com
EMAIL_RECIPIENTS=nurdahliana86@gmail.com,waonex86@gmail.com

SERVER_PORT=8080
SERVER_ENV=development
```

---

## Build & Run

```bash
# Build
go build -o monthly-journal cmd/server/main.go

# Run
./monthly-journal

# Or direct run
go run cmd/server/main.go
```

---

## Testing

```bash
# Create expense
curl -X POST http://localhost:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{"description":"jajan kopi","amount":40000,"sender":"Nur Dahlia"}'

# Get expenses for current month
curl http://localhost:8080/api/expenses

# Get expenses for specific month
curl http://localhost:8080/api/expenses?month=2026-06

# Send report
curl -X POST http://localhost:8080/api/reports/send \
  -H "Content-Type: application/json" \
  -d '{"format":"csv"}'

# Delete expense
curl -X DELETE http://localhost:8080/api/expenses/1
```

---

## Notes
- Start dengan Sprint 0, jangan skip setup
- Verify setiap sprint sebelum lanjut ke sprint berikutnya
- Gunakan file ini sebagai reference untuk structure
- Setiap file dijelaskan purpose-nya untuk memudahkan implementation
