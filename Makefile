run: build
	docker run --name web_crawler -d -p 9000:9000 web_crawler

build: rm
	docker build . -t web_crawler

rm: stop
	-docker rm web_crawler

stop:
	-docker stop web_crawler