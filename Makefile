include .env
MIGRATE=docker-compose exec web migrate -path=migration -database "mysql://${DBUsername}:${DBPassword}@tcp(${DBHost}:${DBPort})/${DBName}" -verbose

dev:
		gin appPort ${ServerPort} -i run server.go
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

crud:
	@read -p "what is the name of crud?" FILENAME; \
	bash automate/scripts/crud.sh $$FILENAME

.PHONY: migrate-up migrate-down force goto drop create

