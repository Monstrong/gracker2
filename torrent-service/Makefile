include config/.env
export

.PHONY: migrate-up migrate-down run

migrate-up:
	migrate -path migrations -database "$(MIGRATE_DSN)" up

migrate-down:
	migrate -path migrations -database "$(MIGRATE_DSN)" down 1

run:
	go run ./cmd/server