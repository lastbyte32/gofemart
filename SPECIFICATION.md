### Сводное HTTP API

Накопительная система лояльности «Гофермарт» должна предоставлять следующие HTTP-хендлеры:

| Метод        | Путь                           | Описание                                                                       |
|--------------|--------------------------------|--------------------------------------------------------------------------------|
| POST         | /api/user/register             | Регистрация пользователя                                                      |
| POST         | /api/user/login                | Аутентификация пользователя                                                   |
| POST         | /api/user/orders               | Загрузка пользователем номера заказа для расчёта                               |
| GET          | /api/user/orders               | Получение списка загруженных пользователем номеров заказов, их статусов и информации о начислениях |
| GET          | /api/user/balance              | Получение текущего баланса счёта баллов лояльности пользователя              |
| POST         | /api/user/balance/withdraw     | Запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа |
| GET          | /api/user/withdrawals          | Получение информации о выводе средств с накопительного счёта пользователем   |



#### **Регистрация пользователя**

Хендлер: `POST /api/user/register`.

Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.
После успешной регистрации должна происходить автоматическая аутентификация пользователя.

Формат запроса:

```
POST /api/user/register HTTP/1.1
Content-Type: application/json
...

{
	"login": "<login>",
	"password": "<password>"
}
```

Возможные коды ответа:

- `200` — пользователь успешно зарегистрирован и аутентифицирован;
- `400` — неверный формат запроса;
- `409` — логин уже занят;
- `500` — внутренняя ошибка сервера.

#### **Аутентификация пользователя**

Хендлер: `POST /api/user/login`.

Аутентификация производится по паре логин/пароль.

Формат запроса:

```
POST /api/user/login HTTP/1.1
Content-Type: application/json
...

{
	"login": "<login>",
	"password": "<password>"
}
```

Возможные коды ответа:

- `200` — пользователь успешно аутентифицирован;
- `400` — неверный формат запроса;
- `401` — неверная пара логин/пароль;
- `500` — внутренняя ошибка сервера.

#### **Загрузка номера заказа**

Хендлер: `POST /api/user/orders`.

Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.

Номер заказа может быть проверен на корректность ввода с помощью [алгоритма Луна](https://ru.wikipedia.org/wiki/Алгоритм_Луна){target="_blank"}.

Формат запроса:

```
POST /api/user/orders HTTP/1.1
Content-Type: text/plain
...

12345678903
```

Возможные коды ответа:

- `200` — номер заказа уже был загружен этим пользователем;
- `202` — новый номер заказа принят в обработку;
- `400` — неверный формат запроса;
- `401` — пользователь не аутентифицирован;
- `409` — номер заказа уже был загружен другим пользователем;
- `422` — неверный формат номера заказа;
- `500` — внутренняя ошибка сервера.

#### **Получение списка загруженных номеров заказов**

Хендлер: `GET /api/user/orders`.

Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых старых к самым новым. Формат даты — RFC3339.

Доступные статусы обработки расчётов:

- `NEW` — заказ загружен в систему, но не попал в обработку;
- `PROCESSING` — вознаграждение за заказ рассчитывается;
- `INVALID` — система расчёта вознаграждений отказала в расчёте;
- `PROCESSED` — данные по заказу проверены и информация о расчёте успешно получена.

Формат запроса:

```
GET /api/user/orders HTTP/1.1
Content-Length: 0
```

Возможные коды ответа:

- `200` — успешная обработка запроса.

  Формат ответа:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    [
    	{
            "number": "9278923470",
            "status": "PROCESSED",
            "accrual": 500,
            "uploaded_at": "2020-12-10T15:15:45+03:00"
        },
        {
            "number": "12345678903",
            "status": "PROCESSING",
            "uploaded_at": "2020-12-10T15:12:01+03:00"
        },
        {
            "number": "346436439",
            "status": "INVALID",
            "uploaded_at": "2020-12-09T16:09:53+03:00"
        }
    ]
    ```

- `204` — нет данных для ответа.
- `401` — пользователь не авторизован.
- `500` — внутренняя ошибка сервера.

#### **Получение текущего баланса пользователя**

Хендлер: `GET /api/user/balance`.

Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов лояльности, а также сумме использованных за весь период регистрации баллов.

Формат запроса:

```
GET /api/user/balance HTTP/1.1
Content-Length: 0
```

Возможные коды ответа:

- `200` — успешная обработка запроса.

  Формат ответа:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    {
    	"current": 500.5,
    	"withdrawn": 42
    }
    ```

- `401` — пользователь не авторизован.
- `500` — внутренняя ошибка сервера.

#### **Запрос на списание средств**

Хендлер: `POST /api/user/balance/withdraw`

Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя в счет оплаты которого списываются баллы.

Примечание: для успешного списания достаточно успешной регистрации запроса, никаких внешних систем начисления не предусмотрено и не требуется реализовывать.

Формат запроса:

```
POST /api/user/balance/withdraw HTTP/1.1
Content-Type: application/json

{
	"order": "2377225624",
    "sum": 751
}
```

Здесь `order` — номер заказа, а `sum` — сумма баллов к списанию в счёт оплаты.

Возможные коды ответа:

- `200` — успешная обработка запроса;
- `401` — пользователь не авторизован;
- `402` — на счету недостаточно средств;
- `422` — неверный номер заказа;
- `500` — внутренняя ошибка сервера.

#### **Получение информации о выводе средств**

Хендлер: `GET /api/user/balance/withdrawals`.

Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых старых к самым новым. Формат даты — RFC3339.

Формат запроса:

```
GET /api/user/withdrawals HTTP/1.1
Content-Length: 0
```

Возможные коды ответа:

- `200` — успешная обработка запроса.

  Формат ответа:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    [
        {
            "order": "2377225624",
            "sum": 500,
            "processed_at": "2020-12-09T16:09:57+03:00"
        }
    ]
    ```

- `204` - нет ни одного списания.
- `401` — пользователь не авторизован.
- `500` — внутренняя ошибка сервера.

### Взаимодействие с системой расчёта начислений баллов лояльности

Для взаимодействия с системой доступен один хендлер:

- `GET /api/orders/{number}` — получение информации о расчёте начислений баллов лояльности.

Формат запроса:

```
GET /api/orders/{number} HTTP/1.1
Content-Length: 0
```

Возможные коды ответа:

- `200` — успешная обработка запроса.

  Формат ответа:

    ```
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    {
        "order": "<number>",
        "status": "PROCESSED",
        "accrual": 500
    }
    ```

  Поля объекта ответа:

    - `order` — номер заказа;
    - `status` — статус расчёта начисления:

        - `REGISTERED` — заказ зарегистрирован, но не начисление не рассчитано;
        - `INVALID` — заказ не принят к расчёту, и вознаграждение не будет начислено;
        - `PROCESSING` — расчёт начисления в процессе;
        - `PROCESSED` — расчёт начисления окончен;

    - `accrual` — рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.

- `204` - заказ не зарегистрирован в системе расчета.

- `429` — превышено количество запросов к сервису.

  Формат ответа:

    ```
    429 Too Many Requests HTTP/1.1
    Content-Type: text/plain
    Retry-After: 60
    
    No more than N requests per minute allowed
    ```

- `500` — внутренняя ошибка сервера.
