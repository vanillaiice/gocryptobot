# Determine the operating system
ifeq ($(OS), Windows_NT)
	OS_TARGET = Windows
else
	OS_TARGET := $(shell uname -s)
endif

# Compile depending on OS
os: $(OS_TARGET)

# Compile for windows
Windows:
	mkdir -p bin/windows
	CGO_ENABLED=1 GOOS=windows go build -ldflags="-s -w" -o bin/windows/gocryptobot.exe cmd/gocryptobot/main.go

# Compile for linux
Linux:
	mkdir -p bin/linux
	CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o bin/linux/gocryptobot-linux cmd/gocryptobot/main.go

# Compile for darwin
Darwin:
	mkdir -p bin/darwin
	CGO_ENABLED=1 GOOS=darwin go build -ldflags="-s -w" -o bin/darwin/gocryptobot-darwin cmd/gocryptobot/main.go

# Compile for all OS
all: Windows Linux Darwin
