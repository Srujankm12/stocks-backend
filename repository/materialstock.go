package repository

import (
	"database/sql"
	"net/http"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/Srujankm12/SRstocks/pkg/database"
)

type MaterialStockRepo struct {
	db *sql.DB
}

func NewMaterialStockRepo(db *sql.DB) *MaterialStockRepo {
	return &MaterialStockRepo{
		db: db,
	}
}

func (msr *MaterialStockRepo) SubmitMaterialStock(material models.MaterialStock) error {
	query := database.NewQuery(msr.db)
	err := query.SubmitMaterialStock(material)
	if err != nil {
		return err
	}
	return nil
}
func (msr *MaterialStockRepo) FetchAllMaterialStock(r *http.Request) ([]models.MaterialStock, error) {
	query := database.NewQuery(msr.db)
	formData, err := query.FetchAllMaterialStock()
	if err != nil {
		return nil, err
	}
	return formData, nil
}
func (msr *MaterialStockRepo) UpdateMaterialStock(material models.MaterialStock) error {
	query := database.NewQuery(msr.db)
	err := query.UpdateMaterialStock(material)
	if err != nil {
		return err
	}
	return nil
}

func (msr *MaterialStockRepo) FetchMaterialDropdownData() ([]models.MaterialStockDropDown, error) {
	query := database.NewQuery(msr.db)
	formData, err := query.FetchMaterialDropdownData()
	if err != nil {
		return nil, err
	}
	return formData, nil
}
