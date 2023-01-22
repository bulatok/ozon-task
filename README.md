# Укорачиватель ссылок

## Info
API позволяет укорачивать или возможно хранить ссылки, используя REST интерфейс.


В качестве базы данных можно выбрать
 1) postgres
 2) redis
 3) cache
## Запуск
```
    git clone https://github.com/bulatok/ozon-task.git
    cd ozon-task
    docker-compose up --build
```
При запуске можно указать флаг для выбора хранилища (Postgres или redis). Например:
```bash
    $ make build
    $ ./don -db cache     # запускаем in memory вариант
    $ ./don -db postgres  # запускаем postgres вариант
    $ ./don               # запускаем postgres вариант (если без флага)
    $ ./don -db reids     # запускаем redis вариант
```


Если запуск через докер, то можно поставить флаги в ```Dockerfile```.
## Пример

### Запрос на создание
Пример POST запроса через curl. Ответом является сокращенная ссылка, полученная путем хеширования.

```bash
    $ curl -X POST http://127.0.0.1:8080/new -H 'Content-Type: application/json' -d '{"original_link":"https://www.kinopoisk.ru/"}'
    $ {"short_link":"http://127.0.0.1:8080/6PQggdj6v8"} 
```

### Запрос на получение
```bash
    $ curl -X GET http://127.0.0.1:8080/6PQggdj6v8
    $ {"original_link":"https://www.kinopoisk.ru/"}
```

### GRPC
В качестве доставки также реализовал grpc сервер в `./internal/ozon-task/grpc/v1`, который инициализируется при сборке

Для него же написал простенький клиент в `./client/cmd/grpc`, 
чтобы запустить нужно прописать `make grpc-client && ./client/cmd/grpc/don`.

.proto файл находится в `./pkg/pb/proto/links.proto`

### Тесты
Весь функционал api покрыл тестами, которые лежат `./internal/ozon-task/delivery/http/api/v1/links_test.go`