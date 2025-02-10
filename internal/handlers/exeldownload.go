package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Srujankm12/SRstocks/repository"
	"github.com/xuri/excelize/v2"
)

type ExcelDownloadController struct {
	repo *repository.ExcelDownloadMSRepo
}

func NewExcelDownloadController(repo *repository.ExcelDownloadMSRepo) *ExcelDownloadController {
	return &ExcelDownloadController{repo: repo}
}

func (edc *ExcelDownloadController) DownloadMaterialStock(w http.ResponseWriter, r *http.Request) {

	data, err := edc.repo.FetchAllData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	file := excelize.NewFile()
	sheetName := "MaterialStock"
	file.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"ID", "Timestamp", "Supplier", "Category", "LeadTime", "StdNonStd", "PartCode",
		"Unit", "Rate", "MinimumRetain", "MaximumRetain", "Received", "Issue",
		"ReservedStock", "Stock", "Value", "ReorderStatus", "ExcessStock", "ExcessStockValue",
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
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), record.Category)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), record.LeadTime)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), record.StdNonStd)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), record.PartCode)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), record.Unit)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), record.Rate)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), record.MinimumRetain)
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), record.MaximumRetain)
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), record.Received)
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", rowNum), record.Issue)
		file.SetCellValue(sheetName, fmt.Sprintf("N%d", rowNum), record.ReservedStock)
		file.SetCellValue(sheetName, fmt.Sprintf("O%d", rowNum), record.Stock)
		file.SetCellValue(sheetName, fmt.Sprintf("P%d", rowNum), record.Value)
		file.SetCellValue(sheetName, fmt.Sprintf("Q%d", rowNum), record.ReorderStatus)
		file.SetCellValue(sheetName, fmt.Sprintf("R%d", rowNum), record.ExcessStock)
		file.SetCellValue(sheetName, fmt.Sprintf("S%d", rowNum), record.ExcessStockValue)
	}

	tempDir := "/tmp"
	if os.Getenv("OS") == "Windows_NT" {
		tempDir = os.Getenv("TEMP")
	}

	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Error creating temp directory: %v", err), http.StatusInternalServerError)
		return
	}

	filePath := fmt.Sprintf("%s/MaterialStock.xlsx", tempDir)

	if err := file.SaveAs(filePath); err != nil {
		http.Error(w, fmt.Sprintf("Error saving file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=MaterialStock.xlsx")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, filePath)
}
