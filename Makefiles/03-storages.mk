.PHONY: sqlc-generate
sqlc-generate: ## Генерация Go-кода из SQL (sqlc)
	@echo -e "$(YELLOW)[INFO][SQLC][STARTED]$(RESET) Генерация sqlc"
	@$(SQLC) generate -f $(SQLC_CONFIG) && \
		echo -e "$(GREEN)[INFO][SQLC][SUCCESS]$(RESET) Генерация sqlc завершена" || \
		echo -e "$(RED)[ERROR][SQLC][FAIL]$(RESET) Ошибка генерации sqlc"