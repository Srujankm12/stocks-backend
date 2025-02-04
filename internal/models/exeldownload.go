package models

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type ExcelDownloadMS struct {
	ID               int     `json:"id"`
	Timestamp        string  `json:"timestamp"`
	Supplier         string  `json:"supplier"`
	Category         string  `json:"category"`
	LeadTime         int     `json:"lead_time"`
	StdNonStd        string  `json:"std_non_std"`
	PartCode         string  `json:"part_code"`
	Unit             string  `json:"unit"`
	Rate             float64 `json:"rate"`
	MinimumRetain    int     `json:"minimum_retain"`
	MaximumRetain    int     `json:"maximum_retain"`
	Received         int     `json:"received"`
	Issue            int     `json:"issue"`
	ReservedStock    int     `json:"reserved_stock"`
	Stock            int     `json:"stock"`
	Value            float64 `json:"value"`
	ReorderStatus    bool    `json:"reorder_status"`
	ExcessStock      int     `json:"excess_stock"`
	ExcessStockValue float64 `json:"excess_stock_value"`
}

type ExcelDownloadMSInterface interface {
	CreateMaterialStock(*io.ReadCloser) (*excelize.File, error)
}
