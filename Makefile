.PHONY:build test

run: 
	sudo ./docker run -t -i /bin/sh && echo $?

build:
	go build -o docker main.go

test:
	@echo "Hello world"