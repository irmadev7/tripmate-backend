run:
	go run main.go

migration:
	@goose -dir=docs/migrations create $(filename) sql

migrate:
	goose -dir=docs/migrations postgres "host=localhost port=5432 user=postgres dbname=postgres password= sslmode=disable" up

