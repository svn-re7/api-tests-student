package favorites

import (
	"net/http"
	"testing"

	client "api-tests-template/internal/client/http/favorites"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/response"
	"github.com/stretchr/testify/require"
)

const MaxFavorites = 50

type Item struct {
	ID            string                 `json:"favorite_id"`
	Advertisement adModels.Advertisement `json:"advertisement"`
}

type Page struct {
	Items []Item `json:"items"`
	Total int    `json:"total"`
}

// Add отправляет тело запроса на добавление в избранное.
func Add(t *testing.T, token string, payload any) response.Result {
	t.Helper()
	t.Logf("Шаг: добавление в избранное %v", payload)
	return response.Read(t, client.Add(t, token, payload))
}

// AddAdvertisement проверяет создание избранного и его связь с объявлением.
func AddAdvertisement(t *testing.T, token, id string) Item {
	t.Helper()
	item := response.Decode[Item](t, Add(t, token, map[string]string{"advertisement_id": id}), http.StatusCreated)
	require.NotEmpty(t, item.ID, "Пустой favorite_id")
	require.Equal(t, id, item.Advertisement.ID, "В избранное добавлено другое объявление")
	return item
}

// ListRaw возвращает ответ списка для проверки авторизации.
func ListRaw(t *testing.T, token, search string) response.Result {
	t.Helper()
	return response.Read(t, client.List(t, token, search))
}

// List возвращает типизированную страницу избранного.
func List(t *testing.T, token, search string) Page {
	t.Helper()
	t.Logf("Шаг: получение избранного, search=%q", search)
	return response.Decode[Page](t, ListRaw(t, token, search), http.StatusOK)
}

// Delete возвращает ответ удаления для позитивных и негативных сценариев.
func Delete(t *testing.T, token, id string) response.Result {
	t.Helper()
	t.Logf("Шаг: удаление из избранного %s", id)
	return response.Read(t, client.Delete(t, token, id))
}

// Remove проверяет успешное удаление записи избранного.
func Remove(t *testing.T, token, id string) {
	t.Helper()
	result := Delete(t, token, id)
	require.Equal(t, http.StatusNoContent, result.Status, "Удаление из избранного: %s", result.Body)
}

// RequireMembership проверяет состав отфильтрованного списка по точному id.
func RequireMembership(t *testing.T, token string, ad adModels.Advertisement, present bool) {
	t.Helper()
	page := List(t, token, ad.Title)
	ids := make([]string, 0, len(page.Items))
	for _, item := range page.Items {
		ids = append(ids, item.Advertisement.ID)
	}
	if present {
		require.Contains(t, ids, ad.ID, "Объявление отсутствует в избранном")
	} else {
		require.NotContains(t, ids, ad.ID, "Объявление осталось в избранном")
	}
	t.Logf("Проверка: объявление %s, присутствие в избранном=%t", ad.ID, present)
}
