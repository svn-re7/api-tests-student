package homework

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

func TestFavoritesRejectAboveLimit(t *testing.T) {
	owner, viewer := fixtures.Users(t)
	initial := favorites.List(t, viewer.Token, "")
	require.Zero(t, initial.Total, "Для проверки лимита нужен отдельный FAVORITES-аккаунт с пустым избранным")
	for range favorites.MaxFavorites {
		ad := fixtures.Advertisement(t, owner.Token)
		favorites.AddAdvertisement(t, viewer.Token, ad.ID)
	}
	page := favorites.List(t, viewer.Token, "")
	require.Equal(t, favorites.MaxFavorites, page.Total, "Не подготовлено 50 записей")
	require.Len(t, page.Items, favorites.MaxFavorites, "Не получена полная страница из 50 записей")
	extra := fixtures.Advertisement(t, owner.Token)
	response.RequireError(t, favorites.AddID(t, viewer.Token, extra.ID), http.StatusBadRequest, response.LimitReached, "50")
	require.Equal(t, favorites.MaxFavorites, favorites.List(t, viewer.Token, "").Total, "Отклоненное добавление изменило лимит")
	favorites.RequireMembership(t, viewer.Token, extra, false)
}
