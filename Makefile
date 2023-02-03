ifeq ($(OS),Windows_NT)
    SHELL := powershell.exe #change shell for windows
    .SHELLFLAGS := -Command
    ending := exe
else
    ending := out
endif


run:
	go run ./cmd

build:
	go build -o ./build/main.$(ending) ./cmd

fmt: 
	go fmt ./cmd ./internal/db ./internal/shortlinkgen ./internal/api