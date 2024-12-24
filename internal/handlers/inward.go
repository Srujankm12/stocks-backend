package handlers

import (
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

func (ic *InwardController) FetchInwardData(w http.ResponseWriter, r *http.Request) {
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
