compose-all:
	docker-compose -f docker-compose.svc.yml up -d --build && \
	docker-compose -f docker-compose.svc.yml run -T wait-for-svc

compose-down:
	@docker-compose -f docker-compose.svc.yml down -t 10 --remove-orphans
