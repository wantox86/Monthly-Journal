# ⚡ Quick Start - Deploy to 192.168.50.131

**Fastest way to deploy. Just follow 3 steps.**

## 3-Step Deployment

### Step 1: Copy Project (5 min)
```bash
scp -r . user@192.168.50.131:/home/user/monthly-journal/
```

### Step 2: SSH & Deploy (2 min)
```bash
ssh user@192.168.50.131
cd /home/user/monthly-journal
docker-compose up -d
```

### Step 3: Test (1 min)
```bash
# From server
curl http://localhost:8080/health

# Or from local
curl http://192.168.50.131:8080/health
```

**✅ DONE! Service is running**

---

## What's Running?

- **API Server:** http://192.168.50.131:8080
- **Database:** 192.168.50.131:3306 (monthly_bill)
- **Email:** SMTP configured for reports
- **Status:** Auto-restart enabled

---

## Test the API

```bash
# Create expense
curl -X POST http://192.168.50.131:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{"description":"Kopi Pagi","amount":35000,"sender":"Wawan"}'

# List expenses
curl http://192.168.50.131:8080/api/expenses

# Send monthly report
curl -X POST http://192.168.50.131:8080/api/reports/send \
  -H "Content-Type: application/json" \
  -d '{"format":"csv"}'
```

---

## If Something Goes Wrong

```bash
# Check logs
docker logs -f monthly-journal-api

# Check container status
docker ps

# Restart
docker-compose restart

# Stop & remove
docker-compose down
```

---

## Configuration

All configuration is in `.env` file:
```env
DB_HOST=192.168.50.131           # Database
SMTP_USER=waonex86@gmail.com     # Email sender
EMAIL_RECIPIENTS=...             # Email recipients
SERVER_ENV=production            # Environment
```

**No changes needed - everything is pre-configured!**

---

## Full Documentation

- **[DOCKER.md](DOCKER.md)** - Docker details
- **[DEPLOY_TO_SERVER.md](DEPLOY_TO_SERVER.md)** - Detailed guide
- **[README.md](README.md)** - API documentation
- **[SPRINTS.md](SPRINTS.md)** - Implementation details

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Container won't start | `docker logs monthly-journal-api` |
| Connection refused | Verify IP/Port in .env |
| Database error | Check MySQL is running |
| Email not sending | Verify SMTP credentials |
| Port 8080 in use | Change port in docker-compose.yml |

---

**Status: PRODUCTION-READY ✅**

Estimated setup time: **~10 minutes**
