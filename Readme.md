# Overview

gominesweeper is a fun project of mine that implements [minesweeper](https://en.wikipedia.org/wiki/Minesweeper_(video_game)) using Golang + WASM

![gomineseeper in action](/snapshots/sc1.png?raw=true "gomineseeper in action")

# Prerequisite
- Golang 1.15
- Depends on syscall/js package which is still in alpha and can have breaking changes across releases of Go

# Compile and run

- Clone the workspace 
- Compile the WASM
```
GOOS=js GOARCH=wasm go build -o resources/game.wasm
```
- Run the webserver
```
go run ./server/server.go -listen=:8080
```
- Point your browser to http://localhost:8080 and enjoy sweeping some mines :)