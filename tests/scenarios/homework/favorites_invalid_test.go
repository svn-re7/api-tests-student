package homework

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesInvalidAdvertisementID(t *testing.T) {
	viewer := fixtures.Login(t, "FAVORITES")
	cases := []struct {
		name    string
		payload map[string]string
		message string
	}{
		{name: "missing_id", payload: map[string]string{}, message: "AdvertisementID"},
		{name: "empty_id", payload: map[string]string{"advertisement_id": ""}, message: "invalid UUID length: 0"},
		{name: "invalid_id", payload: map[string]string{"advertisement_id": "invalid"}, message: "invalid UUID"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			response.RequireError(t, favorites.Add(t, viewer.Token, tc.payload), http.StatusBadRequest, response.InvalidRequest, tc.message)
		})
	}
}
