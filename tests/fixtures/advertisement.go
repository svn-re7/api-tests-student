package fixtures

import (
	_ "embed"
	"net/http"
	"os"
	"testing"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/response"
	"api-tests-template/internal/utils"

	"github.com/stretchr/testify/require"
)

const (
	InvalidToken = "invalid_token"
	Description  = "API homework fixture"
	Price        = "1500"
	Quantity     = "2"
	TitleLimit   = 50
)

//go:embed testdata/photo.jpg
var photo []byte

// Login авторизует один из двух независимых аккаунтов, заданных в локальном env.
func Login(t *testing.T, prefix string) authModels.LoginOkResponse {
	t.Helper()
	utils.LoadEnv()
	t.Logf("Подготовка: вход тестового пользователя %s", prefix)
	result := auth.Login(t, os.Getenv(prefix+"_LOGIN"), os.Getenv(prefix+"_PASSWORD"))
	require.NotEmpty(t, result.Token, "Авторизация не вернула токен")
	require.NotEmpty(t, result.User.Id, "Авторизация не вернула id пользователя")
	return result
}

// Users возвращает владельца объявления и отдельного пользователя избранного.
func Users(t *testing.T) (authModels.LoginOkResponse, authModels.LoginOkResponse) {
	t.Helper()
	owner, viewer := Login(t, "TEST"), Login(t, "FAVORITES")
	require.NotEqual(t, owner.User.Id, viewer.User.Id, "TEST и FAVORITES должны быть разными пользователями")
	return owner, viewer
}

// Request создает независимые поля и фотографию для каждого теста.
func Request() models.CreateRequest {
	return models.CreateRequest{Fields: map[string]string{"title": "QA API " + utils.RandomString(20), "description": Description, "price": Price, "quantity": Quantity}, Photos: []models.PhotoUpload{{Name: "photo.jpg", Data: append([]byte(nil), photo...)}}}
}

// Advertisement создает собственные тестовые данные с автоматическим удалением.
func Advertisement(t *testing.T, token string) adModels.Advertisement {
	t.Helper()
	return response.Decode[adModels.Advertisement](t, advertisement.Create(t, token, Request()), http.StatusCreated)
}

// DeletedAdvertisementID получает гарантированно отсутствующий UUID через удаление своей фикстуры.
func DeletedAdvertisementID(t *testing.T, token string) string {
	t.Helper()
	ad := Advertisement(t, token)
	advertisement.Delete(t, token, ad.ID)
	return ad.ID
}
