# Укорачиватель ссылок

## Info
Есть всего один хендлер ```/```, который и выполняет всю основную логику.
При **POST** запросе в формате json ```{"url":"someURL"}``` в бд добавляется оригинальная ссылка и укороченная. Укорченная ссылка получается путем её **хеширования** (исполюзуется FNV64). При **GET** запросе программа обращается к бд и при нахождении возвращает исходную ссылку в **JSON** формате.
Также и ошибки(несуществующие ссылки и тд) приходят в **JSON** формате.

При запуске можно указать флаг для выбора хранилища (Postgres или in-memory). Например:
```bash
    $ make build
    $ ./bin/ozon-task -store_type in-memory # запускаем in-memory вариант
```
или запустить без флага и будет поствален Postgres, который по умолчанию.
```bash
    $ make build
    $ ./bin/ozon-task -store_type Postgres # запускаем Postgres вариант
```
Пример POST запроса через curl
```bash
    $ curl -X POST http://localhost:8080 -H 'Content-Type: application/json' -d '{"url":"https://www.google.com/"}'
```

## Запуск как докер контейнер
```bash
    $ git clone https://github.com/bulatok/ozon-task.git
    $ make start
```