package advertisementSearch

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisement"
	advList "api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/response"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

// Подстроки относятся к префиксу собственной фикстуры каждого прогона.
const (
	existingAdTitle = "Тест названия"
	searchBeginning = "Тест"
	searchMiddle    = "назва"
	searchEnd       = "вания"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
	adTitle   string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Проверяем, что константы для поиска заданы")
	require.NotEmpty(s.T(), existingAdTitle, "Задайте existingAdTitle")
	require.NotEmpty(s.T(), searchBeginning, "Задайте searchBeginning")
	require.NotEmpty(s.T(), searchMiddle, "Задайте searchMiddle")
	require.NotEmpty(s.T(), searchEnd, "Задайте searchEnd")

	base.Precondition("Проверяем, что подстроки соответствуют названию объявления")
	require.True(s.T(), strings.HasPrefix(existingAdTitle, searchBeginning),
		"searchBeginning должен совпадать с началом existingAdTitle")
	require.Contains(s.T(), existingAdTitle, searchMiddle,
		"searchMiddle должен входить в existingAdTitle")
	require.True(s.T(), strings.HasSuffix(existingAdTitle, searchEnd),
		"searchEnd должен совпадать с концом existingAdTitle")

	base.Precondition("Авторизуемся для поиска объявлений")
	s.loginData = fixtures.Login(s.T(), "TEST")
	request := fixtures.Request()
	request.Fields["title"] = existingAdTitle + " " + utils.RandomString(16)
	ad := response.Decode[adModels.Advertisement](s.T(), advertisement.Create(s.T(), s.loginData.Token, request), http.StatusCreated)
	s.adTitle = ad.Title
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestSearchByTitle() {
	list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, searchMiddle)

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items)
	})
	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.adTitle)
	})
}

func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchMiddle)},
		{"верхний регистр", strings.ToUpper(searchMiddle)},
		{"смешанный регистр", utils.MixCase(searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.adTitle)
		})
	}
}
