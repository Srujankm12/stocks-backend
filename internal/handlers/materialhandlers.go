package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Srujankm12/SRstocks/repository"
)

type MaterialHandler struct {
	Repo *repository.MaterialRepository
}

func NewMaterialHandler(repo *repository.MaterialRepository) *MaterialHandler {
	return &MaterialHandler{Repo: repo}
}

func (h *MaterialHandler) GetStockHandler(w http.ResponseWriter, r *http.Request) {
	stock, err := h.Repo.GetMaterialStock()
	if err != nil {
		http.Error(w, "Failed to fetch stock data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}
