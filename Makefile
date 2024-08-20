.PHONY: install

bin/llpm: src/main.go
	go build -o $@ $^

DEPS = src/lib/args.go src/lib/messages.go src/lib/processManager.go src/lib/serialization.go

src/main.go: $(DEPS)

install: $(HOME)/bin/llpm

$(HOME)/bin/llpm: bin/llpm
	cp bin/llpm ~/bin
