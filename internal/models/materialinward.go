package models

import "net/http"

type MaterialInward struct {
	Timestamp       string  `json:"timestamp"`
	Supplier        string  `json:"supplier"`
	Buyer           string  `json:"buyer"`
	PartCode        string  `json:"partcode"`
	SerialNumber    string  `json:"serial_number"`
	Quantity        int     `json:"qty"`
	PONo            string  `json:"po_no"`
	PODate          string  `json:"po_date"`
	InvoiceNo       string  `json:"invoice_no"`
	InvoiceDate     string  `json:"invoice_date"`
	ReceivedDate    string  `json:"received_date"`
	UnitPricePerQty float64 `json:"unit_price_per_qty"`
	Category        string  `json:"category"`
	Warranty        int     `json:"warranty"`
	WarrantyDueDays int     `json:"warranty_due_days"`
}

type InwardDropDown struct {
	Supplier string `json:"supplier"`
	Buyer    string `json:"buyer"`
}

type MaterialInwardInterface interface {
	FetchFormData() ([]InwardDropDown, error)
	SubmitFormData(material MaterialInward) error
	FetchAllFormDataInward(r *http.Request) ([]MaterialInward, error)
}
