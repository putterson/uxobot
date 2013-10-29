all: build

run: build
	./uxobot

build: uxobot.go
	go build

clean:
	rm uxobot
