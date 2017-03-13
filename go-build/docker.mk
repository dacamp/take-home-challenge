IMAGE_LOG_FILE = /tmp/challenge.docker
IMAGE_ID ?= $(shell tail -1 $(IMAGE_LOG_FILE)  | cut -d " " -f 3)
IMAGES = $(shell docker ps -a -q)

RUNNING_IMAGES = $(shell docker ps -a -q --filter status=running)
PAUSED_IMAGES = $(shell docker ps -a -q --filter status=paused)

RANGE =  $(words $(RUNNING_IMAGES))
NUMBER = $(shell number=$$RANDOM; let "number %= $(RANGE)" 2> /dev/null; echo "$$(($$number +1))")
RANDOM_IMAGE = $(word $(NUMBER), $(RUNNING_IMAGES) )


deploy: stop docker-build start

start: validate-docker-build
	docker run -d -p 127.0.0.1:1234:7777 $(IMAGE_ID);
	docker run -d -p 127.0.0.1:1235:7777 $(IMAGE_ID);
	docker run -d -p 127.0.0.1:1236:7777 $(IMAGE_ID);

stop:
	@if [[ "$(IMAGES)" ]]; then  docker stop $(IMAGES); fi

remove: stop
	@if [[ "$(IMAGES)" ]]; then  docker rm $(IMAGES); fi

docker-build: package
	docker build . | tee $(IMAGE_LOG_FILE)

pause: show-images
	@if [[ "$(RUNNING_IMAGES)" ]]; then image=$(RANDOM_IMAGE); echo "Pausing container: $$image" && \
		docker pause $$image; else echo "No running images"; fi

unpause:
	@if [[ "$(PAUSED_IMAGES)" ]]; then echo "Unpausing: $(PAUSED_IMAGES)" && \
		docker unpause $(PAUSED_IMAGES); else echo "No paused images"; fi

show-images:
	@echo $(IMAGES)

validate-docker-build:
ifndef IMAGE_ID
    $(error unable to determine docker image $(IMAGE_LOG_FILE) does not exist)
endif

open-debug-requests:
	open "http://localhost:1234/debug/requests" "http://localhost:1235/debug/requests" "http://localhost:1236/debug/requests"
