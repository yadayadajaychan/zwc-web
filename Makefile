export GOOS=js
export GOARCH=wasm

.PHONY: all
all:
	go build -o zwc-web.wasm zwc-web.go

.PHONY: clean
clean:
	rm zwc-web.wasm
