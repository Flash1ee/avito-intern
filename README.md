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
- Логирование операций в файл в папку `/logs`
- Реализованы `Middlewares` для: отслеживания паники, логирования.
- Валидация реализована с помощью `ozzo-validation`.
- Конфигурирование приложения с использованием .toml файла.
- Сервис поднимается в `Docker` контейнерах: база данных и основное приложение.
- Контейнеры конфигурируются в  `docker-compose`.
- В качестве СУБД используется `PostgreSQL`.
- API задокументировано с использованием Swagger по адресу `http:://localhost:8080/api/v1/swagger/`
- Взаимодействие с проектом организовано посредством `Makefile`.

## Дополнительные задания

- Реализовано доп. задание на конвертацию баланса пользователя в другие валюты.
- Использовано апи ЦБ РФ  
 **Доступные валюты**
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
### Коды ответа
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

- user_id - id пользователя Тело запроса:

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
### Коды ответа
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
### Коды ответа
- 200 - OK;
- 404 - отправитель/получатель не найден;
- 422 - ошибка в теле запроса;
- 422 - недостаточно средств для перевода;
- 500 - ошибка сервера.
