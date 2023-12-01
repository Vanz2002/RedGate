package handlers

import (
	"errors"
	"log"
	"net/http"

	"redgate.com/b/db/sqlc"
	"redgate.com/b/utils"
)

func NewPlateIDHandler(l *log.Logger, q *sqlc.Queries, u *AuthedUser) *PlateHandler {
	var c uint = 0
	return &PlateHandler{&Handler{l, q, &c, u}}
}

func (plate_h *PlateHandler) CreatePlateHandler(w http.ResponseWriter, r *http.Request) {
	hp := HandlerParam{w, r, http.MethodPost, plate_h.createPlateId}
	plate_h.h.handleRequest(hp, plate_h.h.u)
}

func (plate_h *PlateHandler) VerifyPlateHandler(w http.ResponseWriter, r *http.Request) {
	hp := HandlerParam{w, r, http.MethodPost, plate_h.verifyPlateId}
	plate_h.h.handleRequest(hp, nil)
}

func (ph *PlateHandler) createPlateId(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return err
	}

	// Retrieve form values
	accountID := r.FormValue("account_id")
	plateNumber := r.FormValue("plate_number")
	if !utils.PlateNumberIsValid(plateNumber) {
		http.Error(w, "Invalid plate number", http.StatusInternalServerError)
		return errors.New("Invalid plate number")
	}
	vID := utils.GeneratePlateID(plateNumber)

	// Create cardParams using retrieved form values
	vehicleParams := sqlc.InsertVehicleParams{
		VID:         vID,
		AccountID:   utils.StringToNullString(accountID),
		PlateNumber: utils.StringToNullString(plateNumber),
	}

	veh, err := ph.h.q.InsertVehicle(r.Context(), vehicleParams)
	if err != nil {
		http.Error(w, "Error inserting data", http.StatusInternalServerError)
		return err
	}

	// after creating ID assume that user is also subscribing
	_, err = ph.h.q.AccountSubscribe(r.Context(), accountID)
	if err != nil {
		http.Error(w, "Error subscribing account", http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	toJSON(w, veh)
	return nil
}

func (ph *PlateHandler) verifyPlateId(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return err
	}
	// check if the plate id exist
	vID := r.FormValue("v_id")
	response, err := ph.h.q.VerifyVehicle(r.Context(), vID)
	if err != nil {
		http.Error(w, "Verifying failed", http.StatusInternalServerError)
		return err
	}

	plate_number := response.String
	if utils.GeneratePlateID(plate_number) != vID {
		http.Error(w, "Mismatch plate number with the hash value", http.StatusInternalServerError)
		return errors.New("Mismatch")
	}

	w.WriteHeader(http.StatusOK)
	toJSON(w, vID)
	return nil
}
