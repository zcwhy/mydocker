.PHONY:build test

run: 
	sudo ./docker run -t -i /bin/bash && echo $?

build:
	go build -o docker main.go

test:
	@echo "Hello world"