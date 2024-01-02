.PHONY: css views sqlc

generate: css views

css:
	npx tailwindcss -i ./style.css -o ./public/style.css

views:
	templ generate

sqlc:
	sqlc generate

