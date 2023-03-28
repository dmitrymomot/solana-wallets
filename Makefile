.PHONY: help, migrate

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
%:
	@:

prepare-migrations: ## Prepare migrations: clean old files and copy all migrations from services to migrations folder.
	@rm -Rvf migrations/*.sql
	@cp -Rvf ./svc/**/repository/sql/migrations/*.sql migrations/
