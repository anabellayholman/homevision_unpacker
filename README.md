# HomeVision Unpacker

This repository contains a Go parser and CLI and a browser frontend for unpacking HomeVision .env archive format.

Repository layout:
- cli/: Go module with parser and CLI
- tests/: separated unit, integration and business tests (importing the cli module)
- frontend/: browser UI with JS fallback; place main.wasm and wasm_exec.js here if you build WASM

Build and test (Linux / macOS / Windows with Go installed):
cd cli
go test ./...
go run . sample.env

Build WASM (optional):
cd cli
GOOS=js GOARCH=wasm go build -o ../frontend/main.wasm main.go parser.go
copy %GOROOT%\misc\wasm\wasm_exec.js ..\frontend\wasm_exec.js   (Windows)
cp $(go env GOROOT)/misc/wasm/wasm_exec.js ../frontend/wasm_exec.js (macOS/Linux)

Run frontend locally (no Go required):
cd frontend
python3 -m http.server 8080
open http://localhost:8080
