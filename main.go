package main

import "syscall/js"

func main() {
	registerCallbacks()
	initGame(js.Null(), []js.Value{js.ValueOf("1")})
	println("go wasm initialized")
	// block for ever, so wasm functions stay available
	select {}
}
