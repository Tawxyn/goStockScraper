run:
	@npm run build:css
	@templ generate
	@go run cmd/app/main.go