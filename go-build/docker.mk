IMAGE_LOG_FILE = /tmp/challenge.docker
IMAGE_ID ?= $(shell tail -1 $(IMAGE_LOG_FILE)  | cut -d " " -f 3)
IMAGES = $(shell docker ps -a -q)

deploy: stop docker-build start

start: validate-docker-build
	docker run -d -p 127.0.0.1:1234:7777 $(IMAGE_ID);
	docker run -d -p 127.0.0.1:1235:7777 $(IMAGE_ID);
	docker run -d -p 127.0.0.1:1236:7777 $(IMAGE_ID);

stop:
	@if [[ "$(IMAGES)" ]]; then  docker stop $(IMAGES); fi

remove: stop
	@if [[ "$(IMAGES)" ]]; then  docker rm $(IMAGES); fi

docker-build:
	docker build . | tee $(IMAGE_LOG_FILE)

validate-docker-build:
ifndef IMAGE_ID
    $(error unable to determine docker image $(IMAGE_LOG_FILE) does not exist)
endif
