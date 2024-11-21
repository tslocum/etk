//go:build linux

package etk

import "os/exec"

// Open opens a file, directory or URI using the default application registered
// in the OS to handle it. Only URIs are supported on WebAssembly.
func Open(target string) error {
	cmd := exec.Command("xdg-open", target)
	return cmd.Start()
}
