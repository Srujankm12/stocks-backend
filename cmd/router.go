package main

import (
	"database/sql"

	"github.com/Srujankm12/SRstocks/internal/handlers"
	"github.com/Srujankm12/SRstocks/internal/middlewares"
	"github.com/Srujankm12/SRstocks/repository"
	"github.com/gorilla/mux"
)

func registerRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.CorsMiddleware)
	inwardcon := handlers.NewInwardController(repository.NewMaterialInwardRepo(db))
	router.HandleFunc("/inward", inwardcon.FetchInwardData).Methods("GET")

	return router
}
