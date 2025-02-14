package database

import (
	"database/sql"
	"fmt"
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
		`CREATE TABLE IF NOT EXISTS branch_region (
			region_name VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS issueagainst (
			issueagainst_name VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS seller (
			seller_name VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS submitteddata (
		    id SERIAL PRIMARY KEY,
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
		)`,
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
		)`,
		`CREATE TABLE IF NOT EXISTS material_stock (
			id SERIAL PRIMARY KEY,
			timestamp VARCHAR(255) NOT NULL,
			supplier VARCHAR(255) NOT NULL,
			category VARCHAR(255),
			lead_time INT NOT NULL,
			std_non_std VARCHAR(50),
			part_code VARCHAR(255) NOT NULL,
			unit VARCHAR(50),
			rate DECIMAL(10, 2),
			minimum_retain INT NOT NULL,
			maximum_retain INT NOT NULL,
			received INT DEFAULT 0,
			issue INT DEFAULT 0,
			reserved_stock INT DEFAULT 0,
			stock INT NOT NULL,
			value DECIMAL(10, 2) NOT NULL,
			reorder_status BOOLEAN NOT NULL,
			excess_stock INT NOT NULL,
			excess_stock_value DECIMAL(10, 2) NOT NULL
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
	rows, err := q.db.Query("SELECT b.buyer_name, s.supplier_name ,c.category_name FROM buyer b CROSS JOIN supplier s CROSS JOIN category c;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var formdata models.InwardDropDown
		if err := rows.Scan(&formdata.Buyer, &formdata.Supplier, &formdata.Category); err != nil {
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
	rows, err := q.db.Query("SELECT s.seller_name, b.region_name, i.issueagainst_name , c.category_name FROM seller s CROSS JOIN branch_region b CROSS JOIN issueagainst i CROSS JOIN category c;")
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var formOutwardData models.OutwardDropDown
		if err := rows.Scan(&formOutwardData.Seller, &formOutwardData.BranchRegion, &formOutwardData.IssuesAgainst, &formOutwardData.Category); err != nil {
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

func (q *Query) FetchMaterialDropdownData() ([]models.MaterialStockDropDown, error) {
	var materialdropdown []models.MaterialStockDropDown
	rows, err := q.db.Query("SELECT supplier_name,category_name,unit_name FROM supplier s CROSS JOIN category c CROSS JOIN unit u;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var materialdropdowndata models.MaterialStockDropDown
		if err := rows.Scan(&materialdropdowndata.Supplier, &materialdropdowndata.Category, &materialdropdowndata.Unit); err != nil {
			return nil, err
		}
		materialdropdown = append(materialdropdown, materialdropdowndata)
	}
	if err != nil {
		return nil, err
	}
	return materialdropdown, nil
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
			timestamp, supplier, buyer, partcode, serial_number, qty, po_no, 
			po_date, invoice_no, invoice_date, received_date, unit_price_per_qty, 
			category, warranty, warranty_due_days
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
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
func (q *Query) SubmitMaterialStock(material models.MaterialStock) (int, error) {
	if material.Timestamp == "" {
		material.Timestamp = time.Now().Format(time.RFC3339)
	}

	material.Stock = material.Received - material.Issue

	if material.Stock < material.MinimumRetain {
		material.ReorderStatus = true
	} else {
		material.ReorderStatus = false
	}

	if material.Stock > material.MaximumRetain {
		material.ExcessStock = material.Stock - material.MaximumRetain
		material.ExcessStockValue = float64(material.ExcessStock) * material.Rate
	} else {
		material.ExcessStock = 0
		material.ExcessStockValue = 0
	}

	var id int
	err := q.db.QueryRow(
		`INSERT INTO material_stock (
			timestamp, supplier, category, lead_time, std_non_std, part_code, unit, rate,
			minimum_retain, maximum_retain, received, issue, reserved_stock, stock, value,
			reorder_status, excess_stock, excess_stock_value
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) RETURNING id;`,
		material.Timestamp,
		material.Supplier,
		material.Category,
		material.LeadTime,
		material.StdNonStd,
		material.PartCode,
		material.Unit,
		material.Rate,
		material.MinimumRetain,
		material.MaximumRetain,
		material.Received,
		material.Issue,
		material.ReservedStock,
		material.Stock,
		float64(material.Stock)*material.Rate,
		material.ReorderStatus,
		material.ExcessStock,
		material.ExcessStockValue,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
func (q *Query) FetchAllMaterialStock() ([]models.MaterialStock, error) {
	var materials []models.MaterialStock
	rows, err := q.db.Query(`
    SELECT id, timestamp, supplier, category, lead_time, std_non_std, part_code, unit,
           rate, minimum_retain, maximum_retain, received, issue, reserved_stock, stock,
           value, reorder_status, excess_stock, excess_stock_value
    FROM material_stock;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var material models.MaterialStock
		err := rows.Scan(
			&material.ID,
			&material.Timestamp,
			&material.Supplier,
			&material.Category,
			&material.LeadTime,
			&material.StdNonStd,
			&material.PartCode,
			&material.Unit,
			&material.Rate,
			&material.MinimumRetain,
			&material.MaximumRetain,
			&material.Received,
			&material.Issue,
			&material.ReservedStock,
			&material.Stock,
			&material.Value,
			&material.ReorderStatus,
			&material.ExcessStock,
			&material.ExcessStockValue,
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
            id, timestamp, supplier, buyer, partcode, serial_number, qty, po_no, po_date, 
            invoice_no, invoice_date, received_date, unit_price_per_qty, category, 
            warranty, warranty_due_days 
        FROM submitteddata
        ORDER BY received_date ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var material models.MaterialInward
		err := rows.Scan(
			&material.ID, &material.Timestamp, &material.Supplier, &material.Buyer,
			&material.PartCode, &material.SerialNumber, &material.Quantity, &material.PONo,
			&material.PODate, &material.InvoiceNo, &material.InvoiceDate, &material.ReceivedDate,
			&material.UnitPricePerQty, &material.Category, &material.Warranty, &material.WarrantyDueDays,
		)
		if err != nil {
			return nil, err
		}

		// 🔍 Debugging: Print fetched row
		fmt.Printf("Fetched: %+v\n", material)

		materials = append(materials, material)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return materials, nil
}

func (q *Query) UpdateMaterialStock(material models.MaterialStock) error {
	material.Stock = material.Received - material.Issue
	material.Value = float64(material.Stock) * material.Rate

	material.ReorderStatus = material.Stock < material.MinimumRetain

	if material.Stock > material.MaximumRetain {
		material.ExcessStock = material.Stock - material.MaximumRetain
		material.ExcessStockValue = float64(material.ExcessStock) * material.Rate
	} else {
		material.ExcessStock = 0
		material.ExcessStockValue = 0
	}

	fmt.Printf("Updating Lead Time: %d for ID: %d\n", material.LeadTime, material.ID)

	query := `UPDATE material_stock SET
            timestamp = $1,
            supplier = $2,
            category = $3,
          	 lead_time = $4,
            std_non_std = $5,
            part_code = $6,
            unit = $7,
            rate = $8,
            minimum_retain = $9,
            maximum_retain = $10,
            received = $11,
            issue = $12,
            reserved_stock = $13,
            stock = $14,
            value = $15,	
            reorder_status = $16,
            excess_stock = $17,
            excess_stock_value = $18
        WHERE id = $19;`

	result, err := q.db.Exec(query,
		material.Timestamp,
		material.Supplier,
		material.Category,
		material.LeadTime,
		material.StdNonStd,
		material.PartCode,
		material.Unit,
		material.Rate,
		material.MinimumRetain,
		material.MaximumRetain,
		material.Received,
		material.Issue,
		material.ReservedStock,
		material.Stock,
		material.Value,
		material.ReorderStatus,
		material.ExcessStock,
		material.ExcessStockValue,
		material.ID,
	)
	if err != nil {
		fmt.Printf("Error updating material_stock: %v\n", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("Warning: No rows were updated. Check if ID exists.")
	}

	return nil
}
func (q *Query) FetchAllData() ([]models.ExcelDownloadMS, error) {
	var data []models.ExcelDownloadMS
	rows, err := q.db.Query("SELECT * FROM material_stock")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record models.ExcelDownloadMS
		if err := rows.Scan(&record.ID, &record.Timestamp, &record.Supplier, &record.Category, &record.LeadTime,
			&record.StdNonStd, &record.PartCode, &record.Unit, &record.Rate, &record.MinimumRetain,
			&record.MaximumRetain, &record.Received, &record.Issue, &record.ReservedStock, &record.Stock,
			&record.Value, &record.ReorderStatus, &record.ExcessStock, &record.ExcessStockValue); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		data = append(data, record)
	}

	if len(data) == 0 {
		fmt.Println("No records found")
		return nil, nil
	}

	return data, nil
}
func (q *Query) FetchExelMi() ([]models.ExcelDownloadMI, error) {
	var data []models.ExcelDownloadMI
	rows, err := q.db.Query("SELECT * FROM submitteddata")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record models.ExcelDownloadMI
		if err := rows.Scan(&record.ID, &record.Timestamp, &record.Supplier, &record.Buyer, &record.PartCode, &record.SerialNumber,
			&record.Quantity, &record.PONo, &record.PODate, &record.InvoiceNo, &record.InvoiceDate,
			&record.ReceivedDate, &record.UnitPricePerQty, &record.Category, &record.Warranty, &record.WarrantyDueDays); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		data = append(data, record)
	}
	if len(data) == 0 {
		fmt.Println("No records found")
		return nil, nil
	}
	return data, nil
}

func (q *Query) FetchExelMO() ([]models.ExcelDownloadMO, error) {
	var data []models.ExcelDownloadMO
	rows, err := q.db.Query("SELECT * FROM outwarddata")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
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
		return nil, nil
	}
	return data, nil
}
func (q *Query) UpdateWarrantyDueDays() {

	_, err := q.cron.AddFunc("@every 1s", func() {
		log.Println("Cron job triggered: Updating warranty_due_days")

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

	q.cron.Start()
	log.Println("Cron job for updating warranty due days started successfully.")
}
func (q *Query) GetMaterialStock() (map[string]int, error) {
	stock := make(map[string]int)

	rows, err := q.db.Query(`
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
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var partCode string
		var totalReceived, totalIssued, stockValue int

		err := rows.Scan(&partCode, &totalReceived, &totalIssued, &stockValue)
		if err != nil {
			return nil, err
		}

		stock[partCode] = stockValue
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stock, nil
}
