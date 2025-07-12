//go:build !embed_frontend

package server

import (
	"net/http"
	"os"
)

type notFoundFS struct{}

func (notFoundFS) Open(name string) (http.File, error) {
	return nil, os.ErrNotExist
}

var frontendFS http.FileSystem = notFoundFS{}
