package reporting

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GenerateAuditPDF(ctx context.Context, start, end time.Time) ([]byte, error) {
	count, err := s.repo.GetDepositedItemsCount(ctx, start, end)
	if err != nil {
		return nil, err
	}

	dist, _ := s.repo.GetPredictionDistribution(ctx)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 51, 102)
	pdf.Cell(0, 10, "UPCYCLE CONNECT - AUDIT REPORT")
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 10, "Description de l'audit : Extraction des objets deposes")
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Periode : du %s au %s", start.Format("02/01/2006"), end.Format("02/01/2006")))
	pdf.Ln(7)
	pdf.Cell(0, 10, fmt.Sprintf("Genere le : %s", time.Now().Format("02/01/2006 15:04")))
	pdf.Ln(15)

	pdf.SetFillColor(230, 230, 230)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(120, 10, "Indicateur", "1", 0, "L", true, 0, "")
	pdf.CellFormat(70, 10, "Valeur", "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(120, 10, "Nombre total d'objets deposes", "1", 0, "L", false, 0, "")
	pdf.CellFormat(70, 10, fmt.Sprintf("%d", count), "1", 0, "C", false, 0, "")
	pdf.Ln(20)

	if len(dist) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 10, "Predictions ML (Interets utilisateurs)")
		pdf.Ln(10)

		pdf.SetFont("Arial", "B", 11)
		pdf.CellFormat(120, 10, "Service Predi", "1", 0, "L", true, 0, "")
		pdf.CellFormat(70, 10, "Nombre d'utilisateurs", "1", 0, "C", true, 0, "")
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 11)
		for sType, sCount := range dist {
			pdf.CellFormat(120, 10, sType, "1", 0, "L", false, 0, "")
			pdf.CellFormat(70, 10, fmt.Sprintf("%d", sCount), "1", 0, "C", false, 0, "")
			pdf.Ln(10)
		}
		pdf.Ln(10)
	}

	pdf.SetFont("Arial", "I", 9)
	pdf.SetTextColor(100, 100, 100)
	pdf.MultiCell(0, 5, "Ce document constitue un audit officiel de l'activite de la plateforme UpcycleConnect pour la periode specifiee. Les donnees sont extraites directement de la base de donnees de production.", "", "L", false)

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Service) GetActorStats(ctx context.Context) ([]ActorStats, error) {
	return s.repo.GetActorStats(ctx)
}

func (s *Service) GetPrestationStats(ctx context.Context) ([]PrestationStats, error) {
	return s.repo.GetPrestationStats(ctx)
}

func (s *Service) GetUserPredictions(ctx context.Context, page, limit int) (*PaginatedPredictions, error) {
	return s.repo.GetUserPredictions(ctx, page, limit)
}

func (s *Service) GetMLStatus(ctx context.Context) (*MLStatus, error) {
	return s.repo.GetMLStatus(ctx)
}
