//go:build windows || ((linux || android) && cgo)

package etk

import "golang.design/x/clipboard"

func clipboardBuffer() []byte {
	return clipboard.Read(clipboard.FmtText)
}
