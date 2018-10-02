NS = dev
VERSION ?= latest
NAME = protox-dev
INSTANCE = default

.PHONY: build-remote build run start stop test release shell purge irun spawn

build:
	docker build -t $(NS) .

build-remote:
	docker build -f Dockerfile.build -t $(NS) .

shell:
	docker run --rm --name $(NAME) -i -t $(PORTS) $(VOLUMES) $(ENV) $(NS) /bin/bash

run:
	docker run --rm --name $(NAME) $(PORTS) $(VOLUMES) $(ENV) $(NS)

spawn:
	docker run --rm $(PORTS) $(VOLUMES) $(ENV) -d $(NS):$(VERSION)

irun:
	docker run --rm --name $(NAME) $(PORTS) $(VOLUMES) $(ENV) -i -t $(NS)

start:
	docker run -d --name $(NAME) $(PORTS) $(VOLUMES) $(ENV) $(NS)

stop:
	docker stop $(NAME)

purge:
	-docker stop `docker ps -a -q`
	-docker rm  `docker ps -a -q`

default:
	build
