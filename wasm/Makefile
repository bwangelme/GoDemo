.PHONY: build init

build:
	GOOS=js GOARCH=wasm go build -o static/main.wasm
init:
	cp /home/xuyundong/.gvm/gos/go1.20.2/misc/wasm/wasm_exec.js static