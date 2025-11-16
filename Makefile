include Makefiles/00-common.mk
include Makefiles/01-tools.mk
include Makefiles/02-go.mk
include Makefiles/03-storages.mk
include Makefiles/04-api.mk
include Makefiles/05-docker-dev.mk
include Makefiles/06-docker-prod.mk

.PHONY: help
help: ## Показать список доступных для использования джоб
	@echo -e "$(PURPLE)Доступные джобы:$(RESET)"
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-20s$(RESET) %s\n", $$1, $$2}'