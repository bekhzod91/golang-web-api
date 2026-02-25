include .env

swag:
	rm -rf src/infrastructure/api/dto/*
	swagger generate model --spec=src/openapi/_swagger.yaml --target=src/infrastructure/api --model-package=dto
	swagger flatten src/openapi/_swagger.yaml --output=src/static/swagger/swagger.json

swag-validate:
	swagger validate src/openapi/_swagger.yaml

run:
	make swag && make sqlc && cd src && go run cmd/server/main.go

migration-shared:
	migrate create -ext sql -dir src/infrastructure/migrations/shared $(title)

migration-tenant:
	migrate create -ext sql -dir src/infrastructure/migrations/tenant $(title)

migrate:
	cd src && go run cmd/migrate/main.go

test:
	cd src \
	&& gotestsum --format testname -- ./tests/test_auth/ \
	&& gotestsum --format testname -- ./tests/test_role/ \
	&& gotestsum --format testname -- ./tests/test_user/ \
	&& gotestsum --format testname -- ./tests/test_permissions/ \

sqlc:
	cd src && sqlc generate
