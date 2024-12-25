package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/robfig/cron/v3"
)

type Query struct {
	db *sql.DB
}

func NewQuery(db *sql.DB) *Query {
	return &Query{
		db: db,
	}
}

func (q *Query) CreateTables() error {
	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS unit (
			unit_value VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS category (
			category_value VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS supplier (
			supplier_value VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS leadtime (
			leadtime_value VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS std (
			std_value VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS warranty (
			warranty_value VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS issueagainst (
			issueagainst_value VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS seller (
			seller_value VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS buyer (
			buyer_value VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS region (
			region_value VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS submitteddata (
			timestamp VARCHAR(255) NOT NULL,
			supplier VARCHAR(255) NOT NULL,
			buyer VARCHAR(255) NOT NULL,
			partcode VARCHAR(255) NOT NULL,
			serial_number VARCHAR(255) NOT NULL,
			qty INT NOT NULL,
			po_no VARCHAR(255) NOT NULL,
			po_date DATE NOT NULL,
			invoice_no VARCHAR(255) NOT NULL,
			invoice_date DATE NOT NULL,
			received_date DATE NOT NULL,
			unit_price_per_qty DECIMAL(10, 2) NOT NULL,
			category VARCHAR(255) NOT NULL,
			warranty INT NOT NULL,
			warranty_due_days INT NOT NULL
		);`,
	}

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			log.Printf("Failed to execute query: %s\nError: %v", query, err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	log.Println("All tables created successfully.")
	return nil
}

func (q *Query) FetchFormData() ([]models.InwardDropDown, error) {
	var formdatas []models.InwardDropDown
	rows, err := q.db.Query("SELECT b.buyer_value, s.supplier_value FROM buyer b CROSS JOIN supplier s;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var formdata models.InwardDropDown
		if err := rows.Scan(&formdata.Buyer, &formdata.Supplier); err != nil {
			return nil, err
		}
		formdatas = append(formdatas, formdata)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return formdatas, nil
}

func (q *Query) SubmitFormData(material models.MaterialInward) error {
	_, err := q.db.Exec(
		`INSERT INTO submitteddata (
			timestamp,
			supplier,
			buyer,
			partcode,
			serial_number,
			qty,
			po_no,
			po_date,
			invoice_no,
			invoice_date,
			received_date,
			unit_price_per_qty,
			category,
			warranty,
			warranty_due_days
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		material.Timestamp,
		material.Supplier,
		material.Buyer,
		material.PartCode,
		material.SerialNumber,
		material.Quantity,
		material.PONo,
		material.PODate,
		material.InvoiceNo,
		material.InvoiceDate,
		material.ReceivedDate,
		material.UnitPricePerQty,
		material.Category,
		material.Warranty,
		material.WarrantyDueDays,
	)
	return err
}

func (q *Query) FetchAllFormData() ([]models.MaterialInward, error) {
	var materials []models.MaterialInward
	rows, err := q.db.Query(`
        SELECT 
            timestamp, supplier, buyer, partcode, serial_number, qty, po_no, po_date, 
            invoice_no, invoice_date, received_date, unit_price_per_qty, category, 
            warranty, warranty_due_days 
        FROM submitteddata
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var material models.MaterialInward
		err := rows.Scan(
			&material.Timestamp, &material.Supplier, &material.Buyer, &material.PartCode,
			&material.SerialNumber, &material.Quantity, &material.PONo, &material.PODate,
			&material.InvoiceNo, &material.InvoiceDate, &material.ReceivedDate,
			&material.UnitPricePerQty, &material.Category, &material.Warranty,
			&material.WarrantyDueDays,
		)
		if err != nil {
			return nil, err
		}
		materials = append(materials, material)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return materials, nil
}

func (q *Query) UpdateWarrantyDueDays() {
	// Load the time zone for Bengaluru (Asia/Kolkata)
	loc, err := time.LoadLocation("Asia/Kolkata") // Bengaluru time zone
	if err != nil {
		log.Fatalf("Failed to load time zone: %v", err)
	}

	// Create a new cron job with the specified time zone
	c := cron.New(cron.WithLocation(loc))

	// Schedule the cron job to run at midnight every day
	_, err = c.AddFunc("0 0 * * *", func() {
		log.Println("Cron job triggered: Updating warranty_due_days")

		// Execute the update query to decrement warranty_due_days
		result, err := q.db.Exec(`
			UPDATE submitteddata 
			SET warranty_due_days = warranty_due_days - 1 
			WHERE warranty_due_days > 0
		`)
		if err != nil {
			log.Printf("Error updating warranty_due_days: %v", err)
		} else {
			// Get the number of affected rows
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				log.Printf("Error fetching affected rows: %v", err)
			} else {
				log.Printf("Warranty due days updated successfully. %d rows affected.", rowsAffected)
			}
		}
	})

	if err != nil {
		log.Fatalf("Error scheduling cron job: %v", err)
	}

	// Start the cron job
	c.Start()
	log.Println("Cron job for updating warranty due days started successfully.")
}
