package homework

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesUnauthorized(t *testing.T) {
	owner := fixtures.Login(t, "TEST")
	ad := fixtures.Advertisement(t, owner.Token)
	for _, tc := range fixtures.AuthorizationCases() {
		t.Run(tc.Name, func(t *testing.T) {
			t.Run("list", func(t *testing.T) {
				response.RequireError(t, favorites.ListRaw(t, tc.Token, ""), http.StatusUnauthorized, response.Unauthorized, tc.Message)
			})
			t.Run("add", func(t *testing.T) {
				response.RequireError(t, favorites.AddID(t, tc.Token, ad.ID), http.StatusUnauthorized, response.Unauthorized, tc.Message)
			})
			t.Run("delete", func(t *testing.T) {
				response.RequireError(t, favorites.Delete(t, tc.Token, ad.ID), http.StatusUnauthorized, response.Unauthorized, tc.Message)
			})
		})
	}
}
