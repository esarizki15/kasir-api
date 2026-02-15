package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// untuk /api/report/hari-ini
func (s *ReportService) GetTodaySummary() (*models.SalesSummary, error) {
	return s.repo.GetTodaySummary()
}

// untuk /api/report?start_date=...&end_date=...
func (s *ReportService) GetSummaryByDate(start, end string) (*models.SalesSummary, error) {
	return s.repo.GetSummaryByDate(start, end)
}
