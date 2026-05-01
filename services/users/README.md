## Сервис users

Для запуска выполнить команду из директории screw-it-desk-backend/services/users:

```bash
go run ./cmd/server/main.go -config-path="prod.env"
```

или

```bash
go run ./cmd/server/main.go -config-path="prod.env"
```


Создание миграции:
```bash
bin/goose create migrations/__migrations_name__ sql
```