package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/toleubekov/check-iin-kaz/internal/model"
	"github.com/toleubekov/check-iin-kaz/internal/repository"
	"github.com/toleubekov/check-iin-kaz/internal/service"
)

type Handler struct {
	iinService *service.IINService
	repo       *repository.PersonRepository
}

func NewHandler(iinService *service.IINService, repo *repository.PersonRepository) *Handler {
	return &Handler{
		iinService: iinService,
		repo:       repo,
	}
}

func (h *Handler) CheckIIN(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	iin := vars["iin"]

	correct, sex, dateOfBirth, err := h.iinService.ValidateIIN(iin)

	response := model.IINResponse{
		Correct: correct,
	}

	if correct {
		response.Sex = sex
		response.DateOfBirth = dateOfBirth
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// Modify the CreatePerson function in internal/api/handler.go

func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	log.Println("Received CreatePerson request")

	var person model.Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	log.Printf("Attempting to create person: Name=%s, IIN=%s", person.Name, person.IIN)

	correct, _, _, err := h.iinService.ValidateIIN(person.IIN)
	if !correct || err != nil {
		errorMsg := "Invalid IIN"
		if err != nil {
			errorMsg = err.Error()
		}
		log.Printf("ERROR: IIN validation failed: %s", errorMsg)
		sendErrorResponse(w, http.StatusInternalServerError, errorMsg)
		return
	}

	log.Printf("IIN validated successfully, proceeding to database insertion")

	err = h.repo.Create(&person)
	if err != nil {
		log.Printf("ERROR: Database insertion failed: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("SUCCESS: Person created in database with IIN: %s", person.IIN)

	response := model.PersonResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetPersonByIIN(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	iin := vars["iin"]

	correct, _, _, err := h.iinService.ValidateIIN(iin)
	if !correct || err != nil {
		errorMsg := "Invalid IIN"
		if err != nil {
			errorMsg = err.Error()
		}
		sendErrorResponse(w, http.StatusInternalServerError, errorMsg)
		return
	}

	person, err := h.repo.GetByIIN(iin)
	if err != nil {
		if err.Error() == "person not found" {
			sendErrorResponse(w, http.StatusNotFound, "Person not found")
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

func (h *Handler) FindPeopleByNamePart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namePart := vars["name_part"]

	people, err := h.repo.FindByNamePart(namePart)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := model.PersonResponse{
		Success: false,
		Errors:  errorMsg,
	}

	json.NewEncoder(w).Encode(response)
}
