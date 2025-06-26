run:
	@go run main.go

.PHONY: db
db:
	docker-compose up -d db

.PHONY: down
down:
	docker-compose down db