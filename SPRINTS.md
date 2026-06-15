# Implementation Sprints - Monthly Journal Backend

## Overview
Breakdown implementation menjadi 5 sprint sequential. Setiap sprint adalah deliverable yang bisa dikerjakan & tested standalone.

---

## Sprint 0: Project Setup & Configuration
**Duration:** ~1-2 jam  
**Dependencies:** None  
**Goal:** Setup foundation project, config, dan prepare development environment

### Tasks
- [ ] Initialize Go project structure (`go mod init`)
- [ ] Create `.env` file dari `.env.example` dengan test values
- [ ] Setup package structure:
  - `cmd/server/main.go` - Entry point
  - `internal/config/` - Config loader
  - `internal/models/` - Data structures
  - `internal/handlers/` - HTTP handlers
  - `internal/services/` - Business logic
  - `internal/database/` - DB setup
- [ ] Create `go.mod` dengan dependencies:
  - `github.com/gin-gonic/gin`
  - `gorm.io/gorm`
  - `gorm.io/driver/mysql`
  - `github.com/joho/godotenv`
- [ ] Create basic `config.go` untuk load environment variables
- [ ] Create `main.go` dengan basic Gin server (just listen on :8080)
- [ ] Test: `go run cmd/server/main.go` bisa start tanpa error

### Deliverables
```
project-root/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   ├── handlers/
│   ├── services/
│   └── database/
├── .env (gitignore)
├── .env.example
├── go.mod
└── go.sum
```

---

## Sprint 1: Database Setup & Models
**Duration:** ~2-3 jam  
**Dependencies:** Sprint 0 ✅  
**Goal:** Database connection & define data models

### Tasks
- [ ] Create `internal/database/db.go` - Database connection setup
- [ ] Create database connection function dengan DSN dari config
- [ ] Test database connection dari main
- [ ] Create `internal/models/expense.go`:
  ```go
  type Expense struct {
    ID          int       `gorm:"primaryKey"`
    Date        time.Time
    Description string
    Amount      int
    Sender      string
    MonthYear   string
    CreatedAt   time.Time
  }
  ```
- [ ] Create `internal/models/models.go` - Register models
- [ ] Create database migration function (auto migration via GORM)
- [ ] Run migration & verify table created di database
- [ ] Test: Connect & verify expenses table exists

### Deliverables
- Database connection working
- Expenses table created in `monthly_bill` database
- Models defined & registered

### Verification
```bash
# Check table di MySQL
SELECT * FROM monthly_bill.expenses LIMIT 1;
```

---

## Sprint 2: Expense API - CRUD Operations
**Duration:** ~3-4 jam  
**Dependencies:** Sprint 1 ✅  
**Goal:** Implement full CRUD endpoints untuk expenses

### Tasks

#### 2.1 - POST /api/expenses (Create)
- [ ] Create `internal/handlers/expense.go` - Handler functions
- [ ] Implement `CreateExpense` handler:
  - Parse JSON request (description, amount, sender)
  - Validate: description & amount required
  - Auto-generate month_year dari current date (format: YYYY-MM)
  - Auto-set sender ke "Unknown" jika kosong
  - Save ke database
  - Return 201 + expense object
- [ ] Test dengan curl:
  ```bash
  curl -X POST http://localhost:8080/api/expenses \
    -H "Content-Type: application/json" \
    -d '{"description":"jajan kopi","amount":40000,"sender":"Nur Dahlia"}'
  ```

#### 2.2 - GET /api/expenses (List by Month)
- [ ] Implement `GetExpenses` handler:
  - Query parameter: `month` (optional, default: current month YYYY-MM)
  - Fetch expenses dari database WHERE month_year = month
  - Calculate total amount
  - Return expenses array + total + count
- [ ] Test dengan curl:
  ```bash
  curl http://localhost:8080/api/expenses?month=2026-06
  ```

#### 2.3 - DELETE /api/expenses/:id
- [ ] Implement `DeleteExpense` handler:
  - Parse id dari URL param
  - Delete dari database
  - Return status + id
- [ ] Test dengan curl:
  ```bash
  curl -X DELETE http://localhost:8080/api/expenses/1
  ```

#### 2.4 - Setup Routes
- [ ] Update `main.go` / create `internal/routes/routes.go`:
  - Register POST /api/expenses → CreateExpense
  - Register GET /api/expenses → GetExpenses
  - Register DELETE /api/expenses/:id → DeleteExpense

### Deliverables
- 3 endpoints working (POST, GET, DELETE)
- Validation & error handling
- All 3 endpoints tested & verified

---

## Sprint 3: Email & Report Service
**Duration:** ~3-4 jam  
**Dependencies:** Sprint 2 ✅  
**Goal:** Implement SMTP service & report generation

### Tasks

#### 3.1 - SMTP Service
- [ ] Create `internal/services/email.go`:
  ```go
  type EmailService struct {
    SMTPHost   string
    SMTPPort   int
    SMTPUser   string
    SMTPPass   string
    EmailFrom  string
  }
  
  func NewEmailService(cfg *config.Config) *EmailService
  func (es *EmailService) SendEmail(to []string, subject string, body string) error
  ```
- [ ] Implement `SendEmail` function menggunakan net/smtp
- [ ] Test: Send test email ke configured recipients

#### 3.2 - Report HTML Generator
- [ ] Create `internal/services/report.go`:
  ```go
  type ReportService struct {
    db *gorm.DB
  }
  
  func (rs *ReportService) GenerateHTMLReport(month string) (string, error)
  ```
- [ ] Implement report generation logic:
  - Fetch expenses untuk month dari database
  - Build HTML table:
    - Headers: No | Tanggal | Deskripsi | Amount (Rp) | Pengirim
    - Format tanggal: DD-MM-YYYY HH:mm
    - Format amount dengan comma (82,350)
    - Handle NULL sender
  - Return HTML string
- [ ] Create HTML email template dengan header & footer
- [ ] Test: Generate report untuk current month

#### 3.3 - POST /api/reports/send Handler
- [ ] Create `internal/handlers/report.go`:
  ```go
  func SendReport(db *gorm.DB, emailService *EmailService, reportService *ReportService) gin.HandlerFunc
  ```
- [ ] Implement handler:
  - Get current month (YYYY-MM)
  - Generate HTML report
  - Build email with report as body
  - Send to configured recipients
  - Return response dengan status, month, total, count, recipients, sent_at
- [ ] Register route: POST /api/reports/send

### Deliverables
- SMTP service working
- Report HTML generator working
- /api/reports/send endpoint working
- Email sent successfully ke recipients

### Verification
```bash
curl -X POST http://localhost:8080/api/reports/send \
  -H "Content-Type: application/json" \
  -d '{"format":"csv"}'
```

---

## Sprint 4: Integration & Testing
**Duration:** ~2-3 jam  
**Dependencies:** Sprint 3 ✅  
**Goal:** Integration testing, error handling, & documentation

### Tasks

#### 4.1 - Error Handling & Validation
- [ ] Add comprehensive error handling:
  - Database errors
  - SMTP errors
  - Validation errors
  - Return proper HTTP status codes (400, 500, etc)
- [ ] Add logging untuk debugging

#### 4.2 - Integration Testing
- [ ] Create test flow:
  1. POST /api/expenses (create 3-5 test expenses)
  2. GET /api/expenses?month=YYYY-MM (verify all expenses)
  3. POST /api/reports/send (verify email sent)
  4. DELETE /api/expenses/:id (delete 1 expense)
  5. GET /api/expenses?month=YYYY-MM (verify deleted)
- [ ] Document testing steps

#### 4.3 - Documentation
- [ ] Update README dengan:
  - Setup instructions
  - Environment variables required
  - API endpoint list
  - Example requests/responses
  - Testing guide
- [ ] Add code comments untuk public functions

#### 4.4 - Cleanup
- [ ] Remove test/debug code
- [ ] Verify all endpoints production-ready
- [ ] Security check (no credentials in code, proper validation)

### Deliverables
- All endpoints tested & documented
- Error handling implemented
- README dengan setup & usage guide
- Project ready for Android integration

---

## Timeline Summary

| Sprint | Duration | Status |
|--------|----------|--------|
| Sprint 0: Setup | 1-2h | ⏳ Ready |
| Sprint 1: Database | 2-3h | ⏳ Ready |
| Sprint 2: CRUD APIs | 3-4h | ⏳ Ready |
| Sprint 3: Email & Report | 3-4h | ⏳ Ready |
| Sprint 4: Testing & Docs | 2-3h | ⏳ Ready |
| **Total** | **11-16h** | |

---

## Testing Checklist Template

### Sprint 0 ✅
- [ ] `go run cmd/server/main.go` bisa start
- [ ] Server listen di :8080
- [ ] No build errors

### Sprint 1 ✅
- [ ] Database connection successful
- [ ] Expenses table created
- [ ] `SHOW TABLES FROM monthly_bill;` shows expenses table

### Sprint 2 ✅
- [ ] POST /api/expenses returns 201 + expense object
- [ ] GET /api/expenses returns expenses array
- [ ] DELETE /api/expenses/:id returns status deleted
- [ ] Month auto-populated (YYYY-MM format)
- [ ] Sender defaults to "Unknown" if empty

### Sprint 3 ✅
- [ ] SMTP connection successful
- [ ] Test email received in recipients inbox
- [ ] HTML table format correct (No|Tanggal|Deskripsi|Amount|Pengirim)
- [ ] Amount formatted dengan comma (82,350)
- [ ] POST /api/reports/send returns 200 + report data

### Sprint 4 ✅
- [ ] Error handling tested (invalid input, db error, smtp error)
- [ ] README complete & accurate
- [ ] All endpoints documented
- [ ] Ready untuk Android integration

---

## Notes
- Setiap sprint bisa dikerjakan independent dari yang sebelumnya (dengan dependency dipenuhi)
- Test setiap sprint sebelum lanjut ke sprint berikutnya
- Update `.env` dengan real values saat setup
- Untuk SMTP Gmail: gunakan [App Password](https://support.google.com/accounts/answer/185833), bukan password account
