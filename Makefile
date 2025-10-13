# Makefile

# Подтягиваем SHELL
SHELL := C:/Program Files/Git/bin/bash.exe

# Подтягиваем переменные из .env
ifneq (,$(wildcard ./.env))
include .env
export
endif

DB_DSN := "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

# Таргет для создания новой миграции
migrate-new:
	migrate create -ext sql -dir ./migrations ${NAME}

# Применение миграций
migrate-up:
	$(MIGRATE) up

# Откат миграций
migrate-down:
	$(MIGRATE) down

# Папки для генерации API
TAGS=ofreelancers olikes omatches oprojects oreviews ousers

# Генерация API
gen:
# Очиска и пересоздание папок
	for tag in $(TAGS); do \
		rm -rf ./internal/web/$$tag/; \
		mkdir -p ./internal/web/$$tag/; \
	done

# Герерация кода по каждому тегу
	oapi-codegen -config openapi/.openapi -include-tags freelancers -package ofreelancers openapi/openapi.yaml > ./internal/web/ofreelancers/api.gen.go
	oapi-codegen -config openapi/.openapi -include-tags likes -package olikes openapi/openapi.yaml > ./internal/web/olikes/api.gen.go
	oapi-codegen -config openapi/.openapi -include-tags matches -package omatches openapi/openapi.yaml > ./internal/web/omatches/api.gen.go
	oapi-codegen -config openapi/.openapi -include-tags projects -package oprojects openapi/openapi.yaml > ./internal/web/oprojects/api.gen.go
	oapi-codegen -config openapi/.openapi -include-tags reviews -package oreviews openapi/openapi.yaml > ./internal/web/oreviews/api.gen.go
	oapi-codegen -config openapi/.openapi -include-tags users -package ousers openapi/openapi.yaml > ./internal/web/ousers/api.gen.go

# Запуск приложения
run:
	go run cmd/main.go
