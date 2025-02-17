package repository

import (
	"database/sql"
	"net/http"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/Srujankm12/SRstocks/pkg/database"
)

type MaterialOutwardRepo struct {
	db *sql.DB
}

func NewMaterialOutwardRepo(db *sql.DB) *MaterialOutwardRepo {
	return &MaterialOutwardRepo{
		db: db,
	}
}

func (mir *MaterialOutwardRepo) FetchFormDropdownData() ([]models.OutwardDropDown, error) {
	query := database.NewQuery(mir.db)
	res, err := query.FetchFormOutwardDropDown()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return []models.OutwardDropDown{}, nil
	}
	return res, nil
}

func (mir *MaterialOutwardRepo) SubmitFormOutwardData(material models.MaterialOutward) error {
	query := database.NewQuery(mir.db)
	err := query.SubmitFormOutwardData(material)
	if err != nil {
		return err
	}
	return nil
}
func (mir *MaterialOutwardRepo) FetchAllFormOutwardData(r *http.Request) ([]models.MaterialOutward, error) {
	query := database.NewQuery(mir.db)
	formData, err := query.FetchAllFormOutwardData()
	if err != nil {
		return nil, err
	}
	return formData, nil
}

func (mir *MaterialOutwardRepo) UpdateFormOutwardData(material models.MaterialOutward) error {
	query := database.NewQuery(mir.db)
	err := query.UpdateMaterialOutward(material)
	if err != nil {
		return err
	}
	return nil
}
