include .env

.PHONY: dev
dev:
	go run ./cmd/api -env dev

.PHONY: goose-up
goose-up:
	goose -dir=migration postgres ${DB_DSN} up

.PHONY: goose-down
goose-down:
	goose -dir=migration postgres ${DB_DSN} down

.PHONY: goose-up-one
goose-up-one:
	goose -dir=migration postgres ${DB_DSN} up-by-one

.PHONY: goose-down-one
goose-down-one:
	goose -dir=migration postgres ${DB_DSN} down-by-one

.PHONY: goose-status
goose-status:
	goose -dir=migration postgres ${DB_DSN} status

.PHONY: goose-validate
goose-reset:
	goose -dir=migration postgres ${DB_DSN} reset

.PHONY: jet-gen
jet-gen:
	jet -dsn=${DB_DSN} -path=./api/model/

.PHONY: build
build:
	go build -o ./bin/app github.com/DaoVuDat/trackpro-api/cmd/api
