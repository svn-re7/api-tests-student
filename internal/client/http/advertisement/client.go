package advertisement

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisement/models"

	"github.com/stretchr/testify/require"
)

const advertisementPath = "/advertisement"

// Create отправляет multipart-форму, сохраняя различие отсутствующего и пустого поля.
func Create(t *testing.T, token string, data models.CreateRequest) *http.Response {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for name, value := range data.Fields {
		require.NoError(t, writer.WriteField(name, value), "Запись поля %s", name)
	}
	for _, photo := range data.Photos {
		part, err := writer.CreateFormFile("photos", photo.Name)
		require.NoError(t, err, "Создание multipart-файла")
		_, err = part.Write(photo.Data)
		require.NoError(t, err, "Запись фотографии")
	}
	require.NoError(t, writer.Close(), "Закрытие multipart-формы")
	return apiRunner.GetRunner().Auth(token).Create().Post(advertisementPath).
		ContentType(writer.FormDataContentType()).Body(body.String()).Expect(t).End().Response
}

// Get получает объявление через query-параметр id согласно Swagger.
func Get(t *testing.T, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Create().Get(advertisementPath).Query("id", id).Expect(t).End().Response
}

// Delete удаляет объявление указанного владельца.
func Delete(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(advertisementPath).Query("id", id).Expect(t).End().Response
}

// Photos получает список фотографий объявления.
func Photos(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get("/advertisements/" + id + "/photos").Expect(t).End().Response
}
