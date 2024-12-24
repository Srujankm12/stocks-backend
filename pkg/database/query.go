package database

import (
	"database/sql"
	"log"
)

type Query struct {
	db *sql.DB
}

func NewQuery(db *sql.DB) *Query {
	return &Query{
		db,
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
		`CREATE TABLE IF NOT EXISTS suppliers (
			suppliers_value VARCHAR(100) NOT NULL UNIQUE
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
func (q *Query) InsertSampleData() error {
	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return err
	}

	queriess := []string{
		`INSERT INTO unit (unit_value) VALUES 
			('Nos'), ('Mtrs'), ('Set'), ('Day'), ('Lot'), ('Kgs'), ('Ltrs'), ('Days')`,
		`INSERT INTO category (category_value) VALUES 
			('FX PLC'), ('HMI'), ('SERVO'), ('VFD'), ('S-Tech'), ('SR Systech')`,
		`INSERT INTO suppliers (suppliers_value) VALUES 
			('Mitsubishi Electric India'), ('Pheonix India'), ('Beijjer'), ('Altronix'), ('S-Tech'), ('SR Systech')`,
		`INSERT INTO leadtime(leadtime_value) VALUES
		('3-4months')`,
		`INSERT INTO std (std_value)VALUES
		('std'),('non std')`,
		`INSERT INTO issueagainst (issueagainst_value)VALUES
		('Invoice'),('Returnable DC '),('Non Returnable DC'),('Project'),('Other')`,
		`INSERT INTO seller (seller_value)VALUES
		('SRAB'),('SRCPL'),('EXL'),('Other')`,
		`INSERT INTO buyer (buyer_value)VALUES
		('SRAB'),('SRCPL'),('EXL'),('Other')`,
		`INSERT INTO region(region_value)VALUES
		('Hyderabad'),('Vizag'),('Chennai'),('Bangalore')`,
	}

	for _, query := range queriess {
		if _, err := tx.Exec(query); err != nil {
			log.Printf("Failed to execute query:\n%s\nError: %v", query, err)
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	log.Println("Sample data inserted successfully.")
	return nil
}
