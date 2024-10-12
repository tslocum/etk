//go:build windows || (linux && !android && cgo)

package etk

import "code.rocket9labs.com/tslocum/clipboard"

func clipboardBuffer() []byte {
	return clipboard.Read(clipboard.FmtText)
}
