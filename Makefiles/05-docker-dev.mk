.PHONY: dev-up
dev-up: ## Поднятие dev окружения
	@echo -e "$(YELLOW)[INFO][DEV-UP][STARTED]$(RESET) Запуск dev окружения"
	@docker compose -f $(DOCKER_DEV_COMPOSE) up --build -d && \
		echo -e "$(GREEN)[INFO][DEV-UP][SUCCESS]$(RESET) Dev окружение успешно запущено" || \
		echo -e "$(RED)[ERROR][DEV-UP][FAIL]$(RESET) Ошибка при запуске dev окружения"

.PHONY: dev-down
dev-down: ## Остановка и удаление контейнеров dev окружения
	@echo -e "$(YELLOW)[INFO][DEV-DOWN][STARTED]$(RESET) Остановка dev окружения"
	@docker compose -f $(DOCKER_DEV_COMPOSE) down && \
		echo -e "$(GREEN)[INFO][DEV-DOWN][SUCCESS]$(RESET) Dev окружение остановлено" || \
		echo -e "$(RED)[ERROR][DEV-DOWN][FAIL]$(RESET) Ошибка при остановке dev окружения"

.PHONY: dev-logs
dev-logs: ## Просмотр логов dev окружения
	@echo -e "$(YELLOW)[INFO][DEV-LOGS][STARTED]$(RESET) Просмотр логов контейнеров dev окружения"
	@docker compose -f $(DOCKER_DEV_COMPOSE) logs -f

.PHONY: dev-restart
dev-restart: ## Перезапуск dev окружения
	@$(MAKE) dev-down && $(MAKE) dev-up
