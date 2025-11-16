up:
	docker compose -f docker-compose.yaml up -d --build
	
down:
	docker compose -f docker-compose.yaml down

test-e2e:
	@echo "Running E2E tests..."
	@go test -v ./e2e/...

lint:
	@golangci-lint --version; \
	CGO_ENABLED=0 golangci-lint run -v