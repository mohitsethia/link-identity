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
	var linkedId uint
	var secondaryIds []uint
	var secondaryEmails []string
	var secondaryPhones []string
	var primaryEmail string
	var primaryPhone string

	for _, v := range contacts {
		if v.LinkedPrecedence == "primary" {
			linkedId = v.ContactId
			primaryEmail = v.Email.String
			primaryPhone = v.Phone.String
		} else {
			secondaryIds = append(secondaryIds, v.ContactId)
			secondaryEmails = append(secondaryEmails, v.Email.String)
			secondaryPhones = append(secondaryPhones, v.Phone.String)
		}
	}

	secondaryEmails = remove(secondaryEmails, primaryEmail)
	secondaryPhones = remove(secondaryPhones, primaryPhone)

	return &ResponseDTO{Contact: struct {
		PrimaryContactID    uint     `json:"PrimaryContactID"`
		Emails              []string `json:"emails"`
		PhoneNumbers        []string `json:"phoneNumbers"`
		SecondaryContactIds []uint   `json:"secondaryContactIds"`
	}{
		PrimaryContactID:    linkedId,
		Emails:              append([]string{primaryEmail}, secondaryEmails...),
		PhoneNumbers:        append([]string{primaryPhone}, secondaryPhones...),
		SecondaryContactIds: secondaryIds,
	}}
}

func remove(slice []string, val string) []string {
	var arr []string
	for _, v := range slice {
		if v != val {
			arr = append(arr, v)
		}
	}
	return arr
}
