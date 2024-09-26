run:
	air

gen:
	go run cmd/gen/main.go
	pyrin go -o cmd/dwebble-dl/api misc/pyrin.json
	pyrin ts -o web/src/lib/api misc/pyrin.json
