BIN := /usr/local/bin

install: build;
	cp -f git-secrets-hooks-cleaner $(BIN)

build:
	go build

uninstall:
	rm -f $(BIN)/git-secrets-hooks-cleaner
