package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Srujankm12/SRstocks/internal/models"
	"github.com/Srujankm12/SRstocks/pkg/utils"
)

type OutwardController struct {
	outwardRepo models.MaterialOutwardInterface
}

func NewOutwardController(outwardRepo models.MaterialOutwardInterface) *OutwardController {
	return &OutwardController{
		outwardRepo: outwardRepo,
	}
}
func (oc *OutwardController) FetchOutwardDataController(w http.ResponseWriter, r *http.Request) {
	outwardData, err := oc.outwardRepo.FetchFormDropdownData()
	if err != nil {
		log.Printf("Error fetching outward data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Internal Server Error"})
		return
	}

	if len(outwardData) == 0 {
		log.Println("No outward data found.")
		w.WriteHeader(http.StatusOK)
		utils.Encode(w, []models.OutwardDropDown{})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Encode(w, outwardData)
}
func (oc *OutwardController) SubmitOutwardDataController(w http.ResponseWriter, r *http.Request) {
	var material models.MaterialOutward
	err := json.NewDecoder(r.Body).Decode(&material)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Invalid request body"})
		return
	}

	err = oc.outwardRepo.SubmitFormOutwardData(material)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Error submitting outward data"})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Encode(w, map[string]string{"message": "Outward data submitted successfully"})
}

func (ic *OutwardController) FetchFormOutwardDataController(w http.ResponseWriter, r *http.Request) {
	outwardData, err := ic.outwardRepo.FetchAllFormOutwardData(r)
	if err != nil {
		log.Printf("Error fetching all outward data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Internal Server Error"})
		return
	}

	if len(outwardData) == 0 {
		log.Println("No outward form data found.")
		w.WriteHeader(http.StatusOK)
		utils.Encode(w, []models.MaterialOutward{})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Encode(w, outwardData)
}
