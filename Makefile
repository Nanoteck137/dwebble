SWAG=swag

.PHONY: sqlc build swagger

all: build

build:
	go build -o build/dwebble main.go 

build-dev:
	go build -tags dev -o build/dwebble main.go 

sqlc:
	sqlc generate

swagger:
	$(SWAG) init -g api/api.go -d .,handlers
