package homework

import (
	"testing"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/fixtures"
)

func TestFavoritesFullCycle(t *testing.T) {
	owner, viewer := fixtures.Users(t)
	ad := fixtures.Advertisement(t, owner.Token)
	favorites.RequireMembership(t, viewer.Token, ad, false)
	favorites.AddAdvertisement(t, viewer.Token, ad.ID)
	favorites.RequireMembership(t, viewer.Token, ad, true)
	favorites.Remove(t, viewer.Token, ad.ID)
	favorites.RequireMembership(t, viewer.Token, ad, false)
	favorites.AddAdvertisement(t, viewer.Token, ad.ID)
	favorites.RequireMembership(t, viewer.Token, ad, true)
	advertisement.Delete(t, owner.Token, ad.ID)
	favorites.RequireMembership(t, viewer.Token, ad, false)
}
