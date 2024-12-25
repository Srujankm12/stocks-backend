package database

import (
	"database/sql"
	"log"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/robfig/cron/v3"
)

type Query struct {
	db *sql.DB
}

func NewQuery(db *sql.DB) *Query {
	return &Query{
		db,
	}
}

// CreateTables creates all required tables in the database.
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
		_, err := tx.Exec(query)
		if err != nil {
			tx.Rollback()
			log.Printf("Failed to execute query: %s\nError: %v", query, err)
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("All tables created successfully.")
	return nil
}

// FetchFormData retrieves the suppliers and buyers for dropdown options.
func (q *Query) FetchFormData() ([]models.InwardDropDown, error) {
	var formdata models.InwardDropDown
	var formdatas []models.InwardDropDown
	res, err := q.db.Query("SELECT b.buyer_value, s.supplier_value FROM buyer b CROSS JOIN supplier s;")
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		if err := res.Scan(&formdata.Buyer, &formdata.Supplier); err != nil {
			return nil, err
		}
		formdatas = append(formdatas, formdata)
	}
	if res.Err() != nil {
		return nil, err
	}
	return formdatas, nil
}

// SubmitFormData submits the form data into the database.
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
		material.WarrantyDueDays, // warranty_due_days is set to the initial value here
	)
	return err
}

// UpdateWarrantyDueDays decrements the warranty due days daily.
func (q *Query) UpdateWarrantyDueDays() {
	// Set up cron job to run every day at midnight
	c := cron.New()
	_, err := c.AddFunc("0 0 * * *", func() {
		// Decrement warranty_due_days for records where warranty_due_days > 0
		result, err := q.db.Exec(`UPDATE submitteddata SET warranty_due_days = warranty_due_days - 1 WHERE warranty_due_days > 0`)
		if err != nil {
			log.Printf("Error updating warranty_due_days: %v", err)
		} else {
			affectedRows, _ := result.RowsAffected() // Get the number of rows affected by the update
			log.Printf("Warranty due days updated successfully. %d rows affected.", affectedRows)
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling cron job: %v", err)
	}

	// Start the cron job
	c.Start()
}
