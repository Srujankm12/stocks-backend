package models

import "net/http"

type MaterialOutward struct {
	ID              int     `json:"id"`
	Timestamp       string  `json:"timestamp"`
	Customer        string  `json:"customer"`
	Seller          string  `json:"seller"`
	BranchRegion    string  `json:"branch_region"`
	PartCode        string  `json:"partcode"`
	SerialNumber    string  `json:"serial_number"`
	Quantity        int     `json:"qty"`
	CusPONo         string  `json:"cus_po_no"`
	CusPODate       string  `json:"cus_po_date"`
	CusInvoiceNo    string  `json:"cus_invoice_no"`
	CusInvoiceDate  string  `json:"cus_invoice_date"`
	DeliveredDate   string  `json:"delivery_date"`
	UnitPricePerQty float64 `json:"unit_price_per_qty"`
	IssuesAgainst   string  `json:"issue_against"`
	Notes           string  `json:"notes"`
	Category        string  `json:"category"`
	Warranty        int     `json:"warranty"`
	WarrantyDueDays int     `json:"warranty_due_days"`
}

type OutwardDropDown struct {
	Seller        string `json:"seller"`
	BranchRegion  string `json:"branch_region"`
	IssuesAgainst string `json:"issue_against"`
	Category      string `json:"category"`
}
type MaterialOutwardInterface interface {
	FetchFormDropdownData() ([]OutwardDropDown, error)
	SubmitFormOutwardData(material MaterialOutward) error
	FetchAllFormOutwardData(r *http.Request) ([]MaterialOutward, error)
	UpdateFormOutwardData(material MaterialOutward) error
}
