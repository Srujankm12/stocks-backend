package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/Srujankm12/SRstocks/pkg/utils"
)

type MaterialStockController struct {
	materialStockRepo models.MaterialStockInterface
}

func NewMaterialStockController(materialStockRepo models.MaterialStockInterface) *MaterialStockController {
	return &MaterialStockController{
		materialStockRepo: materialStockRepo,
	}
}
func (msc *MaterialStockController) SubmitMaterialStockController(w http.ResponseWriter, r *http.Request) {
	var material models.MaterialStock
	err := json.NewDecoder(r.Body).Decode(&material)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Invalid request body"})
		return
	}
	err = msc.materialStockRepo.SubmitMaterialStock(material)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Failed to submit material stock"})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, map[string]string{"message": "Material stock submitted successfully"})
}

func (msc *MaterialStockController) FetchAllMaterialStockController(w http.ResponseWriter, r *http.Request) {
	formData, err := msc.materialStockRepo.FetchAllMaterialStock(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Failed to fetch material stock"})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, formData)
}

func (msc *MaterialStockController) UpdateMaterialStockController(w http.ResponseWriter, r *http.Request) {
	var material models.MaterialStock
	err := json.NewDecoder(r.Body).Decode(&material)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Invalid request body"})
		return
	}
	err = msc.materialStockRepo.UpdateMaterialStock(material)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Failed to update material stock"})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, map[string]string{"message": "Material stock updated successfully"})
}
