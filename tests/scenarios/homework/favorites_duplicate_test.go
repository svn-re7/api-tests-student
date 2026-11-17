package homework

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

func TestFavoritesRejectDuplicate(t *testing.T) {
	owner, viewer := fixtures.Users(t)
	ad := fixtures.Advertisement(t, owner.Token)
	first := favorites.AddAdvertisement(t, viewer.Token, ad.ID)
	response.RequireError(t, favorites.AddID(t, viewer.Token, ad.ID), http.StatusConflict, response.AlreadyExists, "Advertisement is already in favorites")
	page := favorites.List(t, viewer.Token, ad.Title)
	require.Equal(t, 1, page.Total, "Повторное добавление изменило количество записей")
	require.Len(t, page.Items, 1, "В списке появился дубликат")
	require.Equal(t, first.ID, page.Items[0].ID, "Изменилась исходная запись избранного")
	t.Log("Проверка: повторное добавление отклонено, запись сохранилась в одном экземпляре")
}
