package homework

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

func TestCreateAdvertisementAllParameters(t *testing.T) {
	owner := fixtures.Login(t, "TEST")
	data := fixtures.Request()
	created := response.Decode[adModels.Advertisement](t, advertisement.Create(t, owner.Token, data), http.StatusCreated)
	saved := advertisement.Get(t, created.ID)
	t.Log("Проверка: GET возвращает все переданные параметры и владельца")
	require.Equal(t, data.Fields["title"], saved.Title, "title не сохранен")
	require.Equal(t, data.Fields["description"], saved.Description, "description не сохранен")
	price, err := strconv.ParseFloat(data.Fields["price"], 64)
	require.NoError(t, err)
	quantity, err := strconv.Atoi(data.Fields["quantity"])
	require.NoError(t, err)
	require.Equal(t, price, saved.Price, "price не сохранен")
	require.Equal(t, quantity, saved.Quantity, "quantity не сохранен")
	require.Equal(t, owner.User.Id, saved.UserID, "Неверный владелец")
	require.Equal(t, created.ID, saved.ID, "Неверный id")
	require.NotEmpty(t, saved.CreatedAt, "Не задана дата создания")
	require.NotEmpty(t, saved.UpdatedAt, "Не задана дата обновления")
	t.Log("Шаг: поиск по полному уникальному названию")
	list := advertisements.GetAdvertisementsList(t, owner.Token, data.Fields["title"])
	require.Len(t, list.Items, 1, "Полное уникальное название должно найти только созданное объявление")
	require.Equal(t, created.ID, list.Items[0].ID, "Поиск вернул другое объявление")
	require.Equal(t, saved.Title, list.Items[0].Title, "Поиск вернул другое название")
	photos := advertisement.Photos(t, owner.Token, created.ID)
	require.Len(t, photos, len(data.Photos), "Количество фотографий в API не совпадает с загрузкой")
	require.Len(t, saved.Photos, len(data.Photos), "GET объявления потерял фотографии")
	require.Len(t, created.Photos, len(data.Photos), "Создание не вернуло фотографии")
	require.NotEmpty(t, photos[0].ID, "У фото нет id")
	require.NotEmpty(t, photos[0].URL, "У фото нет URL")
	require.Equal(t, created.Photos[0].ID, photos[0].ID, "Фото не связано с результатом создания")
	require.Equal(t, saved.Photos[0].URL, photos[0].URL, "GET возвращает другой URL фото")
	require.Equal(t, 0, photos[0].SortOrder, "Неверный порядок фото")
	t.Log("Проверка: фото загружено и доступно через GET /advertisements/{id}/photos")
}
