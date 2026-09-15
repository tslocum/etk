//go:build windows || (linux && !android && cgo)

package etk

import (
	"context"

	"golang.design/x/clipboard"
)

func clipboardBuffer() []byte {
	buf, err := clipboard.Read(context.Background(), clipboard.FmtText)
	if err != nil {
		return nil
	}
	return buf
}
