package main

import (
	"database/sql"
	"net/http"
	"os"

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

	excelconin := handlers.NewExcelDownloadMIController(repository.NewExcelDownloadMIRepo(db))
	router.HandleFunc("/downloadinward", excelconin.DownloadMaterialInward).Methods("GET")

	excelconout := handlers.NewExcelDownloadMOController(repository.NewExcelDownloadMORepo(db))
	router.HandleFunc("/downloadoutward", excelconout.DownloadMaterialOutward).Methods("GET")

	mat := handlers.NewMaterialHandler(repository.NewMaterialRepository(db))
	router.HandleFunc("/material", mat.GetStockHandler).Methods("POST")

	tempDir := "/tmp"
	if os.Getenv("OS") == "Windows_NT" {
		tempDir = os.Getenv("TEMP")
	}

	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir(tempDir)))).Methods("GET")

	return router
}
