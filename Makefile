.PHONY: db-up migrate-up migrate-down run-back run-front

# Run PostgreSQL via Docker Compose
db-up:
	docker compose up -d

# Run database migrations
migrate-up:
	cd backend && migrate -path migrations -database "postgres://lifecap:lifecap_secret@localhost:5433/life_capital?sslmode=disable" -verbose up

migrate-down:
	cd backend && migrate -path migrations -database "postgres://lifecap:lifecap_secret@localhost:5433/life_capital?sslmode=disable" -verbose down

# Run backend server
run-back:
	cd backend && go run cmd/server/main.go

down-back:
	fuser -k 8080/tcp || true

# Run frontend dev server
run-front:
	cd frontend && npm run dev

down-front:
	fuser -k 5173/tcp || true
