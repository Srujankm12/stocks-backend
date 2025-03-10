package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Srujankm12/SRstocks/repository"
	"github.com/xuri/excelize/v2"
)

type ExcelDownloadMIController struct {
	repo *repository.ExcelDownloadMI
}

func NewExcelDownloadMIController(repo *repository.ExcelDownloadMI) *ExcelDownloadMIController {
	return &ExcelDownloadMIController{repo: repo}
}

func (edc *ExcelDownloadMIController) DownloadMaterialInward(w http.ResponseWriter, r *http.Request) {
	data, err := edc.repo.FetchExelMi()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	file := excelize.NewFile()
	sheetName := "MaterialInward"
	file.NewSheet(sheetName)

	headers := []string{
		"ID",
		"Timestamp", "Supplier", "Buyer", "PartCode", "SerialNumber", "Quantity", "PONo", "PODate",
		"InvoiceNo", "InvoiceDate", "ReceivedDate", "UnitPricePerQty", "Category", "Warranty", "WarrantyDueDays",
	}
	for colIndex, header := range headers {
		cell := fmt.Sprintf("%s1", string(65+colIndex))
		file.SetCellValue(sheetName, cell, header)
	}

	for i, record := range data {
		rowNum := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), record.ID)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), record.Timestamp)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), record.Supplier)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), record.Buyer)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), record.PartCode)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), record.SerialNumber)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), record.Quantity)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), record.PONo)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), record.PODate)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), record.InvoiceNo)
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), record.InvoiceDate)
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), record.ReceivedDate)
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", rowNum), record.UnitPricePerQty)
		file.SetCellValue(sheetName, fmt.Sprintf("N%d", rowNum), record.Category)
		file.SetCellValue(sheetName, fmt.Sprintf("O%d", rowNum), record.Warranty)
		file.SetCellValue(sheetName, fmt.Sprintf("P%d", rowNum), record.WarrantyDueDays)
	}
	tempDir := "/tmp"
	if os.Getenv("OS") == "Windows_NT" {
		tempDir = os.Getenv("TEMP")
	}

	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Error creating temp directory: %v", err), http.StatusInternalServerError)
		return
	}

	filePath := fmt.Sprintf("%s/MaterialInward.xlsx", tempDir)

	if err := file.SaveAs(filePath); err != nil {
		http.Error(w, fmt.Sprintf("Error saving file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=MaterialInward.xlsx")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, filePath)
}
