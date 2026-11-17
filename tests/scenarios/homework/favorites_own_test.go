package homework

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesRejectOwnAdvertisement(t *testing.T) {
	owner := fixtures.Login(t, "TEST")
	ad := fixtures.Advertisement(t, owner.Token)
	response.RequireError(t, favorites.AddID(t, owner.Token, ad.ID), http.StatusForbidden, response.Forbidden, "Cannot add own advertisement to favorites")
	favorites.RequireMembership(t, owner.Token, ad, false)
}
