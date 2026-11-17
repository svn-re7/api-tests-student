package homework

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementUnauthorized(t *testing.T) {
	fixtures.Environment(t)
	for _, tc := range fixtures.AuthorizationCases() {
		t.Run(tc.Name, func(t *testing.T) {
			t.Log("Подготовка: валидная форма и проверяемый вариант токена")
			response.RequireError(t, advertisement.Create(t, tc.Token, fixtures.Request()), http.StatusUnauthorized, response.Unauthorized, tc.Message)
		})
	}
}
