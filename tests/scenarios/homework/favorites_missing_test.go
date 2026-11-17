package homework

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesNonexistentAdvertisement(t *testing.T) {
	owner, viewer := fixtures.Users(t)
	id := fixtures.DeletedAdvertisementID(t, owner.Token)
	t.Run("add", func(t *testing.T) {
		response.RequireError(t, favorites.AddID(t, viewer.Token, id), http.StatusNotFound, response.NotFound, "Advertisement not found")
	})
	t.Run("delete", func(t *testing.T) {
		response.RequireError(t, favorites.Delete(t, viewer.Token, id), http.StatusNotFound, response.NotFound, "Favorite entry not found")
	})
}
