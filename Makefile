run:
	go run cmd/user-service/main.go

migrate-file:
	migrate create -ext sql -dir migrations/ -seq user_temp_table

DB_URL := "postgres://postgres:+_+diyor2005+_+@localhost:5432/user_service?sslmode=disable"

migrate-up:
	migrate -path migrations -database $(DB_URL) -verbose up

migrate-down:
	migrate -path migrations -database $(DB_URL) -verbose down

migrate-force:
	migrate -path migrations -database $(DB_URL) -verbose force 1