compose:
	docker-compose -f docker-compose.yml up -d --build && \
	docker-compose -f docker-compose.yml run -T wait-for-infra && \
	docker-compose -f docker-compose.seed.yml up -d --build

compose-all:
	docker-compose -f docker-compose.yml up -d --build && \
	docker-compose -f docker-compose.yml run -T wait-for-infra && \
	docker-compose -f docker-compose.seed.yml up -d --build && \
	docker-compose -f docker-compose.svc.yml up -d --build && \
	docker-compose -f docker-compose.svc.yml run -T wait-for-svc

compose-down:
	@docker-compose -f docker-compose.yml -f docker-compose.seed.yml -f docker-compose.svc.yml down -t 10 --remove-orphans

unit-test:
	docker run --rm \
		-v $(PWD):/app \
		-w /app \
		golang:1.24.3-alpine \
		sh -c '\
			set -e ; \
			go test -v ./pkg/... ; \
			for mod in $$(find services -name go.mod); do \
				dir=$$(dirname $$mod); \
				echo "Running tests in $$dir"; \
				(cd $$dir && go test -v ./...); \
			done \
		'

integration-test:
	docker run --rm \
		--network backend-network \
		--env-file ./tests/.env.dist \
		-v $(PWD):/app \
		-w /app \
		golang:1.24.3-alpine \
		sh -c "go test -v ./tests/..."
