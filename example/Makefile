.PHONY: build install

build:
	@command -v tinygo >/dev/null || (echo "错误: TinyGo未安装，请运行 make install" && exit 1)
	@tinygo build -o plugin.wasm -target=wasi -opt=z main.go
	@command -v wasm-opt >/dev/null && wasm-opt -Oz plugin.wasm -o plugin.optimized.wasm || true

install:
	@brew install tinygo binaryen
