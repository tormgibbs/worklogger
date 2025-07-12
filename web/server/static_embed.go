//go:build embed_frontend

package server

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var embeddedFiles embed.FS

var frontendFS http.FileSystem

func init() {
	subFS, err := fs.Sub(embeddedFiles, "dist")
	if err != nil {
		panic("failed to load embedded frontend: " + err.Error())
	}
	frontendFS = http.FS(subFS)
}

