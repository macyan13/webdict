.DEFAULT_GOAL := help
.PHONY: *

GITTAG=$(shell git describe --abbrev=0 --tags)

help: ## Display available commands
	@printf "\033[0;36mAvailable commands:\033[0m\n"
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
	for help_line in $${help_lines[@]}; do \
		IFS=$$'#' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//' -e 's/:.*//' -e 's/://'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		if [ "$$help_command" = "" ]; then \
			printf "\033[0;35m %-15s \033[0m \t\n" $$help_info ; \
		else \
			printf "\033[0;32m  %-20s \033[0m \t %s\n" $$help_command $$help_info ; \
		fi; \
	done

race_test: ## Runs unit tests with -race parameter
	cd backend && go test -race -mod=vendor -timeout=60s -count 1 ./...

backend: ## Runs containers from compose-dev-backend.yml
	docker compose -f compose-dev-backend.yml build
	docker compose -f compose-dev-backend.yml up -d

stop: ## Stops ran containers
	docker compose -f compose-dev-backend.yml stop

clean: ## Stops ran containers and destroys
	docker compose -f compose-dev-backend.yml stop && docker compose -f compose-dev-backend.yml rm -f

restart: ## Stops ran containers, destroys them and restart
	docker compose -f compose-dev-backend.yml stop ;\
 	docker compose -f compose-dev-backend.yml rm -f ;\
	docker compose -f compose-dev-backend.yml build  ;\
	docker compose -f compose-dev-backend.yml up -d

release_latest: ## Creates and pushes imaged to docker hub using the last git tag
	- docker buildx build --push --platform linux/amd64,linux/arm/v7,linux/arm64 \
 		-t macyan/webdict:${GITTAG} -t macyan/webdict:latest .