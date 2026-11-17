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
4. Склонируйте свой fork через HTTPS: `git clone https://github.com/<ваш-логин>/api-tests-student.git`.
5. Перейдите в каталог: `cd api-tests-student`.
6. Подтяните необходимые зависимости через команду в корне проекта:

```bash
   go mod tidy
```

7. В файле `.env.override` (создайте его в корне, если нет) выставьте правильные значения для переменных окружения:

```
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
FAVORITES_LOGIN=another_user_login
FAVORITES_PASSWORD=another_user_password
DEBUG=false
```

Зарегистрируйте два отдельных тестовых аккаунта на TestBoard. `TEST` создает объявления,
`FAVORITES` работает с избранным другого пользователя. Аккаунты должны различаться.
Для теста лимита избранное аккаунта `FAVORITES` перед запуском должно быть пустым.
Тесты не удаляют постороннее избранное. `.env.override` исключен из Git; пароли и токены
не нужно коммитить. Адрес API уже задан в `.env`: `https://testboard.avito.com/api/v1`.

8. Для запуска всех тестов достаточно выполнить команду в корне проекта:

```bash
  go test -v -count=1 -p=1 ./...
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

Тесты создают собственное объявление с префиксом "Тест названия" и уникальным суффиксом.
После завершения оно удаляется через `t.Cleanup`. Тест профиля также создает свою фикстуру.
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

1. Установите [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/) версии 2.6.1.
   Например, при установленном Go:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.1
export PATH="$(go env GOPATH)/bin:$PATH"
```

2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```

## API-домашняя работа: создание объявления и избранное

Новые сценарии находятся в `tests/scenarios/homework`, по одному сценарию в файле.
Слои: `internal/client/http` отправляет запросы, `internal/managers` разбирает ответы
и проверяет протокол, `tests/scenarios` задает последовательность и бизнес-проверки.
Общие данные и авторизация находятся в `tests/fixtures`; JPG включен через `go:embed`.

| Файл | Сценарий |
| --- | --- |
| `create_positive_test.go` | Все параметры, получение по id, поиск полного названия, фотографии |
| `create_invalid_test.go` | Отсутствующие и пустые поля, превышение длины, некорректные числа, 0 и 4 фото |
| `create_unauthorized_test.go` | Создание без токена и с неверным токеном |
| `favorites_cycle_test.go` | Добавление, удаление из избранного, удаление самого объявления |
| `favorites_missing_test.go` | Добавление и удаление отсутствующего объявления |
| `favorites_unauthorized_test.go` | GET, POST, DELETE без токена и с неверным токеном |
| `favorites_own_test.go` | Запрет добавления своего объявления, 403 |
| `favorites_duplicate_test.go` | Повторное добавление, 409, отсутствие дубликатов |
| `favorites_limit_test.go` | 50 записей допустимы, 51-я отклоняется с 400 |
| `favorites_invalid_test.go` | Отсутствующий, пустой и некорректный advertisement_id, 400 |

Актуальный Swagger задает `GET /advertisement?id={id}` и `DELETE /advertisement?id={id}`.
Фото получаются через `GET /advertisements/{id}/photos`. Ошибки проверяются по HTTP-статусу,
полю `error` и стабильной части `message`. Для `POST /favorites` воспроизводятся 400,
401, 403, 404 и 409. Ответ 500 требует внутреннего сбоя сервера; стабильного способа вызвать
его пользовательскими данными Swagger не описывает, поэтому такой тест не добавлен.

Каждый тест готовит свои объявления с уникальными названиями. Очистка регистрируется сразу
после успешного создания и выполняется при падении проверки. Удаление объявления также
удаляет связанные записи избранного. Уже удаленная запись при очистке определяется через GET.
Сценарии не используют `t.Parallel`: аккаунт `FAVORITES` предназначен для одного прогона
за раз. Порядок тестов значения не имеет. Аккаунты сохраняются для следующих запусков.
При внезапном завершении процесса очистка Go не выполняется; удалите оставшиеся тестовые
объявления владельца перед следующим запуском проверки лимита.

Запуск новых сценариев и трех последовательных прогонов всего проекта без кеширования:

```bash
go test -v -count=1 ./tests/scenarios/homework/
go test -v -count=3 -p=1 ./...
golangci-lint run
```

Флаг `-p=1` запускает пакеты последовательно, снижая нагрузку на общий стенд. При одновременном запуске пакетов наблюдались ответы nginx 503.
Флаг `-v` показывает подготовку, шаги, проверки и очистку. С `DEBUG=true` дополнительно
печатаются HTTP-запросы и ответы, включая учетные данные и токены; такие логи не публикуйте.
