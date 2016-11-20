IMAGE_LOG_FILE = /tmp/challenge.docker
IMAGE_ID ?= $(shell tail -1 $(IMAGE_LOG_FILE)  | cut -d " " -f 3)

stop:
	docker stop `docker ps -a -q`

remove: stop
	docker rm `docker ps -a -q`

docker:
	docker build . | tee $(IMAGE_LOG_FILE)

validate-docker-build:
ifndef IMAGE_ID
    $(error unable to determine docker image $(IMAGE_LOG_FILE) does not exist)
endif

deploy: docker validate-docker-build
	docker run -d -p 127.0.0.1:1234:7777 $(IMAGE_ID);
	docker run -d -p 127.0.0.1:1235:7777 $(IMAGE_ID);
	docker run -d -p 127.0.0.1:1236:7777 $(IMAGE_ID);

show-requests:
	open "http://localhost:1234/debug/requests?fam=main.counterHandler&b=0"
