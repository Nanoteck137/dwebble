SWAG=swag

.PHONY: sqlc

generate: sqlc

sqlc:
	sqlc generate

swagger:
	$(SWAG) init -g api/api.go -d .,handlers
