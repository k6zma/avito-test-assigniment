.PHONY: api-regen
api-regen: ## Сгенерировать код API (oapi-codegen v2)
	@echo -e "$(YELLOW)[INFO][API-REGEN][STARTED]$(RESET) Генерация кода из OpenAPI"
	@mkdir -p $(OPENAPI_GEN_DIR)

	@echo -e "$(PURPLE)[INFO][API-REGEN]$(RESET) Генерация: модели"
	@$(OAPI_CODEGEN) -config $(OPENAPI_GEN_CONFIG_DIR)/$(OPENAPI_GEN_MODELS) $(OPENAPI_FILE) || { \
		echo -e "$(RED)[ERROR][API-REGEN][FAIL]$(RESET) Ошибка генерации models"; exit 1; }

	@echo -e "$(PURPLE)[INFO][API-REGEN]$(RESET) Генерация: сервер"
	@$(OAPI_CODEGEN) -config $(OPENAPI_GEN_CONFIG_DIR)/$(OPENAPI_GEN_SERVER) $(OPENAPI_FILE) || { \
		echo -e "$(RED)[ERROR][API-REGEN][FAIL]$(RESET) Ошибка генерации server"; exit 1; }

	@echo -e "$(PURPLE)[INFO][API-REGEN]$(RESET) Генерация: spec"
	@$(OAPI_CODEGEN) -config $(OPENAPI_GEN_CONFIG_DIR)/$(OPENAPI_GEN_SPEC) $(OPENAPI_FILE) || { \
		echo -e "$(RED)[ERROR][API-REGEN][FAIL]$(RESET) Ошибка генерации spec"; exit 1; }

	@echo -e "$(GREEN)[INFO][API-REGEN][SUCCESS]$(RESET) Код успешно сгенерирован"

.PHONY: api-lint
api-lint: ## Линтинг OpenAPI спецификации (vacuum)
	@echo -e "$(YELLOW)[INFO][API-LINT][STARTED]$(RESET) Линтинг OpenAPI через vacuum"
	@$(VACUUM) lint -r $(OPENAPI_LINT_RULES) $(OPENAPI_FILE) && \
		echo -e "$(GREEN)[INFO][API-LINT][SUCCESS]$(RESET) Линт прошел успешно" || \
		echo -e "$(RED)[ERROR][API-LINT][FAIL]$(RESET) Ошибки в OpenAPI спецификации"

.PHONY: api-lint-ui
api-lint-ui: ## Интерактивный дашборд OpenAPI (vacuum dashboard)
	@echo -e "$(YELLOW)[INFO][API-LINT-UI][STARTED]$(RESET) Запуск интерактивного дашборда vacuum"
	@$(VACUUM) dashboard -r $(OPENAPI_LINT_RULES) $(OPENAPI_FILE) || \
		echo -e "$(RED)[ERROR][API-LINT-UI][FAIL]$(RESET) Ошибка при запуске vacuum dashboard"
