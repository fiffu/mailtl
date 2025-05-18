all: deps tooling test

MIGRATIONS_DIR := app/storage/migrations
migration:
	@NAME=
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migration NAME='Your migration name'"; \
		exit 1; \
	fi
	@mkdir -p $(MIGRATIONS_DIR)
	@timestamp=$$(date +%s); \
	safe_name=$$(echo $(NAME) | tr ' ' '_'); \
	filename=$$timestamp"_"$$safe_name".sql"; \
	touch "$(MIGRATIONS_DIR)/$$filename"; \
	echo "Created migration: $(MIGRATIONS_DIR)/$$filename"

deps:
	go get

tooling:
	go install github.com/vektra/mockery/v2@v2.38.0
	go install gotest.tools/gotestsum@v1.8.1
	cp pre-commit .git/hooks/pre-commit

test:
	gotestsum -f dots -- -failfast -covermode=count -coverprofile coverage.out ./...

	@go tool cover -func=coverage.out | grep 'total' | sed -e 's/\t\+/ /g' | sed -e 's/total: (statements)/TEST COVERAGE:/'

CONFIG ?= sample.config.json
run:
	go run main.go $(CONFIG)
