lint:
	@echo " > linting..."
	@staticcheck ./...
	@echo " > linting successfully finished"

test:
	@echo " > testing..."
	@go test -gcflags="-l" -race -v ./...
	@echo " > successfully finished"

vet:
	@echo " > go vet..."
	@go vet ./...
	@echo " > go vet successfully finished"

.PHONY: lint test vet