# Docker Setup & Deployment

Backend ini sudah **Docker Ready** dengan multi-stage build untuk optimalisasi ukuran image.

## Prerequisites

- Docker & Docker Compose installed
- Access ke database di 192.168.50.131:3306
- SMTP credentials configured

## Quick Start - Local

### 1. Build Docker Image

```bash
docker build -t monthly-journal:latest .
```

**Output:** Image size ~15-20MB (menggunakan Alpine Linux)

### 2. Run Container

**Dengan environment variables inline:**
```bash
docker run -d \
  --name monthly-journal-api \
  -p 8080:8080 \
  -e DB_HOST=192.168.50.131 \
  -e DB_PORT=3306 \
  -e DB_NAME=monthly_bill \
  -e DB_USER=copilot \
  -e DB_PASS=copilot123 \
  -e SMTP_HOST=smtp.gmail.com \
  -e SMTP_PORT=587 \
  -e SMTP_USER=waonex86@gmail.com \
  -e SMTP_PASS="rycs xeoy laly wenh" \
  -e EMAIL_FROM=waonex86@gmail.com \
  -e EMAIL_RECIPIENTS=nurdahliana86@gmail.com,waonex86@gmail.com \
  -e SERVER_PORT=8080 \
  -e SERVER_ENV=production \
  monthly-journal:latest
```

**Atau gunakan .env file:**
```bash
docker run -d \
  --name monthly-journal-api \
  -p 8080:8080 \
  --env-file .env \
  monthly-journal:latest
```

### 3. Verify Container

```bash
# Check container running
docker ps | grep monthly-journal

# Check logs
docker logs monthly-journal-api

# Test health
curl http://localhost:8080/health
```

## Using Docker Compose (Recommended)

### 1. Start Services

```bash
docker-compose up -d
```

**Apa yang terjadi:**
- Build image dari Dockerfile
- Start container `monthly-journal-api` di port 8080
- Health check configured
- Network `monthly-journal-net` created

### 2. Check Status

```bash
docker-compose ps
docker-compose logs -f monthly-journal-api
```

### 3. Test API

```bash
curl http://localhost:8080/health

# Test create expense
curl -X POST http://localhost:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "description": "jajan kopi",
    "amount": 40000,
    "sender": "Nur Dahlia"
  }'
```

### 4. Stop Services

```bash
docker-compose down

# Stop dan remove volumes
docker-compose down -v
```

## Deployment ke 192.168.50.131

### Option 1: Push Image to Remote Docker

**Step 1: Tag image**
```bash
docker tag monthly-journal:latest 192.168.50.131:5000/monthly-journal:latest
```

**Step 2: Push ke remote registry (jika ada)**
```bash
docker push 192.168.50.131:5000/monthly-journal:latest
```

**Step 3: SSH ke server dan run**
```bash
ssh user@192.168.50.131

# Pull image
docker pull 192.168.50.131:5000/monthly-journal:latest

# Run container
docker run -d \
  --name monthly-journal-api \
  -p 8080:8080 \
  --env-file .env \
  192.168.50.131:5000/monthly-journal:latest

# Check status
docker ps
docker logs monthly-journal-api
```

### Option 2: Copy Dockerfile & Build di Server

**Step 1: Copy files ke server**
```bash
scp Dockerfile docker-compose.yml .dockerignore .env user@192.168.50.131:/home/user/monthly-journal/
scp -r cmd internal user@192.168.50.131:/home/user/monthly-journal/
scp go.mod go.sum user@192.168.50.131:/home/user/monthly-journal/
```

**Step 2: SSH ke server**
```bash
ssh user@192.168.50.131
cd monthly-journal
```

**Step 3: Build & Run**
```bash
# Build
docker build -t monthly-journal:latest .

# Run dengan docker-compose
docker-compose up -d

# Atau docker run langsung
docker run -d \
  --name monthly-journal-api \
  -p 8080:8080 \
  --env-file .env \
  monthly-journal:latest
```

**Step 4: Verify**
```bash
docker ps
docker logs monthly-journal-api

# Test dari server
curl http://localhost:8080/health
```

### Option 3: SSH + Docker Compose

**Paling simple - copy dan run:**

```bash
# Copy ke server
scp -r . user@192.168.50.131:/home/user/monthly-journal/

# SSH
ssh user@192.168.50.131

# Deploy
cd /home/user/monthly-journal
docker-compose up -d

# Check
docker-compose ps
curl http://localhost:8080/health
```

## Docker Images

### Build Details

**Dockerfile menggunakan multi-stage:**
1. **Builder stage (golang:1.17-alpine)**
   - Download dependencies
   - Compile Go binary
   - Result: compiled executable

2. **Runtime stage (alpine:latest)**
   - Copy binary dari builder
   - Minimal dependencies hanya untuk runtime
   - **Total size: ~15-20MB**

### Size Comparison

```
monthly-journal:latest  ~15-20MB
vs
golang:1.17-full       ~800MB+
```

## Environment Variables

Semua variables dari `.env` bisa di-override:

```bash
docker run -e DB_HOST=192.168.50.131 monthly-journal:latest
```

Atau gunakan `.env` file:
```bash
docker run --env-file .env monthly-journal:latest
```

## Health Check

Container punya health check built-in:

```bash
docker ps --format "table {{.Names}}\t{{.Status}}"

# Output:
# NAME                    STATUS
# monthly-journal-api     Up 2 hours (healthy)
```

## Logs & Monitoring

```bash
# Real-time logs
docker logs -f monthly-journal-api

# Last 100 lines
docker logs --tail 100 monthly-journal-api

# With timestamps
docker logs -t monthly-journal-api
```

## Troubleshooting

### Container exits immediately

```bash
docker logs monthly-journal-api
# Check: Database connection error? SMTP credentials?
```

### Port already in use

```bash
# Change port mapping
docker run -p 9090:8080 monthly-journal:latest
```

### Database connection failed

```
Make sure 192.168.50.131:3306 is accessible
Check .env credentials (user: copilot, pass: copilot123)
```

### SMTP error

```
Verify SMTP credentials in .env
Check firewall allows outbound SMTP (port 587)
Test SMTP manually if needed
```

## Production Checklist

- [ ] Database credentials setup di .env
- [ ] SMTP credentials setup di .env  
- [ ] Port 8080 accessible / firewall configured
- [ ] Health check responding
- [ ] Container restart policy: `unless-stopped`
- [ ] Logs monitored / collected
- [ ] Backup database regularly

## Cleanup

```bash
# Stop container
docker stop monthly-journal-api

# Remove container
docker rm monthly-journal-api

# Remove image
docker rmi monthly-journal:latest

# Remove everything (compose)
docker-compose down -v
```

## Notes

- ✅ Multi-stage build untuk minimal image size
- ✅ Health check configured
- ✅ Auto restart on crash
- ✅ Alpine Linux untuk security & size
- ✅ No root user inside container
- ✅ Production-ready

## References

- Docker: https://docs.docker.com/
- Docker Compose: https://docs.docker.com/compose/
- Alpine Linux: https://alpinelinux.org/
