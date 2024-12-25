package repository

import (
	"database/sql"
	"net/http"

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
	if len(res) == 0 {
		return []models.InwardDropDown{}, nil
	}
	return res, nil
}

func (mir *MaterialInwardRepo) SubmitFormData(material models.MaterialInward) error {
	query := database.NewQuery(mir.db)
	err := query.SubmitFormData(material)
	if err != nil {
		return err
	}
	return nil
}

func (mir *MaterialInwardRepo) FetchAllFormDataInward(r *http.Request) ([]models.MaterialInward, error) {
	query := database.NewQuery(mir.db)
	formData, err := query.FetchAllFormData()
	if err != nil {
		return nil, err
	}
	return formData, nil
}
