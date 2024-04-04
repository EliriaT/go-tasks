DB_URL=mysql://user:password@tcp(localhost:3306)/sources
MIGRATIONS_DIR := db/migrations

define create_migration
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(1)
endef

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

create_migration:
	$(call create_migration,$(filter-out $@,$(MAKECMDGOALS)))