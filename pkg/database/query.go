package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/robfig/cron/v3"
)

type Query struct {
	db   *sql.DB
	cron *cron.Cron
}

func NewQuery(db *sql.DB) *Query {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatalf("Failed to load time zone: %v", err)
	}

	return &Query{
		db:   db,
		cron: cron.New(cron.WithLocation(loc)),
	}
}

func (q *Query) CreateTables() error {
	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS unit (
			unit_name VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS category (
			category_name VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS supplier (
			supplier_name VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS buyer (
			buyer_name VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS submitteddata (
			timestamp VARCHAR(255) NOT NULL,
			supplier VARCHAR(255) NOT NULL,
			buyer VARCHAR(255) NOT NULL,
			partcode VARCHAR(255) NOT NULL,
			serial_number VARCHAR(255) NOT NULL,
			qty INT NOT NULL,
			remaining_qty INT NOT NULL,
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
		`CREATE TABLE IF NOT EXISTS outwarddata (
			timestamp VARCHAR(255) NOT NULL,
			customer VARCHAR(255) NOT NULL,
			seller VARCHAR(255) NOT NULL,
			branch_region VARCHAR(255) NOT NULL,
			partcode VARCHAR(255) NOT NULL,
			serial_number VARCHAR(255) NOT NULL,
			quantity INT NOT NULL,
			cus_po_no VARCHAR(255) NOT NULL,
			cus_po_date DATE NOT NULL,
			cus_invoice_no VARCHAR(255) NOT NULL,
			cus_invoice_date DATE NOT NULL,
			delivered_date DATE NOT NULL,
			unit_price_per_qty FLOAT NOT NULL,
			issue_against VARCHAR(255) NOT NULL,
			notes TEXT NOT NULL,
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

func (q *Query) FetchFormOutwardDropDown() ([]models.OutwardDropDown, error) {
	var formdatasoutward []models.OutwardDropDown
	rows, err := q.db.Query("SELECT s.seller_value, b.region_value, i.issueagainst_value FROM seller s CROSS JOIN region b CROSS JOIN issueagainst i;")
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var formOutwardData models.OutwardDropDown
		if err := rows.Scan(&formOutwardData.Seller, &formOutwardData.BranchRegion, &formOutwardData.IssuesAgainst); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		formdatasoutward = append(formdatasoutward, formOutwardData)
	}
	if err != nil {
		return nil, err
	}
	return formdatasoutward, nil
}
func (q *Query) SubmitFormOutwardData(material models.MaterialOutward) error {
	result, err := q.db.Exec(
		`INSERT INTO outwarddata (
			timestamp,
			customer,
			seller,
			branch_region,
			partcode,
			serial_number,
			quantity,
			cus_po_no,
			cus_po_date,
			cus_invoice_no,
			cus_invoice_date,
			delivered_date,
			unit_price_per_qty,
			issue_against,
			notes,
			category,
			warranty,
			warranty_due_days
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)`,
		material.Timestamp,
		material.Customer,
		material.Seller,
		material.BranchRegion,
		material.PartCode,
		material.SerialNumber,
		material.Quantity,
		material.CusPONo,
		material.CusPODate,
		material.CusInvoiceNo,
		material.CusInvoiceDate,
		material.DeliveredDate,
		material.UnitPricePerQty,
		material.IssuesAgainst,
		material.Notes,
		material.Category,
		material.Warranty,
		material.Warranty,
	)
	if err != nil {
		log.Printf("Error inserting outward data: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching affected rows: %v", err)
		return err
	}
	log.Printf("Outward data inserted successfully. %d rows affected.", rowsAffected)
	return nil
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
			remaining_qty,
			po_no,
			po_date,
			invoice_no,
			invoice_date,
			received_date,
			unit_price_per_qty,
			category,
			warranty,
			warranty_due_days
		) VALUES (
			$1, $2, $3, $4, $5, $6, $6, $7, $8, $9, $10, $11, $12, $13, $14 ,$15
			
		)`,
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
		material.Warranty,
	)
	return err
}
func (q *Query) FetchAllFormOutwardData() ([]models.MaterialOutward, error) {
	var materials []models.MaterialOutward
	rows, err := q.db.Query(`
    SELECT 
        timestamp, customer, seller, branch_region, partcode, serial_number, quantity,
        cus_po_no, cus_po_date, cus_invoice_no, cus_invoice_date, 
        delivered_date, unit_price_per_qty, issue_against, notes, category, warranty, warranty_due_days
    FROM outwarddata
`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var material models.MaterialOutward
		err := rows.Scan(
			&material.Timestamp, &material.Customer, &material.Seller, &material.BranchRegion,
			&material.PartCode, &material.SerialNumber, &material.Quantity, &material.CusPONo,
			&material.CusPODate, &material.CusInvoiceNo, &material.CusInvoiceDate, &material.DeliveredDate,
			&material.UnitPricePerQty, &material.IssuesAgainst, &material.Notes, &material.Category,
			&material.Warranty, &material.WarrantyDueDays,
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
	// Schedule the cron job to run every 1 second
	_, err := q.cron.AddFunc("@every 1s", func() {
		log.Println("Cron job triggered: Updating warranty_due_days")

		// Update warranty_due_days for submitteddata
		result, err := q.db.Exec(`
			UPDATE submitteddata
			SET warranty_due_days = GREATEST(warranty_due_days - 1, 0)
			WHERE warranty_due_days > 0
		`)
		if err != nil {
			log.Printf("Error updating warranty_due_days in submitteddata: %v", err)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Error fetching affected rows for submitteddata: %v", err)
			return
		}
		log.Printf("Warranty due days updated in submitteddata. %d rows affected.", rowsAffected)

		// Update warranty_due_days for outwarddata
		result, err = q.db.Exec(`
			UPDATE outwarddata
			SET warranty_due_days = GREATEST(warranty_due_days - 1, 0)
			WHERE warranty_due_days > 0
		`)
		if err != nil {
			log.Printf("Error updating warranty_due_days in outwarddata: %v", err)
			return
		}

		rowsAffected, err = result.RowsAffected()
		if err != nil {
			log.Printf("Error fetching affected rows for outwarddata: %v", err)
			return
		}
		log.Printf("Warranty due days updated in outwarddata. %d rows affected.", rowsAffected)
	})
	if err != nil {
		log.Fatalf("Error scheduling cron job: %v", err)
	}

	// Start the cron scheduler
	q.cron.Start()
	log.Println("Cron job for updating warranty due days started successfully.")
}
