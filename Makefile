include .env

migrate-up :
	migrate -source file:db/postgres/migrations -database postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable up $(N)

migrate-down :
	migrate -source file:db/postgres/migrations -database postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable down $(N)