package api

import (
	"github.com/gorilla/mux"
)

func SetupRouter(handler *Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/iin_check/{iin}", handler.CheckIIN).Methods("GET")

	r.HandleFunc("/people/info", handler.CreatePerson).Methods("POST")

	r.HandleFunc("/people/info/iin/{iin}", handler.GetPersonByIIN).Methods("GET")

	r.HandleFunc("/people/info/name/{name_part}", handler.FindPeopleByNamePart).Methods("GET")

	return r
}
