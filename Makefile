include .env
MIGRATE=docker-compose exec web migrate -path=migration -database "mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" -verbose

dev:
		bash automate/scripts/gin-watch.sh ${SERVER_PORT}

migrate-up:
		$(MIGRATE) up

migrate-down:
		$(MIGRATE) down

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
		${MIGRATE} create -ext sql -seq -dir migration  $$NAME

swag-generate:
		swag fmt
		swag init --parseDependency --parseInternal

crud:
	bash automate/scripts/crud.sh

inject-jwt:
	bash automate/scripts/inject_jwt.sh

.PHONY: dev migrate-up migrate-down force goto drop create swag-generate crud
