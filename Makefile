include .env
MIGRATE_LOCAL=migrate -path=database/migration -database "mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" -verbose

MIGRATE=docker-compose exec web ${MIGRATE_LOCAL}

.PHONY: migrate-up migrate-down force goto drop create auto-create swag test-repo

dev:
		bash automate/scripts/gin-watch.sh ${SERVER_PORT}

migrate-up:
		@if [ "$(env)" = "local" ]; then $(MIGRATE_LOCAL) up; else $(MIGRATE) up; fi

migrate-down:
		@if [ "$(env)" = "local" ]; then $(MIGRATE_LOCAL) down; else $(MIGRATE) down; fi

force:
	@read -p "Which version do you want to force?" VERSION; \
	if [ "$(env)" = "local" ]; then \
	 $(MIGRATE_LOCAL) force $$VERSION; \
	else \
		$(MIGRATE) force $$VERSION; \
	fi
goto:
		@read -p  "Which version do you want to migrate?" VERSION; \
		if [ "$(env)" = "local" ]; then \
		  	$(MIGRATE_LOCAL) goto $$VERSION; \
		else \
		    $(MIGRATE) goto $$VERSION; \
		fi

drop:
		@if [ "$(env)" = "local" ]; then $(MIGRATE_LOCAL) drop; else $(MIGRATE) drop; fi

create:
		@read -p  "What is the name of migration?" NAME; \
		if [ "$(env)" = "local" ]; then \
			$(MIGRATE_LOCAL) create -ext sql -dir migration $$NAME; \
		else \
			$(MIGRATE) create -ext sql -dir migration $$NAME; \
		fi

swag:
		@command -v swag >/dev/null 2>&1 || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
		swag fmt
		swag init --parseDependency --parseInternal

crud:
	bash automate/scripts/crud.sh

install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2
	git config core.hooksPath hooks

start: install
	docker-compose up

run:
	docker-compose up

test-repo: TEST_NAME=$(filter-out $@,$(MAKECMDGOALS))
test-repo:
	go test ./tests/repository_test -v -run $(TEST_NAME)

i-test-controller: TEST_NAME=$(filter-out $@,$(MAKECMDGOALS))
i-test-controller:
	go test ./tests/controllers_i_test -v -run $(TEST_NAME)

test-controller: TEST_NAME=$(filter-out $@,$(MAKECMDGOALS))
test-controller:
	go test ./tests/controllers_test -v -run $(TEST_NAME)

