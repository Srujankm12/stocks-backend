package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Srujankm12/SRstocks/internal/handlers"
	"github.com/Srujankm12/SRstocks/internal/middlewares"
	"github.com/Srujankm12/SRstocks/repository"
	"github.com/gorilla/mux"
)

func registerRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.CorsMiddleware)

	inwardcon := handlers.NewInwardController(repository.NewMaterialInwardRepo(db))
	router.HandleFunc("/inward", inwardcon.FetchInwardDataController).Methods("GET")
	router.HandleFunc("/submit", inwardcon.SubmitInwardDataController).Methods("POST")
	router.HandleFunc("/getlist", inwardcon.FetchFormInwardDataController).Methods("GET")

	outwardcon := handlers.NewOutwardController(repository.NewMaterialOutwardRepo(db))
	router.HandleFunc("/outwardDropdown", outwardcon.FetchOutwardDataController).Methods("GET")
	router.HandleFunc("/submitoutward", outwardcon.SubmitOutwardDataController).Methods("POST")
	router.HandleFunc("/fetchoutward", outwardcon.FetchFormOutwardDataController).Methods("GET")

	materialstockcon := handlers.NewMaterialStockController(repository.NewMaterialStockRepo(db))
	router.HandleFunc("/materialstockdropdown", materialstockcon.FetchMaterialDropdownDataController).Methods("GET")
	router.HandleFunc("/materialstock", materialstockcon.SubmitMaterialStockController).Methods("POST")
	router.HandleFunc("/materialstockdata", materialstockcon.FetchAllMaterialStockController).Methods("GET")
	router.HandleFunc("/materialupdate", materialstockcon.UpdateMaterialStockController).Methods("POST")

	excelcon := handlers.NewExcelDownloadController(repository.NewExcelDownloadMSRepo(db))
	router.HandleFunc("/materialstockdownload", excelcon.DownloadMaterialStock).Methods("GET")
	router.HandleFunc("/Users/bunny/Desktop/Finalyear/test", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileName := vars["filename"]
		filePath := fmt.Sprintf("/Users/bunny/Desktop/Finalyear/test%s", fileName) // Update if your temp directory is different

		http.ServeFile(w, r, filePath)
	}).Methods("GET")

	excelconin := handlers.NewExcelDownloadMIController(repository.NewExcelDownloadMIRepo(db))
	router.HandleFunc("/downloadinward", excelconin.DownloadMaterialInward).Methods("GET")
	router.HandleFunc("/Users/bunny/Desktop/Finalyear/test", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileName := vars["filename"]
		filePath := fmt.Sprintf("/Users/bunny/Desktop/Finalyear/test%s", fileName)
		http.ServeFile(w, r, filePath)
	}).Methods("GET")

	excelconout := handlers.NewExcelDownloadMOController(repository.NewExcelDownloadMORepo(db))
	router.HandleFunc("/downloadoutward", excelconout.DownloadMaterialOutward).Methods("GET")
	router.HandleFunc("/Users/bunny/Desktop/Finalyear/test", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileName := vars["filename"]
		filePath := fmt.Sprintf("/Users/bunny/Desktop/Finalyear/test%s", fileName)
		http.ServeFile(w, r, filePath)
	}).Methods("GET")
	return router
}
