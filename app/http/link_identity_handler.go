package http

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/link-identity/app/application"
	"github.com/link-identity/app/domain"
	"github.com/link-identity/app/utils"

	"github.com/nyaruka/phonenumbers"
)

type (
	// LinkIdentityHandler ...
	LinkIdentityHandler struct {
		service application.LinkIdentityService
	}

	// RequestDTO ...
	RequestDTO struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	// ResponseDTO ...
	ResponseDTO struct {
		Contact struct {
			PrimaryContactID    uint     `json:"PrimaryContactID"`
			Emails              []string `json:"emails"`
			PhoneNumbers        []string `json:"phoneNumbers"`
			SecondaryContactIds []uint   `json:"secondaryContactIds"`
		} `json:"contact"`
	}
)

// NewLinkIdentityHandler ...
func NewLinkIdentityHandler(service application.LinkIdentityService) *LinkIdentityHandler {
	return &LinkIdentityHandler{
		service: service,
	}
}

// Identify ...
func (h *LinkIdentityHandler) Identify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	model := new(RequestDTO)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&model)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		utils.ResponseJSON(w, http.StatusBadRequest, resp)
		return
	}

	if v := model.Validate(); v != nil {
		utils.ResponseJSON(w, v.StatusCode, v)
		return
	}

	contacts, err := h.service.Identify(ctx, model.Email, model.Phone)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		utils.ResponseJSON(w, http.StatusInternalServerError, resp)
		return
	}

	resp := utils.ResponseSuccess(http.StatusOK, convertContactsToResponseDTO(contacts))
	utils.ResponseJSON(w, http.StatusOK, resp)

	return
}

// Validate ...
func (v *RequestDTO) Validate() *utils.ErrorResponse {
	if v.Email != "" {
		if _, err := mail.ParseAddress(v.Email); err != nil {
			return utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		}
	}

	if v.Phone == "" {
		return nil
	}

	num, err := phonenumbers.Parse(v.Phone, "")
	if err != nil {
		return utils.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	if !phonenumbers.IsValidNumber(num) {
		return utils.NewErrorResponse(http.StatusBadRequest, "Invalid phone number")
	}

	return nil
}

func convertContactsToResponseDTO(contacts []*domain.Contact) *ResponseDTO {
	var primaryContactID uint
	var secondaryIds []uint
	secondaryEmails := make(map[string]bool)
	secondaryPhones := make(map[string]bool)
	var primaryEmail string
	var primaryPhone string

	for _, v := range contacts {
		if v.LinkedPrecedence == "primary" {
			primaryContactID = v.ContactID
			primaryEmail = v.Email.String
			primaryPhone = v.Phone.String
		} else {
			secondaryIds = append(secondaryIds, v.ContactID)
			secondaryEmails[v.Email.String] = true
			secondaryPhones[v.Phone.String] = true
		}
	}

	delete(secondaryEmails, primaryEmail)
	delete(secondaryPhones, primaryPhone)

	return &ResponseDTO{Contact: struct {
		PrimaryContactID    uint     `json:"PrimaryContactID"`
		Emails              []string `json:"emails"`
		PhoneNumbers        []string `json:"phoneNumbers"`
		SecondaryContactIds []uint   `json:"secondaryContactIds"`
	}{
		PrimaryContactID:    primaryContactID,
		Emails:              append([]string{primaryEmail}, convertMapToArray(secondaryEmails)...),
		PhoneNumbers:        append([]string{primaryPhone}, convertMapToArray(secondaryPhones)...),
		SecondaryContactIds: secondaryIds,
	}}
}

func convertMapToArray(m map[string]bool) []string {
	var arr []string
	for k := range m {
		arr = append(arr, k)
	}
	return arr
}
