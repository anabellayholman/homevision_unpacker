# HomeVision Unpacker

This repository contains a Go parser and CLI and a browser frontend for unpacking HomeVision `.env` archive format.  
It allows you to extract files from HomeVision `.env` archives and view them in the browser or save them locally.

## What the frontend does
The `frontend` folder contains a web UI that lets you:

- Upload a `.env` archive.
- Automatically parse and display the contents.
- Download individual extracted files.
- Work fully in the browser using WebAssembly (WASM) without needing Go installed.

## Repository layout
- `cli/`: Go module with CLI executable
- `pkg/cli/`: Go package with the parser logic
- `tests/`: unit, integration, and business tests (importing the `cli` module)
- `frontend/`: browser UI with WASM support; place `main.wasm` and `wasm_exec.js` here if you build WASM
- `sample.env`: example `.env` archive to test parsing

## Build WASM (optional, for frontend)
cd cli
GOOS=js GOARCH=wasm go build -o ../frontend/main.wasm main.go parser.go
# Windows
copy %GOROOT%\misc\wasm\wasm_exec.js ..\frontend\wasm_exec.js
# macOS/Linux
cp $(go env GOROOT)/misc/wasm/wasm_exec.js ../frontend/wasm_exec.js

## Run frontend locally (no Go required)
cd frontend
python3 -m http.server 8080

## Optional: View via GitHub Pages
https://anabellayholman.github.io/homevision_unpacker/

## Build and test (Linux / macOS / Windows with Go installed)
```bash
cd cli
go test ./...
go run . sample.env
