include .env
export
RUNNER=docker-compose exec web sql-migrate

ifeq ($(p),host)
 	RUNNER=sql-migrate
endif

MIGRATE=$(RUNNER)

dev:
		bash automate/scripts/gin-watch.sh ${SERVER_PORT}

create:
	@read -p  "What is the name of migration?" NAME; \
	$(MIGRATE) new $$NAME

migrate-status:
		$(MIGRATE) status

migrate-up:
		$(MIGRATE) up

migrate-down:
		$(MIGRATE) down 

redo:
	@read -p  "Are you sure to reapply the last migration? [y/n]" -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; \
	then \
		$(MIGRATE) redo; \
	fi

swag-generate:
		swag fmt
		swag init --parseDependency --parseInternal

crud:
	bash automate/scripts/crud.sh

create-app:
	bash automate/scripts/new_app.sh

create-version:
	bash automate/scripts/new_version.sh

inject-auth:
	bash automate/scripts/inject_auth.sh

install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2
	git config core.hooksPath hooks

start: install
	docker-compose up

run:
	docker-compose up


.PHONY: create migrate-status migrate-up migrate-down redo 
