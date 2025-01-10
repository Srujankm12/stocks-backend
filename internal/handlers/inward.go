package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/Srujankm12/SRstocks/pkg/utils"
)

type InwardController struct {
	inwardRepo models.MaterialInwardInterface
}

func NewInwardController(inwardRepo models.MaterialInwardInterface) *InwardController {
	return &InwardController{
		inwardRepo: inwardRepo,
	}
}

func (ic *InwardController) FetchInwardDataController(w http.ResponseWriter, r *http.Request) {
	inwardData, err := ic.inwardRepo.FetchFormData()
	if err != nil {
		log.Printf("Error fetching inward data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Internal Server Error"})
		return
	}

	if len(inwardData) == 0 {
		log.Println("No inward data found.")
		w.WriteHeader(http.StatusOK)
		utils.Encode(w, []models.InwardDropDown{})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Encode(w, inwardData)
}

func (ic *InwardController) SubmitInwardDataController(w http.ResponseWriter, r *http.Request) {
	var material models.MaterialInward
	err := json.NewDecoder(r.Body).Decode(&material)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Invalid request body"})
		return
	}

	err = ic.inwardRepo.SubmitFormData(material)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Encode(w, map[string]string{"message": "success"})
}

func (ic *InwardController) FetchFormInwardDataController(w http.ResponseWriter, r *http.Request) {
	inwardData, err := ic.inwardRepo.FetchAllFormDataInward(r)
	if err != nil {
		log.Printf("Error fetching all inward data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Internal Server Error"})
		return
	}

	if len(inwardData) == 0 {
		log.Println("No inward form data found.")
		w.WriteHeader(http.StatusOK)
		utils.Encode(w, []models.MaterialInward{})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Encode(w, inwardData)
}
