## WASM 学习笔记

- [Go WebAssembly (Wasm) 简明教程](https://geektutu.com/post/quick-go-wasm.html)

## IDE

Goland 需要修改配置 `Go -> Build Tags & Verdoring`， `GOOS=js`, `GOARCH=wasm`, 才能正常编辑 main.go 文件

## 运行

```shell
make init
make build
```

在当前目录执行

```shell
python3 -m http.server
```

访问 http://localhost:8000/