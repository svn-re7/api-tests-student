package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type Result struct {
	Status int
	Body   []byte
}

type Error struct {
	Code    string `json:"error"`
	Message string `json:"message"`
}

const (
	InvalidRequest = "invalid_request"
	Unauthorized   = "unauthorized"
	NotFound       = "not_found"
	Forbidden      = "forbidden"
	AlreadyExists  = "already_exists"
	LimitReached   = "limit_reached"
)

// Read читает и закрывает HTTP-ответ, включая тело в диагностику последующих проверок.
func Read(t *testing.T, raw *http.Response) Result {
	t.Helper()
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Чтение HTTP %d", raw.StatusCode)
	return Result{Status: raw.StatusCode, Body: body}
}

// Decode проверяет HTTP-статус и разбирает JSON в заданную модель.
func Decode[T any](t *testing.T, result Result, status int) T {
	t.Helper()
	require.Equal(t, status, result.Status, "HTTP-ответ: %s", result.Body)
	var value T
	require.NoError(t, json.Unmarshal(result.Body, &value), "Разбор JSON: %s", result.Body)
	return value
}

// RequireError проверяет статус, код ошибки и стабильную часть сообщения.
func RequireError(t *testing.T, result Result, status int, code, message string) {
	t.Helper()
	value := Decode[Error](t, result, status)
	require.Equal(t, code, value.Code, "Код ошибки: %s", result.Body)
	require.Contains(t, value.Message, message, "Сообщение ошибки: %s", result.Body)
	t.Logf("Проверка: HTTP %d, error=%s, message содержит %q", status, code, message)
}
