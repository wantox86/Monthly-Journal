package services

import (
	"fmt"
	"monthly-journal/internal/models"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ReportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

func (rs *ReportService) GenerateHTMLReport(month string) (string, int64, int, error) {
	var expenses []models.Expense
	var totalResult struct {
		Total int64
	}

	rs.db.Where("month_year = ?", month).Order("date DESC").Find(&expenses)
	rs.db.Model(&models.Expense{}).Where("month_year = ?", month).Select("COALESCE(SUM(amount), 0) as total").Scan(&totalResult)

	html := buildHTMLEmail(month, expenses, totalResult.Total)
	return html, totalResult.Total, len(expenses), nil
}

func buildHTMLEmail(month string, expenses []models.Expense, total int64) string {
	html := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body { font-family: Arial, sans-serif; background-color: #f5f5f5; }
		.container { max-width: 600px; margin: 20px auto; background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
		h1 { color: #333; }
		.summary { background-color: #f9f9f9; padding: 15px; border-radius: 5px; margin-bottom: 20px; }
		.summary p { margin: 5px 0; }
		table { width: 100%; border-collapse: collapse; margin-top: 20px; }
		th { background-color: #007bff; color: #fff; padding: 12px; text-align: left; font-weight: bold; }
		td { padding: 10px; border-bottom: 1px solid #ddd; }
		tr:hover { background-color: #f5f5f5; }
		.footer { text-align: center; color: #666; margin-top: 20px; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<h1>Monthly Expense Report</h1>
		<div class="summary">
			<p><strong>Month:</strong> ` + month + `</p>
			<p><strong>Total Expenses:</strong> Rp ` + formatCurrency(total) + `</p>
			<p><strong>Number of Expenses:</strong> ` + strconv.Itoa(len(expenses)) + `</p>
		</div>
		<table>
			<thead>
				<tr>
					<th>No</th>
					<th>Tanggal</th>
					<th>Deskripsi</th>
					<th>Amount (Rp)</th>
					<th>Pengirim</th>
				</tr>
			</thead>
			<tbody>
`

	for i, expense := range expenses {
		sender := expense.Sender
		if sender == "" {
			sender = "NULL"
		}

		html += fmt.Sprintf(`
				<tr>
					<td>%d</td>
					<td>%s</td>
					<td>%s</td>
					<td>%s</td>
					<td>%s</td>
				</tr>
`,
			i+1,
			expense.Date.Format("02-01-2006 15:04"),
			escape(expense.Description),
			formatCurrency(int64(expense.Amount)),
			escape(sender),
		)
	}

	html += `
			</tbody>
		</table>
		<div class="footer">
			<p>This is an automated report. Please do not reply to this email.</p>
			<p>Report generated on ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
		</div>
	</div>
</body>
</html>
`

	return html
}

func formatCurrency(amount int64) string {
	str := strconv.FormatInt(amount, 10)
	var result strings.Builder

	for i, ch := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(ch)
	}

	return result.String()
}

func escape(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	).Replace(s)
}
