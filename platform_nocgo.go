//go:build android || (!windows && !js && !wasm && !cgo)

package etk

func clipboardBuffer() []byte {
	return nil
}
