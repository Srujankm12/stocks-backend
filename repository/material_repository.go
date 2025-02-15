package repository

import (
	"database/sql"
	"log"

	"github.com/Srujankm12/SRstocks/internal/models"
)

type MaterialRepository struct {
	db *sql.DB
}

func NewMaterialRepository(db *sql.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

func (r *MaterialRepository) GetMaterialStock() ([]models.MaterialStockUpdated, error) {
	var stocks []models.MaterialStockUpdated

	rows, err := r.db.Query(`
        WITH received_cte AS (
            SELECT partcode, SUM(qty) AS total_received
            FROM submitteddata
            GROUP BY partcode
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
            (COALESCE(r.total_received, 0) - COALESCE(i.total_issued, 0)) AS stock
        FROM received_cte r
        FULL OUTER JOIN issued_cte i ON r.partcode = i.partcode;
    `)

	if err != nil {
		log.Println("Error fetching material stock:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.MaterialStockUpdated
		err := rows.Scan(&stock.PartCode, &stock.Received, &stock.Issued, &stock.Stock)
		if err != nil {
			return nil, err
		}

		// Determine reorder status (if stock < 0, needs reorder)
		stock.ReorderStatus = stock.Stock < 1
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}
