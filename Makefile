include .envrc

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@echo '--Running application--'
	@go run ./cmd/api --port=$(PORT) --env=$(ENV) --db-dsn=$(TRAINING_DB_DSN) --cors-trusted-origins="$(CORS_TRUSTED_ORIGINS)" --limiter-burst=5 --limiter-rps=2 --limiter-enabled=true --db-max-open-conns=50 --db-max-idle-conns=50  --db-max-idle-time=2h30m



## db/sql:connect to the database using psql(terminal)
.PHONY: db/psql
db/psql:
	psql ${TRAINING_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration filed for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

##db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${TRAINING_DB_DSN} up


## db/migrations/fix: automatically fix dirty migrations
.PHONY: db/migrations/fix
db/migrations/fix:
	@echo 'Checking migration status...'
	@mkdir -p /tmp
	@{ \
        migrate -path ./migrations -database "${TRAINING_DB_DSN}" version > /tmp/migrate_version 2>&1; \
        cat /tmp/migrate_version; \
        if grep -q "dirty" /tmp/migrate_version; then \
            version=$$(grep -o '[0-9]\+' /tmp/migrate_version | head -1); \
            echo "⚠️ Found dirty migration at version $$version"; \
            echo "Forcing version $$version..."; \
            migrate -path ./migrations -database "${TRAINING_DB_DSN}" force $$version; \
            echo "Running down migration..."; \
            migrate -path ./migrations -database "${TRAINING_DB_DSN}" down 1; \
            echo "Running up migration..."; \
            migrate -path ./migrations -database "${TRAINING_DB_DSN}" up 1; \
        else \
            echo "✅ No dirty migration found"; \
        fi; \
        rm -f /tmp/migrate_version; \
    }