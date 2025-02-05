package repository

import (
	"database/sql"
	"fmt"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/xuri/excelize/v2"
)

type ExcelDownloadMI struct {
	db *sql.DB
}

func NewExcelDownloadMIRepo(db *sql.DB) *ExcelDownloadMI {
	return &ExcelDownloadMI{db: db}
}

func (edr *ExcelDownloadMI) FetchExelMi() ([]models.ExcelDownloadMI, error) {
	var data []models.ExcelDownloadMI

	if edr.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := edr.db.Query("SELECT * FROM submitteddata")
	if err != nil {
		fmt.Println("Database query error:", err)
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record models.ExcelDownloadMI
		if err := rows.Scan(&record.Timestamp, &record.Supplier, &record.Buyer, &record.PartCode, &record.SerialNumber,
			&record.Quantity, &record.PONo, &record.PODate, &record.InvoiceNo, &record.InvoiceDate,
			&record.ReceivedDate, &record.UnitPricePerQty, &record.Category, &record.Warranty, &record.WarrantyDueDays); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		data = append(data, record)
	}

	if len(data) == 0 {
		fmt.Println("No records found in submitteddata table")
		return []models.ExcelDownloadMI{}, nil
	}

	fmt.Printf("Fetched %d records\n", len(data))
	return data, nil
}

func (edr *ExcelDownloadMI) CreateMaterialInward() (*excelize.File, error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	data, err := edr.FetchExelMi()
	if err != nil {
		return nil, err
	}

	sheetName := "MaterialInward"
	file.NewSheet(sheetName)

	headers := []string{
		"Timestamp", "Supplier", "Buyer", "PartCode", "SerialNumber", "Quantity", "PONo", "PODate",
		"InvoiceNo", "InvoiceDate", "ReceivedDate", "UnitPricePerQty", "Category", "Warranty", "WarrantyDueDays",
	}
	for colIndex, header := range headers {
		cell, err := excelize.CoordinatesToCellName(colIndex+1, 1)
		if err != nil {
			return nil, err
		}
		file.SetCellValue(sheetName, cell, header)
	}

	for i, record := range data {
		rowNum := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), record.Timestamp)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), record.Supplier)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), record.Buyer)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), record.PartCode)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), record.SerialNumber)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), record.Quantity)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), record.PONo)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), record.PODate)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), record.InvoiceNo)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), record.InvoiceDate)
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), record.ReceivedDate)
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), record.UnitPricePerQty)
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", rowNum), record.Category)
		file.SetCellValue(sheetName, fmt.Sprintf("N%d", rowNum), record.Warranty)
		file.SetCellValue(sheetName, fmt.Sprintf("O%d", rowNum), record.WarrantyDueDays)
	}

	return file, nil
}
