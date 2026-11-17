package favorites

import (
	"net/http"
	"testing"

	apiRunner "api-tests-template/internal/helpers/api-runner"
)

const favoritesPath = "/favorites"

// Add добавляет объявление либо отправляет невалидное тело для негативного сценария.
func Add(t *testing.T, token string, payload any) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Post(favoritesPath).ContentType("application/json").JSON(payload).Expect(t).End().Response
}

// List получает страницу избранного с поиском по уникальному названию.
func List(t *testing.T, token, search string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(favoritesPath).Query("search", search).Query("page_size", "50").Expect(t).End().Response
}

// Delete удаляет запись избранного по идентификатору объявления.
func Delete(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(favoritesPath + "/" + id).Expect(t).End().Response
}
