SHELL = /bin/sh
ROOT=$(shell pwd)
export VERSION := $(version)

.PHONY: up-clean
up-clean: docker-down docker-rebuild docker-up

docker-up:
	docker-compose up -d
docker-down:
	docker-compose kill
	docker-compose rm -f
docker-clean:
	docker image rm forklol_backend:latest | true
	docker image rm forklol_frontend:latest | true

.PHONY: docker-rebuild
docker-rebuild: 
	docker-compose build --force-rm

.PHONY: mysql-clean
mysql-clean:
	rm -rf ./docker/mysql/data/*

.PHONY: docker-prune
docker-prune:
	docker system prune -a

