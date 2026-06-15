# Expense Tracker App - Android + Backend

## Overview
Aplikasi untuk catat belanja harian + kirim report bulanan. Features:
- **Catat Belanja** - Input expense dengan deskripsi & amount
- **Kirim Report** - Generate & send laporan bulanan (CSV/PDF/JSON)

---

## Backend API (Go + Gin)

### Database Schema

```sql
CREATE TABLE expenses (
  id INT PRIMARY KEY AUTO_INCREMENT,
  date DATETIME DEFAULT CURRENT_TIMESTAMP,
  description VARCHAR(255) NOT NULL,
  amount INT NOT NULL,
  sender VARCHAR(100),
  month_year VARCHAR(7),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Notes:**
- Database: `monthly_bill` di `192.168.50.131:3306`
- User: `copilot` | Pass: `copilot123`
- Table `reports` tidak digunakan (report generated on-the-fly)

### API Endpoints

#### 1. **POST /api/expenses**
Catat belanja baru.

**Request:**
```json
{
  "description": "jajan kopi",
  "amount": 40000,
  "sender": "Nur Dahlia"
}
```

**Response (201):**
```json
{
  "id": 7,
  "date": "2026-06-15T12:23:58Z",
  "description": "jajan kopi",
  "amount": 40000,
  "sender": "Nur Dahlia",
  "month_year": "2026-06"
}
```

#### 2. **GET /api/expenses?month=2026-06**
Ambil list belanja untuk bulan tertentu.

**Response (200):**
```json
{
  "month": "2026-06",
  "total": 532850,
  "count": 9,
  "expenses": [
    {
      "id": 1,
      "date": "2026-06-15T00:05:17Z",
      "description": "astro",
      "amount": 50000,
      "sender": "Unknown"
    },
    ...
  ]
}
```

#### 3. **POST /api/reports/send**
Generate report untuk bulan berjalan & kirim via email.

**Request:**
```json
{
  "format": "csv"
}
```

**Response (200):**
```json
{
  "status": "sent",
  "month": "2026-06",
  "total": 532850,
  "count": 9,
  "recipients": [
    "nurdahliana86@gmail.com",
    "waonex86@gmail.com"
  ],
  "sent_at": "2026-06-15T15:17:29Z"
}
```

**Email Format:**
- Sender: `waonex86@gmail.com`
- Recipients: `nurdahliana86@gmail.com`, `waonex86@gmail.com`
- Body: HTML table dengan kolom `No | Tanggal | Deskripsi | Amount (Rp) | Pengirim`
- Example row: `1 | 15-06-2026 00:14 | jajan jus antara | 82,350 | NULL`

#### 5. **DELETE /api/expenses/:id**
Hapus expense tertentu.

**Response (200):**
```json
{
  "status": "deleted",
  "id": 7
}
```

---

## Android App (React Native)

### Tech Stack
- **Framework:** React Native / Expo
- **State:** Redux atau Context API
- **UI:** React Native Paper / NativeBase
- **HTTP:** Axios / Fetch API
- **Local Storage:** AsyncStorage

### Screen 1: Catat Belanja

```jsx
import React, { useState } from 'react';
import { View, TextInput, TouchableOpacity, Alert } from 'react-native';
import axios from 'axios';

export default function CatatBelanja() {
  const [description, setDescription] = useState('');
  const [amount, setAmount] = useState('');
  const [sender, setSender] = useState('Nur Dahlia');
  const [loading, setLoading] = useState(false);

  const handleSave = async () => {
    if (!description || !amount) {
      Alert.alert('Error', 'Deskripsi dan amount harus diisi');
      return;
    }

    setLoading(true);
    try {
      const response = await axios.post(
        'https://api.example.com/api/expenses',
        {
          description,
          amount: parseInt(amount),
          sender,
        }
      );
      Alert.alert('Sukses', `${description} Rp${amount} tercatat`);
      setDescription('');
      setAmount('');
    } catch (error) {
      Alert.alert('Error', error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={{ padding: 20 }}>
      <TextInput
        placeholder="Deskripsi (misal: jajan kopi)"
        value={description}
        onChangeText={setDescription}
        style={{ borderBottomWidth: 1, marginBottom: 15, paddingVertical: 10 }}
      />
      <TextInput
        placeholder="Amount (misal: 40000)"
        value={amount}
        onChangeText={setAmount}
        keyboardType="numeric"
        style={{ borderBottomWidth: 1, marginBottom: 15, paddingVertical: 10 }}
      />
      <TextInput
        placeholder="Pengirim"
        value={sender}
        onChangeText={setSender}
        style={{ borderBottomWidth: 1, marginBottom: 20, paddingVertical: 10 }}
      />
      <TouchableOpacity
        onPress={handleSave}
        disabled={loading}
        style={{
          backgroundColor: '#007AFF',
          padding: 15,
          borderRadius: 8,
          alignItems: 'center',
        }}
      >
        <Text style={{ color: '#fff', fontSize: 16, fontWeight: 'bold' }}>
          {loading ? 'Menyimpan...' : 'Catat Belanja'}
        </Text>
      </TouchableOpacity>
    </View>
  );
}
```

### Screen 2: Laporan Belanja

```jsx
import React, { useState, useEffect } from 'react';
import { View, FlatList, Text, TouchableOpacity, Alert } from 'react-native';
import axios from 'axios';

export default function LaporanBelanja() {
  const [expenses, setExpenses] = useState([]);
  const [total, setTotal] = useState(0);
  const [month, setMonth] = useState('2026-06');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchExpenses();
  }, [month]);

  const fetchExpenses = async () => {
    setLoading(true);
    try {
      const response = await axios.get(
        `https://api.example.com/api/expenses?month=${month}`
      );
      setExpenses(response.data.expenses);
      setTotal(response.data.total);
    } catch (error) {
      Alert.alert('Error', error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleSendReport = async (format = 'csv') => {
    try {
      await axios.post('https://api.example.com/api/reports/send', {
        format,
      });
      Alert.alert('Sukses', `Report ${format.toUpperCase()} dikirim ke email`);
    } catch (error) {
      Alert.alert('Error', error.message);
    }
  };

  return (
    <View style={{ padding: 20, flex: 1 }}>
      <Text style={{ fontSize: 18, fontWeight: 'bold', marginBottom: 10 }}>
        Laporan {month}
      </Text>
      <Text style={{ fontSize: 16, marginBottom: 15 }}>
        Total: Rp {total.toLocaleString('id-ID')}
      </Text>

      <FlatList
        data={expenses}
        keyExtractor={(item) => item.id.toString()}
        renderItem={({ item }) => (
          <View
            style={{
              padding: 10,
              borderBottomWidth: 1,
              borderBottomColor: '#eee',
            }}
          >
            <Text style={{ fontWeight: 'bold' }}>{item.description}</Text>
            <Text style={{ color: '#666' }}>
              Rp {item.amount.toLocaleString('id-ID')} • {item.sender}
            </Text>
          </View>
        )}
        ListEmptyComponent={<Text>Tidak ada data</Text>}
      />

      <View style={{ marginTop: 20, gap: 10 }}>
        <TouchableOpacity
          onPress={() => handleSendReport('csv')}
          style={{
            backgroundColor: '#34C759',
            padding: 12,
            borderRadius: 8,
            alignItems: 'center',
          }}
        >
          <Text style={{ color: '#fff', fontWeight: 'bold' }}>
            Kirim Report (CSV)
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          onPress={() => handleSendReport('pdf')}
          style={{
            backgroundColor: '#FF9500',
            padding: 12,
            borderRadius: 8,
            alignItems: 'center',
          }}
        >
          <Text style={{ color: '#fff', fontWeight: 'bold' }}>
            Kirim Report (PDF)
          </Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}
```

---

## Backend Implementation (Go + Gin)

### Main Routes

```go
package main

import (
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
)

func setupRoutes(db *gorm.DB) *gin.Engine {
  r := gin.Default()

  // Expenses
  r.POST("/api/expenses", createExpense(db))
  r.GET("/api/expenses", getExpenses(db))
  r.DELETE("/api/expenses/:id", deleteExpense(db))

  // Reports
  r.POST("/api/reports/generate", generateReport(db))
  r.POST("/api/reports/send", sendReport(db))
  r.GET("/api/reports/:month", getReport(db))

  return r
}

func createExpense(db *gorm.DB) gin.HandlerFunc {
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

    expense := Expense{
      Description: req.Description,
      Amount:      req.Amount,
      Sender:      req.Sender,
      MonthYear:   time.Now().Format("2006-01"),
    }

    if err := db.Create(&expense).Error; err != nil {
      c.JSON(500, gin.H{"error": err.Error()})
      return
    }

    c.JSON(201, expense)
  }
}

func getExpenses(db *gorm.DB) gin.HandlerFunc {
  return func(c *gin.Context) {
    month := c.DefaultQuery("month", time.Now().Format("2006-01"))
    var expenses []Expense
    var total int64

    db.Where("month_year = ?", month).Find(&expenses)
    db.Model(&Expense{}).Where("month_year = ?", month).Sum("amount", &total)

    c.JSON(200, gin.H{
      "month":    month,
      "total":    total,
      "count":    len(expenses),
      "expenses": expenses,
    })
  }
}
```

---

## Implementation Notes

### Report Generation
- Report generated on-the-fly (tidak ada file storage)
- Dikirim langsung ke email dalam format HTML table
- Menggunakan SMTP untuk email delivery
- Bulan otomatis sesuai bulan berjalan (current month)

### Future Enhancements
1. **Authentication** - JWT/Bearer token
2. **Multi-user** - Role-based access (admin, viewer)
3. **Categories** - Group expenses by type
4. **Budget tracking** - Set limits per category
5. **Export formats** - PDF, Excel, JSON
6. **Notifications** - Alert when expense added/report sent

---

## Configuration

### SMTP Setup (.env atau config file)
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password

EMAIL_FROM=waonex86@gmail.com
EMAIL_RECIPIENTS=nurdahliana86@gmail.com,waonex86@gmail.com

DB_HOST=192.168.50.131
DB_PORT=3306
DB_NAME=monthly_bill
DB_USER=copilot
DB_PASS=copilot123
```

**Notes:**
- Untuk Gmail: gunakan [App Password](https://support.google.com/accounts/answer/185833)
- Sesuaikan SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS dengan provider email yang digunakan

---

## Setup Instructions

### Backend
```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/joho/godotenv
go run main.go
# Server runs on :8080
```

### Android
```bash
npx expo init ExpenseTracker
cd ExpenseTracker
npm install axios react-native-paper
npm start
```

---

## Testing

### Catat Belanja
```bash
curl -X POST http://localhost:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "description": "jajan kopi",
    "amount": 40000,
    "sender": "Nur Dahlia"
  }'
```

### Ambil Report
```bash
curl http://localhost:8080/api/expenses?month=2026-06
```

### Kirim Report
```bash
curl -X POST http://localhost:8080/api/reports/send \
  -H "Content-Type: application/json" \
  -d '{
    "format": "csv"
  }'
```
Report akan generate untuk bulan berjalan dan dikirim ke recipients yang sudah configured di .env

---

## CSV Format (Current)
```
No,Tanggal,Deskripsi,Amount,Pengirim
1,2026-06-15 00:05:17,astro,50000,Unknown
...
```

✅ **Already compatible with API output**
