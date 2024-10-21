build:
	docker build -t forum:1 .
run:
	docker run --rm -p 8080:8080 forum:1