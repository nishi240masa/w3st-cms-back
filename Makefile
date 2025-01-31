include .env


up: build
	docker-compose up -d

build:
	docker-compose build

down:
	docker-compose down

go:
	docker-compose exec -it w3st-cms /bin/sh

db:
	docker exec -it ${DB_HOST} psql -U ${DB_USER} -d ${DB_NAME}