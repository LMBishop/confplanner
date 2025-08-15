package web

import (
	"embed"
	"io/fs"
	"net/http"
	"regexp"
)

//go:generate npm install
//go:generate npm run generate

//go:embed all:.output/public
var fsys embed.FS
var urlFileRegexp = regexp.MustCompile(`[\w\-/]+\.[a-zA-Z]+$`)

type WebFileServer struct {
	root    fs.FS
	handler http.Handler
}

func NewWebFileServer() *WebFileServer {
	fsys, _ := fs.Sub(fsys, ".output/public")
	return &WebFileServer{
		root:    fsys,
		handler: http.FileServerFS(fsys),
	}
}

func (fs *WebFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p := r.URL.Path; p != "/" && !urlFileRegexp.MatchString(p) {
		http.ServeFileFS(w, r, fs.root, "index.html")
		return
	}
	fs.handler.ServeHTTP(w, r)
}
