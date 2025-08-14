package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:generate npm ci
//go:generate npm run generate

//go:embed all:.output/public
var fsys embed.FS

type WebFileServer struct {
	server http.Handler
}

func NewWebFileServer() *WebFileServer {
	fsys, _ := fs.Sub(fsys, ".output/public")
	return &WebFileServer{
		server: http.FileServerFS(fsys),
	}
}

func (fs *WebFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs.server.ServeHTTP(w, r)
}
