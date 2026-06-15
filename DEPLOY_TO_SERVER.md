# Deployment ke 192.168.50.131

Step-by-step guide untuk deploy Monthly Journal Backend ke server 192.168.50.131.

## Prerequisites

- SSH access ke 192.168.50.131
- Docker & Docker Compose installed di server
- Database sudah running di 192.168.50.131:3306 (monthly_bill)

## Method 1: Copy Files & Build di Server (Recommended)

**Paling simple dan reliable untuk first-time deployment.**

### Step 1: Copy Project Files ke Server

```bash
# Dari local machine
scp -r . user@192.168.50.131:/home/user/monthly-journal/
```

### Step 2: SSH ke Server

```bash
ssh user@192.168.50.131
```

### Step 3: Navigate & Deploy

```bash
cd /home/user/monthly-journal

# Verify files exist
ls -la

# Start with docker-compose
docker-compose up -d
```

### Step 4: Verify Deployment

```bash
# Check container status
docker-compose ps
docker logs -f monthly-journal-api

# Test API (dari server)
curl http://localhost:8080/health

# Test dari local (replace IP dengan server IP)
curl http://192.168.50.131:8080/health
```

**Expected Output:**
```json
{"status":"ok"}
```

### Step 5: Test Full Flow

```bash
# Create expense
curl -X POST http://192.168.50.131:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "description": "test belanja",
    "amount": 50000,
    "sender": "Test User"
  }'

# List expenses
curl http://192.168.50.131:8080/api/expenses

# Send report (akan kirim email jika SMTP configured)
curl -X POST http://192.168.50.131:8080/api/reports/send \
  -H "Content-Type: application/json" \
  -d '{"format":"csv"}'
```

---

## Method 2: Manual Docker Commands

Jika tidak ingin pakai docker-compose:

### Step 1-2: Copy & SSH (sama seperti Method 1)

### Step 3: Build Image

```bash
cd /home/user/monthly-journal
docker build -t monthly-journal:latest .

# Check image
docker images | grep monthly-journal
```

### Step 4: Run Container

```bash
docker run -d \
  --name monthly-journal-api \
  -p 8080:8080 \
  --env-file .env \
  --restart unless-stopped \
  monthly-journal:latest
```

### Step 5: Verify

```bash
docker ps
docker logs monthly-journal-api
curl http://localhost:8080/health
```

---

## Method 3: Using Docker Registry (For CI/CD)

Jika ingin setup automated deployment:

### Step 1: Tag & Push Image

```bash
# Local machine
docker build -t monthly-journal:latest .
docker tag monthly-journal:latest 192.168.50.131:5000/monthly-journal:latest
docker push 192.168.50.131:5000/monthly-journal:latest
```

### Step 2: Pull & Run di Server

```bash
# Di server
docker pull 192.168.50.131:5000/monthly-journal:latest
docker run -d \
  --name monthly-journal-api \
  -p 8080:8080 \
  --env-file .env \
  192.168.50.131:5000/monthly-journal:latest
```

---

## Configuration di Server

### .env File

Pastikan `.env` sudah di-setup dengan credentials yang benar:

```bash
# Di server, edit .env
nano /home/user/monthly-journal/.env
```

Verify credentials:
```env
DB_HOST=192.168.50.131
DB_PORT=3306
DB_NAME=monthly_bill
DB_USER=copilot
DB_PASS=copilot123

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=waonex86@gmail.com
SMTP_PASS=rycs xeoy laly wenh

EMAIL_FROM=waonex86@gmail.com
EMAIL_RECIPIENTS=nurdahliana86@gmail.com,waonex86@gmail.com

SERVER_PORT=8080
SERVER_ENV=production
```

---

## Troubleshooting

### Container not starting

```bash
docker logs monthly-journal-api

# Common issues:
# 1. Database not accessible
#    -> Verify DB_HOST, DB_PORT, credentials
#    -> mysql -h 192.168.50.131 -u copilot -p

# 2. Port already in use
#    -> docker ps (check what using port 8080)
#    -> Change port: -p 9090:8080

# 3. SMTP error
#    -> Check SMTP credentials in .env
#    -> Verify outbound SMTP access (port 587)
```

### Database connection error

```bash
# From server, test MySQL connection
mysql -h 192.168.50.131 -u copilot -p copilot123 -e "USE monthly_bill; SHOW TABLES;"

# Should show: expenses table
```

### Port not accessible

```bash
# Check firewall di server
sudo ufw status
sudo ufw allow 8080

# Test dari local
curl http://192.168.50.131:8080/health
```

### Email not sending

```bash
# Check logs
docker logs monthly-journal-api | grep -i "smtp\|email"

# Verify SMTP credentials
# Test manually jika needed
```

---

## Maintenance

### View Logs

```bash
# Real-time
docker logs -f monthly-journal-api

# Last 100 lines
docker logs --tail 100 monthly-journal-api
```

### Restart Container

```bash
docker-compose restart monthly-journal-api
# atau
docker restart monthly-journal-api
```

### Update Code

```bash
# Pull latest changes
cd /home/user/monthly-journal
git pull origin main

# Rebuild image
docker-compose down
docker build -t monthly-journal:latest .
docker-compose up -d
```

### Check Container Status

```bash
docker-compose ps
docker ps
docker stats monthly-journal-api
```

---

## Monitoring

### Health Check

```bash
# Container health
docker ps --format "table {{.Names}}\t{{.Status}}"

# API health
curl http://192.168.50.131:8080/health
```

### Logs Monitoring

```bash
# Follow logs
docker logs -f monthly-journal-api

# Search for errors
docker logs monthly-journal-api | grep -i error
```

### Database Verification

```bash
# Check expenses table
mysql -h 192.168.50.131 -u copilot -p copilot123 monthly_bill \
  -e "SELECT COUNT(*) as 'Total Expenses' FROM expenses;"
```

---

## Backup & Cleanup

### Backup Database

```bash
# From server
mysqldump -h 192.168.50.131 -u copilot -p copilot123 monthly_bill > backup_$(date +%Y%m%d).sql

# From local
mysql -h 192.168.50.131 -u copilot -p copilot123 monthly_bill < backup.sql
```

### Cleanup Old Containers

```bash
# Remove stopped containers
docker container prune

# Remove unused images
docker image prune

# Full cleanup
docker system prune -a
```

---

## Rollback (if needed)

```bash
# Stop current
docker-compose down

# Checkout previous version
git log --oneline
git checkout <previous-commit>

# Rebuild & restart
docker-compose up -d
```

---

## Production Checklist

- [ ] SSH access ke 192.168.50.131 working
- [ ] Docker & Docker Compose installed
- [ ] .env file configured dengan credentials
- [ ] Database test: `mysql -h 192.168.50.131 -u copilot -p copilot123`
- [ ] Container running: `docker ps`
- [ ] Health check: `curl http://192.168.50.131:8080/health`
- [ ] Database accessible from container
- [ ] SMTP credentials working
- [ ] Port 8080 accessible from local
- [ ] Logs monitored
- [ ] Backup strategy in place

---

## Quick Reference Commands

```bash
# Deploy
docker-compose up -d

# Status
docker-compose ps

# Logs
docker logs -f monthly-journal-api

# Restart
docker-compose restart

# Stop
docker-compose stop

# Remove
docker-compose down

# Rebuild
docker-compose up -d --build

# Test API
curl http://192.168.50.131:8080/health
curl http://192.168.50.131:8080/api/expenses
```

---

**Notes:**
- Semua files sudah Docker-ready
- Tidak perlu Go compiler di server (image include binary)
- Multi-stage build = minimal image size (~15-20MB)
- Auto-restart on crash (unless-stopped policy)
- Health check configured

**Waktu deployment:** ~5-10 menit (dari copy files sampai running)
