# Overview

gominesweeper is a fun project that implements [minesweeper](https://en.wikipedia.org/wiki/Minesweeper_(video_game)) using Golang + WASM

# Prerequisite
- Golang 1.15
- depends on syscall/js package which is still in alpha and can have breaking changes across releases of Go

# Compile and run

- Clone the workspace 
- Compile the WASM using
```
GOOS=js GOARCH=wasm go build -o resources/game.wasm
```
- run the webserver
```
go run ./server/server.go -listen=:8080
```
- now point your browser to http://localhost:8080 and enjoy sweeping mines :)