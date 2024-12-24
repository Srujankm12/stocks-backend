package repository

import (
	"database/sql"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/Srujankm12/SRstocks/pkg/database"
)

type MaterialInwardRepo struct {
	db *sql.DB
}

func NewMaterialInwardRepo(db *sql.DB) *MaterialInwardRepo {
	return &MaterialInwardRepo{
		db: db,
	}
}

func (mir *MaterialInwardRepo) FetchFormData() ([]models.InwardDropDown, error) {
	query := database.NewQuery(mir.db)
	res, err := query.FetchFormData()
	if err != nil {
		return nil, err
	}

	// Ensure we don't return nil for empty data
	if len(res) == 0 {
		return []models.InwardDropDown{}, nil
	}
	return res, nil
}
