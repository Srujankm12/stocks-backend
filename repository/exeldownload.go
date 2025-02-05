package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/xuri/excelize/v2"
)

type ExcelDownloadMSRepo struct {
	db *sql.DB
}

func NewExcelDownloadMSRepo(db *sql.DB) *ExcelDownloadMSRepo {
	return &ExcelDownloadMSRepo{db: db}
}

func (edr *ExcelDownloadMSRepo) FetchAllData() ([]models.ExcelDownloadMS, error) {
	var data []models.ExcelDownloadMS
	rows, err := edr.db.Query("SELECT * FROM material_stock")
	if err != nil {
		return nil, fmt.Errorf("Error executing query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record models.ExcelDownloadMS
		if err := rows.Scan(&record.ID, &record.Timestamp, &record.Supplier, &record.Category, &record.LeadTime,
			&record.StdNonStd, &record.PartCode, &record.Unit, &record.Rate, &record.MinimumRetain,
			&record.MaximumRetain, &record.Received, &record.Issue, &record.ReservedStock, &record.Stock,
			&record.Value, &record.ReorderStatus, &record.ExcessStock, &record.ExcessStockValue); err != nil {
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}
		data = append(data, record)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("No records found")
	}

	return data, nil
}

func (edr *ExcelDownloadMSRepo) CreateMaterialStock() (*excelize.File, error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	data, err := edr.FetchAllData()
	if err != nil {
		return nil, err
	}

	sheetName := "MaterialStock"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	file.SetActiveSheet(index)

	headers := []string{"ID", "Timestamp", "Supplier", "Category", "Lead Time", "Std/Non-Std", "Part Code", "Unit", "Rate", "Min Retain", "Max Retain", "Received", "Issue", "Reserved Stock", "Stock", "Value", "Reorder Status", "Excess Stock", "Excess Stock Value"}
	for colIndex, header := range headers {
		col := columnIndexToLetter(colIndex+1) + "1"
		file.SetCellValue(sheetName, col, header)
	}

	for rowIndex, record := range data {
		rowNum := strconv.Itoa(rowIndex + 2)
		file.SetCellValue(sheetName, "A"+rowNum, record.ID)
		file.SetCellValue(sheetName, "B"+rowNum, record.Timestamp)
		file.SetCellValue(sheetName, "C"+rowNum, record.Supplier)
		file.SetCellValue(sheetName, "D"+rowNum, record.Category)
		file.SetCellValue(sheetName, "E"+rowNum, record.LeadTime)
		file.SetCellValue(sheetName, "F"+rowNum, record.StdNonStd)
		file.SetCellValue(sheetName, "G"+rowNum, record.PartCode)
		file.SetCellValue(sheetName, "H"+rowNum, record.Unit)
		file.SetCellValue(sheetName, "I"+rowNum, record.Rate)
		file.SetCellValue(sheetName, "J"+rowNum, record.MinimumRetain)
		file.SetCellValue(sheetName, "K"+rowNum, record.MaximumRetain)
		file.SetCellValue(sheetName, "L"+rowNum, record.Received)
		file.SetCellValue(sheetName, "M"+rowNum, record.Issue)
		file.SetCellValue(sheetName, "N"+rowNum, record.ReservedStock)
		file.SetCellValue(sheetName, "O"+rowNum, record.Stock)
		file.SetCellValue(sheetName, "P"+rowNum, record.Value)
		file.SetCellValue(sheetName, "Q"+rowNum, record.ReorderStatus)
		file.SetCellValue(sheetName, "R"+rowNum, record.ExcessStock)
		file.SetCellValue(sheetName, "S"+rowNum, record.ExcessStockValue)
	}

	return file, nil
}

func columnIndexToLetter(index int) string {
	result := ""
	for index > 0 {
		index--
		result = string(rune('A'+(index%26))) + result
		index /= 26
	}
	return result
}
