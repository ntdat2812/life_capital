ifneq (,$(wildcard backend/.env))
    include backend/.env
    export
endif

.PHONY: db-up migrate-up migrate-down run-back run-front

# Run PostgreSQL via Docker Compose
db-up:
	docker compose up -d

# Run backend server
run-back:
	cd backend && go run cmd/server/main.go

down-back:
	fuser -k 8080/tcp || true

down-back-mac:
	sudo lsof -i tcp:8080 -t | xargs kill -9 || true

# Run frontend dev server
run-front:
	cd frontend && npm run dev

down-front:
	fuser -k 5173/tcp || true

# --- MIGRATION COMMANDS (Uses DATABASE_URL from backend/.env) ---
migrate-up:
	@echo "Running up migrations to database..."
	cd backend && migrate -path migrations -database "$${DATABASE_URL}" -verbose up

migrate-down:
	@echo "Running down migrations to database..."
	cd backend && migrate -path migrations -database "$${DATABASE_URL}" -verbose down
