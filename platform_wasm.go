//go:build js && wasm

package etk

import (
	"fmt"
	"syscall/js"
)

func clipboardBuffer() []byte {
	global := js.Global()
	if !global.Get("getClipboard").Truthy() {
		return nil
	}
	promise := global.Call("getClipboard", nil)
	if !promise.Truthy() {
		return nil
	}
	result := make(chan []byte, 1)
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 || !args[0].Truthy() {
			result <- nil
			return nil
		}
		result <- []byte(args[0].String())
		return nil
	}))
	return <-result
}

// Open opens a file, directory or URI using the default application registered
// in the OS to handle it. Only URIs are supported on WebAssembly.
func Open(target string) error {
	window := js.Global().Get("window")
	if !window.Truthy() {
		return fmt.Errorf("failed to get window object")
	} else if !window.Get("open").Truthy() {
		return fmt.Errorf("failed to get window.open")
	}
	windowProxy := window.Call("open", target)
	if !windowProxy.Truthy() || !windowProxy.Get("focus").Truthy() {
		return nil
	}
	windowProxy.Call("focus", nil)
	return nil
}
