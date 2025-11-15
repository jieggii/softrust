.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "generate-api-client  - Generate API client from OpenAPI specification"

.PHONY: generate-api-client
generate-api-client:
	npx openapi-typescript-codegen --input ../backend/openapi.yml --output ./app/client --client axios