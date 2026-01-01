codegen:
	cd app && oapi-codegen -config codegen-config.yml openapi.yml

dev:
	cd app && go run .
	