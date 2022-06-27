include .env
export
RUNNER=docker-compose exec web sql-migrate

ifeq ($(p),host)
 	RUNNER=sql-migrate
endif

MIGRATE=$(RUNNER)

dev:
		gin appPort ${ServerPort} -i run server.go

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

force:
		@read -p  "Which version do you want to force?" VERSION; \
		$(MIGRATE) force $$VERSION

goto:
		@read -p  "Which version do you want to migrate?" VERSION; \
		$(MIGRATE) goto $$VERSION

drop:
		$(MIGRATE) drop

create:
		@read -p  "What is the name of migration?" NAME; \
		$(MIGRATE) new $$NAME

crud:
		bash automate/scripts/crud.sh

.PHONY: migrate-up migrate-down force goto drop create

.PHONY: migrate-up migrate-down force goto drop create auto-create
