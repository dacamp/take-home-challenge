CONTAINERID ?= 8663488f2d94

deploy:
	docker run -d -p 127.0.0.1:1234:7777 $(CONTAINERID);
	docker run -d -p 127.0.0.1:1235:7777 $(CONTAINERID);
	docker run -d -p 127.0.0.1:1236:7777 $(CONTAINERID)
