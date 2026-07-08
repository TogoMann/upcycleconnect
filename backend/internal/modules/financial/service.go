package financial

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jung-kurt/gofpdf"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCommissions() ([]Commission, error) {
	return s.repo.GetCommissions()
}

func (s *Service) GetFinancier() (*FinancierData, error) {
	return s.repo.GetFinancier()
}

func (s *Service) GetReport() (*FinancialReport, error) {
	return s.repo.GetReport()
}

func (s *Service) GetUserInvoices(userId int64) ([]Invoice, error) {
	return s.repo.GetUserInvoices(userId)
}

func (s *Service) GetInvoiceById(id int64) (*InvoiceDetail, error) {
	return s.repo.GetInvoiceById(id)
}

func (s *Service) CreateExpense(e Expense) error {
	return s.repo.CreateExpense(e)
}

func (s *Service) GetAllExpenses() ([]Expense, error) {
	return s.repo.GetAllExpenses()
}

func (s *Service) GeneratePDF(detail InvoiceDetail) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "UPCYCLE CONNECT")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, "Facture de prestation")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Facture n: %s", detail.InvoiceNumber))
	pdf.Ln(7)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Date: %s", detail.CreatedAt))
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 10, "Destinataire:")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, detail.UserName)
	pdf.Ln(5)
	pdf.Cell(0, 10, detail.UserEmail)
	pdf.Ln(20)

	pdf.SetFillColor(240, 240, 240)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(120, 10, "Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 10, "Qte", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Montant HT", "1", 0, "R", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(120, 10, detail.Description, "1", 0, "L", false, 0, "")
	pdf.CellFormat(30, 10, "1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f EUR", detail.Amount), "1", 0, "R", false, 0, "")
	pdf.Ln(20)

	pdf.SetX(140)
	pdf.CellFormat(30, 10, "Total HT", "0", 0, "L", false, 0, "")
	pdf.CellFormat(30, 10, fmt.Sprintf("%.2f EUR", detail.Amount), "0", 0, "R", false, 0, "")
	pdf.Ln(7)

	pdf.SetX(140)
	pdf.CellFormat(30, 10, "TVA (20%)", "0", 0, "L", false, 0, "")
	pdf.CellFormat(30, 10, fmt.Sprintf("%.2f EUR", detail.VatAmount), "0", 0, "R", false, 0, "")
	pdf.Ln(10)

	pdf.SetX(140)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(30, 10, "Total TTC", "0", 0, "L", false, 0, "")
	pdf.CellFormat(30, 10, fmt.Sprintf("%.2f EUR", detail.TotalAmount), "0", 0, "R", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Service) GenerateInvoiceForOrder(userId int64, sellerId *int64, orderId int64, orderType string, amount float64) error {
	vatRate := 0.20
	vatAmount := amount * vatRate
	totalAmount := amount + vatAmount

	invoiceNumber := fmt.Sprintf("INV-%d-%s-%d", userId, orderType[:1], time.Now().Unix())

	inv := Invoice{
		UserId:        userId,
		OrderId:       orderId,
		OrderType:     orderType,
		InvoiceNumber: invoiceNumber,
		Amount:        amount,
		VatAmount:     vatAmount,
		TotalAmount:   totalAmount,
		Status:        "paid",
	}

	if sellerId != nil {
		inv.SellerId = pgtype.Int8{Int64: *sellerId, Valid: true}
	}

	return s.repo.CreateInvoice(inv)
}
