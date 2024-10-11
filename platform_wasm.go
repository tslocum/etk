//go:build js && wasm

package etk

import (
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
