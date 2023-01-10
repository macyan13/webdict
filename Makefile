race_test:
	cd backend && go test -race -mod=vendor -timeout=60s -count 1 ./...

backend:
	docker compose -f compose-dev-backend.yml build
	docker compose -f compose-dev-backend.yml up -d

stop:
	docker compose -f compose-dev-backend.yml stop

.PHONY: backend stop
