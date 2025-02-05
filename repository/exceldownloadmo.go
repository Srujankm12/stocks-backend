package repository

import (
	"database/sql"
	"fmt"

	"github.com/Srujankm12/SRstocks/internal/models"
)

type ExcelDownloadMORepo struct {
	db *sql.DB
}

func NewExcelDownloadMORepo(db *sql.DB) *ExcelDownloadMORepo {
	return &ExcelDownloadMORepo{db: db}
}

func (edr *ExcelDownloadMORepo) FetchExelMO() ([]models.ExcelDownloadMO, error) {

	var data []models.ExcelDownloadMO
	if edr.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	rows, err := edr.db.Query("SELECT * FROM material_outward")
	if err != nil {
		fmt.Println("Database query error:", err)
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var record models.ExcelDownloadMO
		if err := rows.Scan(&record.Timestamp, &record.Customer, &record.Seller, &record.BranchRegion, &record.PartCode, &record.SerialNumber, &record.Quantity, &record.CusPONo, &record.CusPODate, &record.CusInvoiceNo, &record.CusInvoiceDate, &record.DeliveredDate, &record.UnitPricePerQty, &record.IssuesAgainst, &record.Notes, &record.Category, &record.Warranty, &record.WarrantyDueDays); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		data = append(data, record)

	}
	if len(data) == 0 {
		fmt.Println("No records found")
		return []models.ExcelDownloadMO{}, nil

	}
	fmt.Printf("Fetched %d records\n", len(data))
	return data, nil

}
