# Monthly-Journal - Expense Tracker

Multi-platform aplikasi untuk catat belanja harian dan generate laporan bulanan.

## Project Structure

- **Backend:** Go REST API (Gin + MySQL)
- **Android:** React Native app (future)
- **Documentation:** API specs, sprint breakdown, architecture

## Backend Setup

See [Backend Documentation](#backend-documentation) below or follow [SPRINTS.md](SPRINTS.md) for implementation details.

---

# Backend Documentation

REST API untuk expense tracking dengan laporan email bulanan.

## Quick Start

```bash
# Setup environment
cp .env.example .env
# Edit .env dengan database & SMTP credentials

# Run server
go run cmd/server/main.go
```

Server berjalan di `http://localhost:8080`

## API Endpoints

### Health Check
```bash
GET /health
```

### Expenses
```bash
# Create
POST /api/expenses
{
  "description": "jajan kopi",
  "amount": 40000,
  "sender": "Nur Dahlia"
}

# List by month
GET /api/expenses?month=2026-06

# Delete
DELETE /api/expenses/:id
```

### Reports
```bash
# Send monthly report via email
POST /api/reports/send
{
  "format": "csv"
}
```

## Configuration (.env)

```env
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

## Docker Deployment

Aplikasi sudah **Docker Ready** dengan multi-stage build (~15-20MB image size).

### Quick Deploy

```bash
# Local
docker build -t monthly-journal:latest .
docker-compose up -d

# Or single container
docker run -d -p 8080:8080 --env-file .env monthly-journal:latest
```

### Deploy ke 192.168.50.131

```bash
# Option 1: Copy & build di server
scp -r . user@192.168.50.131:/home/user/monthly-journal/
ssh user@192.168.50.131
cd /home/user/monthly-journal
docker-compose up -d

# Option 2: Push image ke registry
docker tag monthly-journal:latest 192.168.50.131:5000/monthly-journal:latest
docker push 192.168.50.131:5000/monthly-journal:latest
```

Lihat [DOCKER.md](DOCKER.md) untuk deployment details.

## Documentation Files

- **[SPRINTS.md](SPRINTS.md)** - 5 sprint breakdown dengan checklist
- **[PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)** - Directory structure & code templates
- **[DOCKER.md](DOCKER.md)** - Docker setup & deployment guide
- **[claude.md](claude.md)** - Project specifications
- **[.env.example](.env.example)** - Environment variables template
