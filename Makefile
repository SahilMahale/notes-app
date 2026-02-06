.PHONY: all backend frontend dev stop clean

# Run both backend and frontend
dev:
	@echo "Starting backend and frontend..."
	@make -j2 backend frontend

# Run backend only
backend:
	@echo "Starting backend on localhost:8001..."
	cd notes-backend && APP_AUTH=./secrets go run ./cmd/main.go

# Run frontend only
frontend:
	@echo "Starting frontend..."
	cd notes-ui && bun run dev

# Stop all running services
stop:
	@echo "Stopping services..."
	@-pkill -f "go run ./cmd/main.go" 2>/dev/null || true
	@-pkill -f "bun.*src/index.ts" 2>/dev/null || true
	@echo "Services stopped"

# Install dependencies
install:
	@echo "Installing backend dependencies..."
	cd notes-backend && go mod download
	@echo "Installing frontend dependencies..."
	cd notes-ui && bun install

# Build backend
build-backend:
	@echo "Building backend..."
	cd notes-backend && go build -o bin/notes-server ./cmd/main.go

# Run built backend binary
run-backend-bin:
	@echo "Running backend binary..."
	./notes-backend/bin/notes-server
