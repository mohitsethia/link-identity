package http_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/link-identity/app/domain"
	httpHandler "github.com/link-identity/app/http"
	mockObject "github.com/link-identity/app/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testStruct struct {
	IsCalled bool
	Response interface{}
	Error    error
}

// TestLinkIdentityHandler_Identify ...
func TestLinkIdentityHandler_Identify(t *testing.T) {
	tests := []struct {
		Name               string
		RequestPayload     *httpHandler.RequestDTO
		ExpectedResponse   string
		ExpectedStatusCode int
		Service            testStruct
	}{
		{
			Name: "Happy path",
			RequestPayload: &httpHandler.RequestDTO{
				Email: "test1@gmail.com",
				Phone: "+4917611111111",
			},
			ExpectedResponse: `{
				"status_code": 200,
				"data": {
					"contact": {
						"PrimaryContactID": 1,
						"emails": [
							"test1@gmail.com"
						],
						"phoneNumbers": [
							"+4917611111111"
						],
						"secondaryContactIds": null
					}
    			}}`,
			Service: testStruct{
				IsCalled: true,
				Response: []*domain.Contact{
					{
						ContactID:        1,
						Email:            sql.NullString{String: "test1@gmail.com", Valid: true},
						Phone:            sql.NullString{String: "+4917611111111", Valid: true},
						LinkedPrecedence: "primary",
					},
				},
				Error: nil,
			},
			ExpectedStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ctx := context.Background()

			serviceMock := new(mockObject.LinkIdentityServiceMock)
			if tt.Service.IsCalled {
				if tt.Service.Response == nil {
					tt.Service.Response = ([]*domain.Contact)(nil)
				}
				serviceMock.On("Identify", ctx, mock.Anything, mock.Anything).
					Return(tt.Service.Response, tt.Service.Error)
			}

			handler := httpHandler.NewLinkIdentityHandler(serviceMock)

			jsonPayload, _ := json.Marshal(tt.RequestPayload)
			if tt.RequestPayload == nil {
				jsonPayload, _ = json.Marshal("random")
			}
			req, err := http.NewRequest("POST", "/identify", bytes.NewBuffer(jsonPayload))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			AuthorizedCustomerHandler := http.HandlerFunc(handler.Identify)
			AuthorizedCustomerHandler.ServeHTTP(rr, req)

			var body map[string]interface{}
			json.NewDecoder(rr.Body).Decode(&body)
			assert.Equal(t, tt.ExpectedStatusCode, rr.Code)
		})
	}
}
