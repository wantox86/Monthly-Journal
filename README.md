# Monthly-Journal - Expense Tracker

Multi-platform aplikasi untuk catat belanja harian dan generate laporan bulanan.

## Project Structure

- **Backend:** Go REST API (Gin + MySQL)
- **Android:** React Native app (future)
- **Documentation:** API specs, sprint breakdown, architecture

## Quick Start (Docker)

```bash
# 1. Start containers (API + MySQL)
docker-compose up -d

# 2. Test health
curl http://localhost:8080/health

# 3. Stop containers
docker-compose down
```

See [Backend Documentation](#backend-documentation) below or follow [SPRINTS.md](SPRINTS.md) for implementation details.

---

# Backend Documentation

REST API untuk expense tracking dengan laporan email bulanan.

## Quick Start

### Option 1: Docker (Recommended)

```bash
# Start all services (API + MySQL)
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f monthly-journal-api

# Stop
docker-compose down
```

### Option 2: Local Development

```bash
# Setup environment
cp .env.example .env

# Update .env dengan database & SMTP credentials
# Make sure MySQL running on port 3306

# Run server
go run cmd/server/main.go
```

Server berjalan di `http://localhost:8080`

## API Endpoints

All endpoints running on `http://localhost:8080`

### Health Check
```bash
GET /health

# Response
{
  "status": "ok"
}
```

### Expenses

**Create Expense**
```bash
POST /api/expenses
Content-Type: application/json

{
  "description": "jajan kopi",
  "amount": 40000,
  "sender": "Nur Dahlia"
}

# Response (201)
{
  "id": 1,
  "date": "2026-06-15T15:06:26.920760003Z",
  "description": "jajan kopi",
  "amount": 40000,
  "sender": "Nur Dahlia",
  "month_year": "2026-06",
  "created_at": "2026-06-15T15:06:26.921Z"
}
```

**List Expenses**
```bash
GET /api/expenses?month=2026-06

# Response (200)
{
  "month": "2026-06",
  "total": 40000,
  "count": 1,
  "expenses": [
    {
      "id": 1,
      "date": "2026-06-15T15:06:26.921Z",
      "description": "jajan kopi",
      "amount": 40000,
      "sender": "Nur Dahlia"
    }
  ]
}
```

**Delete Expense**
```bash
DELETE /api/expenses/:id

# Response (200)
{
  "status": "deleted",
  "id": 1
}
```

### Reports

**Send Monthly Report**
```bash
POST /api/reports/send
Content-Type: application/json

{
  "format": "csv"
}

# Response (200)
{
  "status": "sent",
  "month": "2026-06",
  "total": 40000,
  "count": 1,
  "recipients": [
    "nurdahliana86@gmail.com",
    "waonex86@gmail.com"
  ],
  "sent_at": "2026-06-15T15:06:26.921Z"
}
```

## Configuration

### Environment Variables (.env)

```env
# Database (Docker container: mysql)
DB_HOST=mysql
DB_PORT=3306
DB_NAME=monthly_bill
DB_USER=copilot
DB_PASS=copilot123

# Email / SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=waonex86@gmail.com
SMTP_PASS=rycs xeoy laly wenh

EMAIL_FROM=waonex86@gmail.com
EMAIL_RECIPIENTS=nurdahliana86@gmail.com,waonex86@gmail.com

# Server
SERVER_PORT=8080
SERVER_ENV=production
```

**Note:** `.env` file sudah ada dengan konfigurasi yang benar. Untuk local development, ubah `DB_HOST=localhost`

## Docker Deployment

Aplikasi sudah **Docker Ready** dengan multi-stage build:
- **Image size:** ~20MB (Alpine-based)
- **Services:** API (Go) + MySQL 8.0
- **Health checks:** Built-in for both
- **Auto-restart:** Enabled for production

### Local Development

```bash
# Start all services
docker-compose up -d

# Services running:
# - API: http://localhost:8080
# - MySQL: localhost:3306

# Check status
docker-compose ps

# View logs
docker-compose logs -f monthly-journal-api

# Stop all
docker-compose down
```

### Database Access

```bash
# Via mysql client
mysql -h 127.0.0.1 -u copilot -p copilot123 monthly_bill

# Via docker exec
docker exec -it monthly-journal-mysql mysql -u copilot -p copilot123 monthly_bill
```

### Deploy ke Server

```bash
# Copy to server
scp -r . user@192.168.50.131:/home/user/monthly-journal/

# SSH to server
ssh user@192.168.50.131
cd /home/user/monthly-journal

# Start
docker-compose up -d

# Check status
docker-compose ps
curl http://localhost:8080/health
```

Lihat [DOCKER.md](DOCKER.md) untuk deployment details lengkap.

## Testing

### Quick Test

```bash
# All services running?
docker-compose ps

# Health check
curl http://localhost:8080/health

# Create expense
curl -X POST http://localhost:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{"description":"test","amount":10000,"sender":"Test User"}'

# Get expenses
curl http://localhost:8080/api/expenses
```

### Database Testing

```bash
# Connect to database
mysql -h 127.0.0.1 -u copilot -p copilot123 monthly_bill

# Show tables
SHOW TABLES;

# Check expenses
SELECT * FROM expenses;
```

## Troubleshooting

**Container fails to start?**
```bash
# Check logs
docker-compose logs monthly-journal-api

# Restart
docker-compose restart monthly-journal-api
```

**Port already in use?**
```bash
# Change port in docker-compose.yml
# ports:
#   - "9090:8080"  # Use 9090 instead
```

**Database connection error?**
```bash
# Make sure MySQL container is healthy
docker-compose logs monthly-journal-mysql

# Check port 3306
netstat -an | grep 3306
```

**Reset everything?**
```bash
# Remove all containers and volumes
docker-compose down -v

# Rebuild and restart
docker-compose up -d
```

---

## Documentation Files

- **[SPRINTS.md](SPRINTS.md)** - 5 sprint breakdown dengan checklist
- **[PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)** - Directory structure & code templates
- **[DOCKER.md](DOCKER.md)** - Detailed Docker setup & deployment guide
- **[CLAUDE.md](CLAUDE.md)** - Project specifications & implementation notes
- **[.env.example](.env.example)** - Environment variables template
