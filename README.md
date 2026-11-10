# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Запуск тестов

1. Установите любое IDE для разработки: [Jetbrains GoLand](https://www.jetbrains.com/ru-ru/go/) (Платное или по студенческой лицензии), [VS Code](https://code.visualstudio.com/) (бесплатное).
2. Установите [Git](https://git-scm.com/install/) на компьютер
3. [Установите](https://go.dev/doc/install) Go версии 1.25.1
4. [Настройте](https://docs.github.com/ru/authentication/connecting-to-github-with-ssh) доступ к своему Github-аккаунту
   по SSH-ключу.
5. Склонируйте себе репозиторий: `git clone <путь к вашему репозиторию>`
6. Подтяните необходимые зависимости через команду в корне проекта:

```bash
   go mod tidy
```

7. В файле `.env.override` (создайте его в корне, если нет) выставьте правильные значения для переменных окружения:

```
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
```

8. Для запуска всех тестов достаточно выполнить команду в корне проекта:

```bash
  go test -v ./...
```

9. Отдельный пакет с тестами можно запустить командой:

```bash
  go test -v ./tests/scenarios/myAdvertisement/
```

## Практика 08: поиск объявлений по названию

Поиск использует `GET /advertisements?search=...` из Swagger и проходит через слои:
`internal/client/http/advertisements` → `internal/managers/advertisements` →
`tests/scenarios/advertisementSearch`.
Клиент проверяет HTTP 200, менеджер возвращает типизированный `AdvertisementsList`,
а общий матчер `assertions.ContainsAdvertisementByTitle` проверяет точное название в `Items`.

Тесты используют существующее объявление **«Тест названия»**. Они не создают и не удаляют данные.
В `search_test.go` заданы константы полного названия и подстрок из начала, середины и конца.
Если объявление переименовано или удалено, обновите эти константы под другое существующее
объявление. Оно должно попадать на первую страницу результатов поиска (по умолчанию 20 записей).
Для авторизации используется локальный `.env.override` из инструкции выше.

```bash
go test -v -count=1 ./tests/scenarios/advertisementSearch/
```

Помимо базового сценария, параметризованный тест проверяет шесть случаев: начало,
середину, конец, нижний, верхний и смешанный регистр. `utils.MixCase` работает с Unicode,
включая кириллицу; его отдельные проверки запускаются вместе с `go test -v ./...`.

## Логирование

Для включения логирования request/response в консоль: выставьте переменную окружения `DEBUG=true` в
`.env.override`

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```
