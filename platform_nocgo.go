//go:build !windows && !js && !wasm && !cgo

package etk

func clipboardBuffer() []byte {
	return nil
}
