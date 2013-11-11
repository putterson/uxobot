all: build

run: build
	uxobot

build: main.go bot.go moveslice.go
	go install

clean:
	rm uxobot
