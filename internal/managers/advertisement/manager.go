package advertisement

import (
	"net/http"
	"testing"

	client "api-tests-template/internal/client/http/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/response"

	"github.com/stretchr/testify/require"
)

// Create отправляет объявление и регистрирует очистку любой успешно созданной записи.
func Create(t *testing.T, token string, data models.CreateRequest) response.Result {
	t.Helper()
	t.Logf("Шаг: создание объявления %q", data.Fields["title"])
	result := response.Read(t, client.Create(t, token, data))
	if result.Status == http.StatusCreated {
		ad := response.Decode[adModels.Advertisement](t, result, http.StatusCreated)
		require.NotEmpty(t, ad.ID, "Созданное объявление должно иметь id")
		t.Cleanup(func() {
			t.Logf("Очистка: удаление объявления %s", ad.ID)
			current := response.Read(t, client.Get(t, ad.ID))
			if current.Status == http.StatusNotFound {
				return
			}
			require.Equal(t, http.StatusOK, current.Status, "Проверка перед очисткой %s: %s", ad.ID, current.Body)
			deleted := response.Read(t, client.Delete(t, token, ad.ID))
			require.Equal(t, http.StatusNoContent, deleted.Status, "Очистка объявления %s: %s", ad.ID, deleted.Body)
		})
	}
	return result
}

// Get возвращает сохраненные данные объявления и ожидает HTTP 200.
func Get(t *testing.T, id string) adModels.Advertisement {
	t.Helper()
	t.Logf("Шаг: получение объявления %s", id)
	return response.Decode[adModels.Advertisement](t, response.Read(t, client.Get(t, id)), http.StatusOK)
}

// Delete проверяет успешное удаление объявления владельцем.
func Delete(t *testing.T, token, id string) {
	t.Helper()
	result := response.Read(t, client.Delete(t, token, id))
	require.Equal(t, http.StatusNoContent, result.Status, "Удаление %s: %s", id, result.Body)
	t.Logf("Проверка: объявление %s удалено", id)
}

// Photos возвращает все фотографии, доступные через API объявления.
func Photos(t *testing.T, token, id string) []adModels.Photo {
	t.Helper()
	t.Logf("Шаг: получение фотографий объявления %s", id)
	return response.Decode[[]adModels.Photo](t, response.Read(t, client.Photos(t, token, id)), http.StatusOK)
}
