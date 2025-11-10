//go:build dev
// +build dev

package web

import (
	"io/fs"
	"os"
)

var assets fs.FS = os.DirFS("web")

func Assets() fs.FS {
	return assets
}
