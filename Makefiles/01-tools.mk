.PHONY: install-all-tools
install-all-tools: ## Установить все тулзы (форматирование, линтинг, OpenAPI)
	@echo -e "$(YELLOW)[INFO][DEPS-ALL][STARTED]$(RESET) Установка всех необходимых тулзов"
	@$(MAKE) install-format-tools
	@$(MAKE) install-api-tools
	@echo -e "$(GREEN)[INFO][DEPS-ALL][SUCCESS]$(RESET) Все тулзы успешно установлены"

.PHONY: install-format-tools
install-format-tools: ## Установка тулзов для форматирования и линта кода (gofumpt, golines, gci, golangci-lint)
	@echo -e "$(GREEN)[INFO][DEPS-TOOLS][STARTED]$(RESET) Установка тулзов для форматирования и линта кода"
	@for tool in $(FORMAT_TOOLS); do \
		echo -e "$(PURPLE)  - Устанавливается $$tool$(RESET)"; \
		go install $$tool || { \
			echo -e "$(RED)[ERROR][DEPS-TOOLS][FAIL]$(RESET) Ошибка при установке $$tool"; exit 1; }; \
	done
	@echo -e "$(GREEN)[INFO][DEPS-TOOLS][SUCCESS]$(RESET) Все тулзы установлены"

.PHONY: install-api-tools
install-api-tools: ## Установка тулзов для OpenAPI (oapi-codegen, vacuum, oasdiff)
	@echo -e "$(GREEN)[INFO][DEPS-API-TOOLS][STARTED]$(RESET) Установка API тулзов (oapi-codegen, vacuum, oasdiff)"
	@for tool in $(API_TOOLS); do \
		echo -e "$(PURPLE)  - Устанавливается $$tool$(RESET)"; \
		go install $$tool || { \
			echo -e "$(RED)[ERROR][DEPS-API-TOOLS][FAIL]$(RESET) Ошибка при установке $$tool"; exit 1; }; \
	done
	@echo -e "$(GREEN)[INFO][DEPS-API-TOOLS][SUCCESS]$(RESET) API тулзы успешно установлены"
