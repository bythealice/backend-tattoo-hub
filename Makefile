include .env

# Detecta migrate local (brew install golang-migrate)
MIGRATE ?= $(shell command -v migrate 2>/dev/null)
MIGRATIONS_DIR := internal/db/migrations

print-env:
	@echo "DATABASE_URL=$(DATABASE_URL)"

# ===== APPLY =====
migrate-up:
ifeq ($(MIGRATE),)
	$(MAKE) migrate-up-docker
else
	$(MAKE) migrate-up-local
endif

migrate-down:
ifeq ($(MIGRATE),)
	$(MAKE) migrate-down-docker
else
	$(MAKE) migrate-down-local
endif

migrate-up-local:
	$(MIGRATE) -path=$(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down-local:
	$(MIGRATE) -path=$(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

# Docker fallback (use host.docker.internal no DATABASE_URL)
migrate-up-docker:
	docker run --rm \
	  -v $(PWD)/$(MIGRATIONS_DIR):/migrations \
	  -e DATABASE_URL="$(DATABASE_URL)" \
	  migrate/migrate \
	  -path=/migrations -database "$${DATABASE_URL}" up

migrate-down-docker:
	docker run --rm \
	  -v $(PWD)/$(MIGRATIONS_DIR):/migrations \
	  -e DATABASE_URL="$(DATABASE_URL)" \
	  migrate/migrate \
	  -path=/migrations -database "$${DATABASE_URL}" down 1

# ===== CREATE =====
# Uso: make migrate-create name=add_users_index
migrate-create:
	@test -n "$(name)" || (echo "ERRO: passe name=..." && exit 1)
ifeq ($(MIGRATE),)
	$(MAKE) migrate-create-docker name="$(name)"
else
	$(MAKE) migrate-create-local name="$(name)"
endif

migrate-create-local:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq "$(name)"

migrate-create-docker:
	docker run --rm \
	  -v $(PWD)/$(MIGRATIONS_DIR):/migrations \
	  migrate/migrate \
	  create -ext sql -dir /migrations -seq "$(name)"

# Interativo: pergunta o nome e cria a migration
new-migration:
	@read -p "Migration name: " NAME; \
	test -n "$$NAME" || (echo "Nome vazio" && exit 1); \
	$(MAKE) migrate-create name="$$NAME"