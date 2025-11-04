run:
	go run main.go

migration:
	@goose -dir=docs/migrations create $(filename) sql

migrate:
	goose -dir=docs/migrations postgres "user=postgres password= postgres dbname=postgres sslmode=disable" up

