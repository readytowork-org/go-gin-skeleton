include .env

# DSN for Go's database driver (used by gentool)
DB_DSN="${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"
DB_DSN_DOCKER="${DB_USERNAME}:${DB_PASSWORD}@tcp(localhost:33066)/${DB_NAME}"

# Standard URL for Atlas (for use inside docker network)
DATABASE_URL="${DB_TYPE}://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
# Standard URL for Atlas (for local connection to docker db)
DATABASE_URL_LOCAL="${DB_TYPE}://${DB_USERNAME}:${DB_PASSWORD}@localhost:33066/${DB_NAME}"

GEN_TOOL=gentool -fieldNullable -fieldWithIndexTag -fieldWithTypeTag -fieldSignable -onlyModel -outPath './database/dao' -modelPkgName 'dao'

# Capture arguments passed after `migrate` target. Defaults to `apply`.
migrate_args = $(filter-out migrate,$(MAKECMDGOALS))

migrate:
		 @echo "using database: ${DB_NAME}"
		 @if [ "$(env)" = "local" ]; then \
			atlas migrate $(or $(migrate_args),apply) --env mysql \
				--var "DATABASE_URL=${DATABASE_URL}"; \
		 else \
			docker-compose exec web atlas migrate $(or $(migrate_args),apply) --env mysql \
				--var "DATABASE_URL=${DATABASE_URL_LOCAL}"; \
		 fi

dao:
		@command -v gentool >/dev/null 2>&1 || (echo "Installing gentool..." && go install gorm.io/gen/tools/gentool@latest)
		@if [ "$(env)" = "local" ]; then $(GEN_TOOL) -dsn $(DB_DSN); else $(GEN_TOOL) -dsn $(DB_DSN_DOCKER); fi

swagger:
		@command -v swag >/dev/null 2>&1 || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
		swag fmt
		swag init --output ./swagger --parseDependency --parseInternal --requiredByDefault

crud:
		bash automate/scripts/crud.sh

lint-install:
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2
		git config core.hooksPath hooks

run:
		@command -v gin >/dev/null 2>&1 || (go install github.com/codegangsta/gin@latest);
		gin -a $(SERVER_PORT) -i -p $$(($(SERVER_PORT) + 1)) run .

context-upload:
	bash automate/scripts/ci-upload.sh

# This is a catch-all target to prevent make from complaining
# when we pass additional arguments to our targets, like `make migrate diff`.
# It assumes that the extra arguments are for the script and not other make targets.
.PHONY: %
%:
	@# This is a deliberate empty recipe

.PHONY: dao migrate create swagger test-repo lint-install context-upload
