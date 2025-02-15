package models

type MaterialStockUpdated struct {
	PartCode      string  `json:"part_code"`
	Received      int     `json:"received"`
	Issued        int     `json:"issued"`
	Stock         int     `json:"stock"`
	ReorderStatus bool    `json:"reorder_status"`
	StockValue    float64 `json:"stock_value"`
}
