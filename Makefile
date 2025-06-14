DB_DSN=sqlite3://db.sqlite
DB_PATH=db.sqlite
MIGRATIONS_DIR=migrations

# =================================================================================== #
# HELPERS
# =================================================================================== #

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## confirm: confirm before running a command
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# =================================================================================== #
# DB
# =================================================================================== #

## db/migrations/new: create a new database migration
## Usage: make db/migrations/new name=<migration_name>
.PHONY: db/migrations/new
db/migrations/new:
	@test $(name) || (echo "Error: name argument required" && exit 1)
	@echo creating migration files for ${name}...
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
	
## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo Running up migrations...
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_DSN) up

## db/migrations/down: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo Running down migrations...
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_DSN) down

## db/migrations/reset: reset the database and apply all migrations
.PHONY: db/migrations/reset
db/migrations/reset: confirm
	@rm -f $(DB_PATH)
	@make db/migrations/up

## db/migrations/force: force a migration to a specific version
## Usage: make db/migrations/force version=<version>
.PHONY: db/migrations/force
db/migrations/force:
	@test $(version) || (echo "Error: version argument required" && exit 1)
	@echo forcing migration to version $(version)...
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_DSN) force $(version)


# =================================================================================== #
# WEB (Frontend + Embed Setup)
# =================================================================================== #

## web/build: build frontend assets and move them into web/server/dist for embedding
.PHONY: web/build
web/build:
	@echo "Building frontend..."
	@cd web/frontend && pnpm install && pnpm run build
	@rm -rf web/server/dist
	@cp -r web/frontend/dist web/server/dist

## web/clean: remove built frontend assets
.PHONY: web/clean
web/clean:
	@echo "Cleaning frontend build artifacts..."
	@rm -rf web/frontend/dist web/server/dist

# =================================================================================== #
# APP
# =================================================================================== #

## build: build the Go application (after embedding dist)
.PHONY: build
build: web/build
	@echo "Building Go binary..."
	@go build -o worklogger .

## run: build and run the app
.PHONY: run
run: build
	@echo "Running app..."
	./worklogger studio

