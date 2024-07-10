.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	sudo docker build -t pocketer:v0.1 .

start-container:
	sudo docker run --name pocketer -p 80:80  --env-file .env pocketer:v0.1