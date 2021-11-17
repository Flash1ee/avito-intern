# Тестовое задание на позицию стажера-бекендера

## Микросервис для работы с балансом пользователей.

**Проблема:**

В нашей компании есть много различных микросервисов. Многие из них так или иначе хотят взаимодействовать с балансом
пользователя. На архитектурном комитете приняли решение централизовать работу с балансом пользователя в отдельный
сервис.

**Задача:**

Необходимо реализовать микросервис для работы с балансом пользователей (зачисление средств, списание средств, перевод
средств от пользователя к пользователю, а также метод получения баланса пользователя). Сервис должен предоставлять HTTP
API и принимать/отдавать запросы/ответы в формате JSON.

## Описание

- Сервис реализован на языке `Golang` с использованием чистой архитектуры.
- Уровни `delivery`, `usecase`, `repository` с использованием библиотеки `testify` и `gomock`.
- Уровень `repository` протестирован с использованием sqlmock.
- Возвращаемые ошибки `repository`, `usecase` задокументированы с использованием встроенных в `Golang` возможностей.
- Логирование операций в файл в папку `/logs`
- Реализованы `Middlewares` для: отслеживания паники, логирования.
- Валидация реализована с помощью `ozzo-validation`.
- Конфигурирование приложения с использованием .toml файла.
- Сервис поднимается в `Docker` контейнерах: база данных и основное приложение.
- Контейнеры конфигурируются в  `docker-compose`.
- В качестве СУБД используется `PostgreSQL`.
- API задокументировано с использованием Swagger по адресу `http:://localhost:8080/api/v1/swagger/`
- Взаимодействие с проектом организовано посредством `Makefile`.
- Подключен `Github Actions` для проверки стиля и сборки приложения.

## Дополнительные задания
- Реализовано доп задание не получение списка транзакций.
- Реализовано доп. задание на конвертацию баланса пользователя в другие валюты.
- Использовано API ЦБ РФ.
  Доступные валюты:
```text
"AMD", "NOK", "TRY", "USD", "CAD", "CNY", "UAH", "CZK",
"JPY", "GBP", "HUF", "MDL", "UZS", "AUD", "INR", "EUR",
"KGS", "TMT", "ZAR", "BRL", "HKD", "KZT", "SEK", "CHF",
"KRW", "BYN", "BGN", "PLN", "SGD", "AZN", "DKK", "RON",
"XDR", "TJS"
```

- [x] Использование docker и docker-compose для поднятия и развертывания dev-среды.
- [x] Методы АПИ возвращают человеко-читабельные описания ошибок и соответвующие статус коды при их возникновении.
- [x] Все реализовано на GO, все-же мы собеседуем на GO разработчика.
- [x] Написаны unit/интеграционные тесты.

## Запуск

Из корня репозитория запустить следующую команду.  
Должен быть свободен порт 8080.  
PostgreSQL поднимается на 5432 порту, он тоже должен быть свободен.

```bash
docker-compose up
```

## API

По умолчанию валюта кошелька - рубль.

1. Метод получения баланса. Принимает id пользователя и сколько средств зачислить.

```text
GET /api/v1/balance/{:user_id}?currency={currency}
```

### Параметры запроса

- user_id - id пользователя
- currency(опциональный параметр) - при наличии параметра баланс переводится в указанную валюту.

### Ответ
**200 - OK**
```json
{
  "balance": 0,
  "user_id": 0
}
```
PS: если присутсвует `currency`, то баланс конвертируется в указанную валюту.  
**400, 404, 422, 500**
```json
{
  "error": "string"
}
```
### Описание кодов ответа
- 200 - OK;
- 400 - неверный параметр в запросе;
- 404 - пользователь с id = user_id не найден;
- 422 - не поддерживаемая API валюта для конвертации;
- 422 - ошибка обращения к API получения курса;
- 500 - ошибка сервера.

2. Зачислить или снять с баланса кошелька средства.

```text
POST /api/v1/balance/{:user_id}
```

### Параметры запроса

- user_id - id пользователя 

### Тело запроса:

```json
{
  "amount": 0,
  "operation": 0
}
```

- amount - сумма снятия/пополнения;
- operaion - 0 - снятие, 1 - пополнение; 
### Ответ
**200 - OK**
```json
{
  "balance": 0,
  "user_id": 0
}
```
**400, 404, 422, 500**
```json
{
  "error": "string"
}
```
### Описание кодов ответа
- 200 - OK;
- 400 - неверный параметр в запросе;
- 404 - пользователь с id = user_id не найден;
- 422 - ошибка в теле запроса;
- 422 - недостаточно средств для выполнения операции(для снятия);
- 500 - ошибка сервера.

3. Перевод между кошельками.

```text
POST /api/v1/transfer
```

### Тело запроса

```json
{
  "amount": 0,
  "receiver_id": 0,
  "sender_id": 0
}
```

- amount - сумма перевода;
- receiver_id - id получателя;
- sender_id - id отправителя.  
### Ответ
**200 - OK**
```json
{
  "receiver_balance": 0,
  "receiver_id": 0,
  "sender_balance": 0,
  "sender_id": 0
}
```
**404, 422, 500**
```json
{
  "error": "string"
}
```
### Описание кодов ответа
- 200 - OK;
- 404 - отправитель/получатель не найден;
- 422 - ошибка в теле запроса;
- 422 - недостаточно средств для перевода;
- 500 - ошибка сервера.

4. Получение списка транзакций.
```text
GET /transaction/{:user_id}
```
### Параметры запроса
- user_id - id пользователя;
- page - номер страницы списка транзакций;
- count - количество транзакций на 1 странице;
- sort - сортировка по полю, опциональный параметр
  - sort=date - сортировка по дате транзакции;
  - sort=sum - сортировка по размеру транзакции;
- direction - условие сортировки(только если указан параметр sort)
  - asc - по возрастанию;
  - desc - по убыванию

### Ответ
**200 - OK**
```json
{
  "transactions": [
    {
      "amount": 0,
      "created_at": "2021-17-11T20:09:22.058143Z",
      "description": "string",
      "id": 0,
      "receiver_id": 0,
      "type": "string"
    }
  ]
}
```
**204 - NoContent**
```json
{
  "OK": "transactions not found"
}
```
**404, 422, 500**
```json
{
  "error": "string"
}
```
### Описание кодов ответа
- 200 - OK;
- 204 - транзакций пользователя не найдено;
- 400 - ошибка в параметрах запроса
  - в id пользователя;
  - в page;
  - в count;
  - в sort;
  - в direction;
  - использование direction без sort;
- 500 - ошибка сервера.

### Примеры запросов
**Запросы делались с использованием `Postman`**

Запрос на количество транзакций без сортировки:
```text
GET http://localhost:8080/api/v1/transaction/2?page=2&count=3
```
### Ответ
**Код ответа - 200**

```json
{
    "transactions": [
        {
            "sender_id": 2,
            "receiver_id": 1,
            "type": "transfer",
            "description": "user <id = 2> send money to user <id = 1>",
            "created_at": "2021-11-17T19:56:24.099571Z",
            "amount": 100
        },
        {
            "sender_id": 2,
            "receiver_id": 1,
            "type": "transfer",
            "description": "user <id = 2> send money to user <id = 1>",
            "created_at": "2021-11-17T19:56:29.102226Z",
            "amount": 100
        }
    ]
}
```
Запрос на количество транзакций c сортировкой по дате:
```text
GET http://localhost:8080/api/v1/transaction/2?page=2&count=3&sort=date
```
### Ответ
**Код ответа - 200**

```json
{
    "transactions": [
        {
            "sender_id": 2,
            "receiver_id": 1,
            "type": "transfer",
            "description": "user <id = 2> send money to user <id = 1>",
            "created_at": "2021-11-17T19:56:24.099571Z",
            "amount": 100
        },
        {
            "sender_id": 2,
            "receiver_id": 1,
            "type": "transfer",
            "description": "user <id = 2> send money to user <id = 1>",
            "created_at": "2021-11-17T19:56:29.102226Z",
            "amount": 100
        }
    ]
}
```
Запрос на количество транзакций c сортировкой по дате по убыванию:
```text
GET http://localhost:8080/api/v1/transaction/2?page=2&count=3&sort=date&direction=desc
```
### Ответ
**Код ответа - 200**

```json
{
  "transactions": [
    {
      "sender_id": 2,
      "type": "refill",
      "description": "user <id = 2> refill account",
      "created_at": "2021-11-17T19:55:47.277489Z",
      "amount": 100
    },
    {
      "sender_id": 2,
      "type": "refill",
      "description": "user <id = 2> refill account",
      "created_at": "2021-11-17T19:55:45.759548Z",
      "amount": 100
    }
  ]
}
```
Запрос на количество транзакций c сортировкой по дате по убыванию с первой страницы 
и количеством заведомо большим имеющихся в базе транзакций:
```text
GET http://localhost:8080/api/v1/transaction/2?page=1&count=10&sort=date&direction=desc
```
### Ответ
**Код ответа - 200**

```json
{
    "transactions": [
        {
            "sender_id": 2,
            "receiver_id": 1,
            "type": "transfer",
            "description": "user <id = 2> send money to user <id = 1>",
            "created_at": "2021-11-17T19:56:29.102226Z",
            "amount": 100
        },
        {
            "sender_id": 2,
            "receiver_id": 1,
            "type": "transfer",
            "description": "user <id = 2> send money to user <id = 1>",
            "created_at": "2021-11-17T19:56:24.099571Z",
            "amount": 100
        },
        {
            "sender_id": 2,
            "type": "refill",
            "description": "user <id = 2> refill account",
            "created_at": "2021-11-17T19:55:48.480167Z",
            "amount": 100
        },
        {
            "sender_id": 2,
            "type": "refill",
            "description": "user <id = 2> refill account",
            "created_at": "2021-11-17T19:55:47.277489Z",
            "amount": 100
        },
        {
            "sender_id": 2,
            "type": "refill",
            "description": "user <id = 2> refill account",
            "created_at": "2021-11-17T19:55:45.759548Z",
            "amount": 100
        }
    ]
}
```
Запрос для несуществующего пользователя:
```text
GET http://localhost:8080/api/v1/transaction/5?page=1&count=10&sort=date&direction=desc
```
### Ответ
**Код ответа - 204**
```json

```
Запрос с ошибкой в параметре запроса:
```text
GET http://localhost:8080/api/v1/transaction/5?page=1&count=-1&sort=date&direction=desc```
### Ответ
**Код ответа - 400**
```json
{"error":"invalid query param count"}
```
Запрос с ошибкой в параметре запроса - использование направления сортировки в sort:
```text
GET http://localhost:8080/api/v1/transaction/5?page=1&count=5&direction=desc
```
### Ответ
**Код ответа - 400**
```json
{"error":"invalid query param direction"}
```
Запрос без пагинации:
```text
GET http://localhost:8080/api/v1/transaction/5
```
### Ответ
**Код ответа - 400**
```json
{"error":"invalid parameters in query"}
```