package homework

import (
	"net/http"
	"strings"
	"testing"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementInvalidParameters(t *testing.T) {
	owner := fixtures.Login(t, "TEST")
	cases := []struct {
		name, field, value, message string
		omit                        bool
	}{
		{name: "missing_title", field: "title", omit: true, message: "title is required"},
		{name: "empty_title", field: "title", message: "title is required"},
		{name: "long_title", field: "title", value: strings.Repeat("x", fixtures.TitleLimit+1), message: "title must not exceed 50 characters"},
		{name: "missing_description", field: "description", omit: true, message: "description is required"},
		{name: "empty_description", field: "description", message: "description is required"},
		{name: "long_description", field: "description", value: strings.Repeat("x", 501), message: "description must be between 1 and 500 characters"},
		{name: "negative_price", field: "price", value: "-1", message: "price must be between 0 and 1000000 inclusive"},
		{name: "excessive_price", field: "price", value: "1000001", message: "price must be between 0 and 1000000 inclusive"},
		{name: "invalid_price", field: "price", value: "invalid", message: "cannot convert string value invalid to float"},
		{name: "zero_quantity", field: "quantity", value: "0", message: "quantity must be between 1 and 100 inclusive"},
		{name: "excessive_quantity", field: "quantity", value: "101", message: "quantity must be between 1 and 100 inclusive"},
		{name: "invalid_quantity", field: "quantity", value: "invalid", message: "cannot convert string value to int"},
		{name: "missing_photos", field: "photos", omit: true, message: "at least 1 photo is required"},
		{name: "too_many_photos", field: "photos", message: "3"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Подготовка: невалидный параметр %s", tc.field)
			data := fixtures.Request()
			if tc.field == "photos" {
				if tc.omit {
					data.Photos = nil
				} else {
					data.Photos = append(data.Photos, data.Photos[0], data.Photos[0], data.Photos[0])
				}
			} else if tc.omit {
				delete(data.Fields, tc.field)
			} else {
				data.Fields[tc.field] = tc.value
			}
			response.RequireError(t, advertisement.Create(t, owner.Token, data), http.StatusBadRequest, response.InvalidRequest, tc.message)
		})
	}
}
