PREFIX = /usr/local
BINARY = fakebuild

all: build

build:
	go build -o $(BINARY) .

install: build
	install -Dm755 $(BINARY) $(PREFIX)/bin/$(BINARY)

clean:
	rm -f $(BINARY)

.PHONY: all build install clean
