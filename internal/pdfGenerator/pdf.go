package pdfGenerator

import (
	"LinkChecker/internal/models"
	"bytes"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

func (g *PDFGenerator) Generate(groups []models.LinksGroup) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)

	pdf.Cell(0, 10, "Links report")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(20, 8, "ID", "1", 0, "", false, 0, "")
	pdf.CellFormat(90, 8, "Link", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 8, "Status", "1", 0, "", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)

	for _, group := range groups {
		for _, link := range group.Links {
			pdf.CellFormat(20, 8, strconv.Itoa(group.ID), "1", 0, "", false, 0, "")
			pdf.CellFormat(90, 8, link.URL, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 8, link.Status, "1", 0, "", false, 0, "")
			pdf.Ln(-1)
		}
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
