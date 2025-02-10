package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Srujankm12/SRstocks/repository"
	"github.com/xuri/excelize/v2"
)

type ExcelDownloadMOController struct {
	repo *repository.ExcelDownloadMORepo
}

func NewExcelDownloadMOController(repo *repository.ExcelDownloadMORepo) *ExcelDownloadMOController {
	return &ExcelDownloadMOController{repo: repo}
}
func (edc *ExcelDownloadMOController) DownloadMaterialOutward(w http.ResponseWriter, r *http.Request) {
	data, err := edc.repo.FetchExcelMO()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	file := excelize.NewFile()
	sheetName := "MaterialOutward"
	file.NewSheet(sheetName)

	headers := []string{
		"Timestamp", "Customer", "Seller", "BranchRegion", "PartCode", "SerialNumber", "Quantity", "CusPONo",
		"CusPODate", "CusInvoiceNo", "CusInvoiceDate", "DeliveredDate", "UnitPricePerQty", "IssuesAgainst",
		"Notes", "Category", "Warranty", "WarrantyDueDays",
	}

	for colIndex, header := range headers {
		cell := fmt.Sprintf("%s1", string(65+colIndex))
		file.SetCellValue(sheetName, cell, header)
	}

	// Fill data
	for i, record := range data {
		rowNum := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), record.Timestamp)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), record.Customer)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), record.Seller)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), record.BranchRegion)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), record.PartCode)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), record.SerialNumber)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), record.Quantity)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), record.CusPONo)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), record.CusPODate)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), record.CusInvoiceNo)
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), record.CusInvoiceDate)
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), record.DeliveredDate)
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", rowNum), record.UnitPricePerQty)
		file.SetCellValue(sheetName, fmt.Sprintf("N%d", rowNum), record.IssuesAgainst)
		file.SetCellValue(sheetName, fmt.Sprintf("O%d", rowNum), record.Notes)
		file.SetCellValue(sheetName, fmt.Sprintf("P%d", rowNum), record.Category)
		file.SetCellValue(sheetName, fmt.Sprintf("Q%d", rowNum), record.Warranty)
		file.SetCellValue(sheetName, fmt.Sprintf("R%d", rowNum), record.WarrantyDueDays)
	}

	tempDir := "/tmp"
	if os.Getenv("OS") == "Windows_NT" {
		tempDir = os.Getenv("TEMP")
	}

	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Error creating temp directory: %v", err), http.StatusInternalServerError)
		return
	}

	filePath := fmt.Sprintf("%s/MaterialOutward.xlsx", tempDir)

	if err := file.SaveAs(filePath); err != nil {
		http.Error(w, fmt.Sprintf("Error saving file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=MaterialOutward.xlsx")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, filePath)
}
