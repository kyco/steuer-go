# Justfile for steuer-go

# Default recipe to run when just is called without arguments
default:
    @just --list

# Run the application
run:
    go run cmd/tax-calculator/main.go

# Build the application
build:
    go build -o tax-calculator ./cmd/tax-calculator/

# Run tests
test:
    go test ./...

# Run tests with coverage
test-coverage:
    go test -cover ./...

# Clean build artifacts
clean:
    rm -f tax-calculator

# Install dependencies
deps:
    go mod download

# Update dependencies
update-deps:
    go get -u ./...
    go mod tidy

# Format code
fmt:
    go fmt ./...

# Lint code
lint:
    go vet ./...

# Create new branch
branch name:
    git checkout -b {{name}}
