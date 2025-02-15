package repository

import (
	"database/sql"
	"fmt"

	"github.com/Srujankm12/SRstocks/internal/models"
)

type MaterialRepository struct {
	db *sql.DB
}

func NewMaterialRepository(db *sql.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

func (r *MaterialRepository) GetMaterialStock() ([]models.MaterialStockUpdated, error) {
	var stockData []models.MaterialStockUpdated

	rows, err := r.db.Query(`
        WITH received_cte AS (
            SELECT partcode, SUM(qty) AS total_received
            FROM submitteddata
            GROUP BY partcode
        ),
        latest_price_cte AS (
            SELECT DISTINCT ON (partcode) partcode, unit_price_per_qty 
            FROM submitteddata 
            ORDER BY partcode, id DESC  -- Fetch the latest unit price based on ID (assuming 'id' is auto-increment)
        ),
        issued_cte AS (
            SELECT partcode, SUM(quantity) AS total_issued
            FROM outwarddata
            GROUP BY partcode
        )
        SELECT 
            COALESCE(r.partcode, i.partcode) AS partcode,
            COALESCE(r.total_received, 0) AS total_received,
            COALESCE(i.total_issued, 0) AS total_issued,
            (COALESCE(r.total_received, 0) - COALESCE(i.total_issued, 0)) AS stock,
            COALESCE(lp.unit_price_per_qty, 0) * (COALESCE(r.total_received, 0) - COALESCE(i.total_issued, 0)) AS stock_value,
            CASE WHEN (COALESCE(r.total_received, 0) - COALESCE(i.total_issued, 0)) <= 0 THEN true ELSE false END AS reorder_status
        FROM received_cte r
        FULL OUTER JOIN issued_cte i ON r.partcode = i.partcode
        LEFT JOIN latest_price_cte lp ON lp.partcode = COALESCE(r.partcode, i.partcode);
    `)

	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.MaterialStockUpdated

		err := rows.Scan(&stock.PartCode, &stock.Received, &stock.Issued, &stock.Stock, &stock.StockValue, &stock.ReorderStatus)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		stockData = append(stockData, stock)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
		return nil, err
	}

	return stockData, nil
}
