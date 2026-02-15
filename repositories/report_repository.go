package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	DB *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{DB: db}
}

func (r *ReportRepository) GetTodaySummary() (*models.SalesSummary, error) {
	summary := &models.SalesSummary{}

	// Total Revenue
	err := r.DB.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&summary.TotalRevenue)

	if err != nil {
		return nil, err
	}

	// Total Transaksi
	err = r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM transactions
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&summary.TotalTransaksi)

	if err != nil {
		return nil, err
	}

	// Produk Terlaris Hari Ini
	row := r.DB.QueryRow(`
		SELECT p.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`)

	bestSeller := &models.BestSeller{}
	err = row.Scan(&bestSeller.Nama, &bestSeller.QtyTerjual)

	if err == nil {
		summary.ProdukTerlaris = bestSeller
	}

	return summary, nil
}

func (r *ReportRepository) GetSummaryByDate(startDate, endDate string) (*models.SalesSummary, error) {
	summary := &models.SalesSummary{}

	// Total Revenue
	err := r.DB.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions
		WHERE created_at::date BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&summary.TotalRevenue)

	if err != nil {
		return nil, err
	}

	// Total Transaksi
	err = r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM transactions
		WHERE created_at::date BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&summary.TotalTransaksi)

	if err != nil {
		return nil, err
	}

	// Produk Terlaris
	row := r.DB.QueryRow(`
		SELECT p.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at::date BETWEEN $1 AND $2
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`, startDate, endDate)

	bestSeller := &models.BestSeller{}
	err = row.Scan(&bestSeller.Nama, &bestSeller.QtyTerjual)

	if err == nil {
		summary.ProdukTerlaris = bestSeller
	}

	return summary, nil
}
