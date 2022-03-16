start:
	docker-compose up -d

start-build:
	docker-compose up -d --build

stop:
	docker-compose down

restart:
	make stop && make start

restart-build:
	make stop && make start-build

logs:
	 docker-compose logs -f api

# test:
# 	go test ./...

migrate:
	docker-compose up -d --build migration

seed:
	docker-compose up -d --build seed

test:
	echo 'test ./...'

download-modules:
	echo 'go mod download'
