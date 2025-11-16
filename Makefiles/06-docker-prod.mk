.PHONY: prod-up
prod-up: ## Поднятие prod окружения
	@echo -e "$(YELLOW)[INFO][PROD-UP][STARTED]$(RESET) Запуск prod окружения"
	@docker compose -f $(DOCKER_PROD_COMPOSE) up -d --build && \
		echo -e "$(GREEN)[INFO][PROD-UP][SUCCESS]$(RESET) Prod окружение успешно запущено" || \
		echo -e "$(RED)[ERROR][PROD-UP][FAIL]$(RESET) Ошибка при запуске prod окружения"

.PHONY: prod-down
prod-down: ## Остановка и удаление контейнеров prod окружения
	@echo -e "$(YELLOW)[INFO][PROD-DOWN][STARTED]$(RESET) Остановка prod окружения"
	@docker compose -f $(DOCKER_PROD_COMPOSE) down && \
		echo -e "$(GREEN)[INFO][PROD-DOWN][SUCCESS]$(RESET) Prod окружение остановлено" || \
		echo -e "$(RED)[ERROR][PROD-DOWN][FAIL]$(RESET) Ошибка при остановке prod окружения"

.PHONY: prod-logs
prod-logs: ## Просмотр логов prod окружения
	@echo -e "$(YELLOW)[INFO][PROD-LOGS][STARTED]$(RESET) Просмотр логов контейнеров prod окружения"
	@docker compose -f $(DOCKER_PROD_COMPOSE) logs -f

.PHONY: prod-restart
prod-restart: ## Перезапуск prod окружения
	@$(MAKE) prod-down && $(MAKE) prod-up
