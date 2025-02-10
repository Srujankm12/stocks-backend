package repository

import (
	"database/sql"
	"fmt"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/xuri/excelize/v2"
)

type ExcelDownloadMORepo struct {
	db *sql.DB
}

func NewExcelDownloadMORepo(db *sql.DB) *ExcelDownloadMORepo {
	return &ExcelDownloadMORepo{db: db}
}

func (edr *ExcelDownloadMORepo) FetchExcelMO() ([]models.ExcelDownloadMO, error) {
	var data []models.ExcelDownloadMO
	if edr.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := edr.db.Query("SELECT * FROM outwarddata")
	if err != nil {
		fmt.Println("Database query error:", err)
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record models.ExcelDownloadMO
		if err := rows.Scan(
			&record.Timestamp, &record.Customer, &record.Seller, &record.BranchRegion,
			&record.PartCode, &record.SerialNumber, &record.Quantity, &record.CusPONo,
			&record.CusPODate, &record.CusInvoiceNo, &record.CusInvoiceDate,
			&record.DeliveredDate, &record.UnitPricePerQty, &record.IssuesAgainst,
			&record.Notes, &record.Category, &record.Warranty, &record.WarrantyDueDays,
		); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		data = append(data, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	if len(data) == 0 {
		fmt.Println("No records found")
		return []models.ExcelDownloadMO{}, nil
	}

	fmt.Printf("Fetched %d records\n", len(data))
	return data, nil
}

func (edr *ExcelDownloadMORepo) CreateMaterialOutward() (*excelize.File, error) {
	file := excelize.NewFile()
	data, err := edr.FetchExcelMO()
	if err != nil {
		return nil, err
	}

	sheetName := "Material Outward"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	file.SetActiveSheet(index)
	file.DeleteSheet("Sheet1")

	headers := []string{
		"Timestamp", "Customer", "Seller", "Branch Region", "Part Code",
		"Serial Number", "Quantity", "Customer PO No", "Customer PO Date",
		"Customer Invoice No", "Customer Invoice Date", "Delivered Date",
		"Unit Price Per Qty", "Issues Against", "Notes", "Category", "Warranty",
		"Warranty Due Days",
	}

	// Add headers
	for colIndex, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIndex+1, 1)
		file.SetCellValue(sheetName, cell, header)
	}

	// Add data
	for i, record := range data {
		row := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), record.Timestamp)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), record.Customer)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), record.Seller)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), record.BranchRegion)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), record.PartCode)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), record.SerialNumber)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", row), record.Quantity)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", row), record.CusPONo)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", row), record.CusPODate)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", row), record.CusInvoiceNo)
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", row), record.CusInvoiceDate)
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", row), record.DeliveredDate)
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", row), record.UnitPricePerQty)
		file.SetCellValue(sheetName, fmt.Sprintf("N%d", row), record.IssuesAgainst)
		file.SetCellValue(sheetName, fmt.Sprintf("O%d", row), record.Notes)
		file.SetCellValue(sheetName, fmt.Sprintf("P%d", row), record.Category)
		file.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), record.Warranty)
		file.SetCellValue(sheetName, fmt.Sprintf("R%d", row), record.WarrantyDueDays)
	}

	if len(data) == 0 {
		fmt.Println("No records found")
		return file, nil // Return an empty file instead of nil
	}

	return file, nil
}
