package http

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/link-identity/app/application"
	"github.com/link-identity/app/domain"
	"github.com/link-identity/app/utils"
	"net/http"
	"strconv"
)

type LocationHandler struct {
	locationService application.ILocation
}

func NewLocationHandler(locationService application.ILocation) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
	}
}

func (h *LocationHandler) GetLastNLocation(w http.ResponseWriter, r *http.Request) {
	lastN := r.URL.Query().Get("max")
	intLastN, ok := strconv.Atoi(lastN)
	if ok != nil || intLastN == 0 {
		intLastN = 5
	}
	rider := chi.URLParam(r, "rider")
	locations := h.locationService.GetLastNLocation(rider, intLastN)
	resp := utils.ResponseDTO{
		StatusCode: http.StatusOK,
		Data:       locations,
	}
	utils.ResponseJSON(w, http.StatusOK, resp)
}

func (h *LocationHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	//'{
	//"lat": 12.34,
	//"long": 56.78
	//}'
	//'localhost:8080/location/steve/now'
	rider := chi.URLParam(r, "rider")

	var location domain.Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		utils.ResponseJSON(w, http.StatusBadRequest, resp)
		return
	}
	h.locationService.UpdateLocation(rider, location)
	resp := utils.ResponseDTO{
		StatusCode: http.StatusOK,
		Data:       nil,
	}
	utils.ResponseJSON(w, http.StatusOK, resp)
}
