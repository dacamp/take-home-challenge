CONTAINERID ?= 8663488f2d94

stop:
	docker stop `docker ps -a -q`

deploy:
	docker run -d -p 127.0.0.1:1234:7777 $(CONTAINERID);
	docker run -d -p 127.0.0.1:1235:7777 $(CONTAINERID);
	docker run -d -p 127.0.0.1:1236:7777 $(CONTAINERID)

show-requests:
	open "http://localhost:1234/debug/requests?fam=main.counterHandler&b=0"
