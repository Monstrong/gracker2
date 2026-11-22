include config/.env
export

.PHONY: migrate-up migrate-down migrate-step migrate-step-down migrate-force run

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-step:
	go run ./cmd/migrate step 1

migrate-step-down:
	go run ./cmd/migrate step -1

migrate-force:
	go run ./cmd/migrate force $(v)

run:
	go run ./cmd/server