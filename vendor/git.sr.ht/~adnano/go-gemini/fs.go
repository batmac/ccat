package gemini

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/url"
	"path"
	"sort"
	"strings"
)

// FileServer returns a handler that serves Gemini requests with the contents
// of the provided file system.
//
// To use the operating system's file system implementation, use os.DirFS:
//
//	gemini.FileServer(os.DirFS("/tmp"))
func FileServer(fsys fs.FS) Handler {
	return fileServer{fsys}
}

type fileServer struct {
	fs.FS
}

func (fsys fileServer) ServeGemini(ctx context.Context, w ResponseWriter, r *Request) {
	const indexPage = "/index.gmi"

	url := path.Clean(r.URL.Path)

	// Redirect .../index.gmi to .../
	if strings.HasSuffix(url, indexPage) {
		w.WriteHeader(StatusPermanentRedirect, strings.TrimSuffix(url, "index.gmi"))
		return
	}

	name := url
	if name == "/" {
		name = "."
	} else {
		name = strings.TrimPrefix(name, "/")
	}

	f, err := fsys.Open(name)
	if err != nil {
		w.WriteHeader(toGeminiError(err))
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		w.WriteHeader(toGeminiError(err))
		return
	}

	// Redirect to canonical path
	if len(r.URL.Path) != 0 {
		if stat.IsDir() {
			target := url
			if target != "/" {
				target += "/"
			}
			if len(r.URL.Path) != len(target) || r.URL.Path != target {
				w.WriteHeader(StatusPermanentRedirect, target)
				return
			}
		} else if r.URL.Path[len(r.URL.Path)-1] == '/' {
			// Remove trailing slash
			w.WriteHeader(StatusPermanentRedirect, url)
			return
		}
	}

	if stat.IsDir() {
		// Use contents of index.gmi if present
		name = path.Join(name, indexPage)
		index, err := fsys.Open(name)
		if err == nil {
			defer index.Close()
			f = index
		} else {
			// Failed to find index file
			dirList(w, f)
			return
		}
	}

	// Detect mimetype from file extension
	ext := path.Ext(name)
	mimetype := mime.TypeByExtension(ext)
	w.SetMediaType(mimetype)
	io.Copy(w, f)
}

// ServeFile responds to the request with the contents of the named file
// or directory. If the provided name is constructed from user input, it
// should be sanitized before calling ServeFile.
func ServeFile(w ResponseWriter, fsys fs.FS, name string) {
	const indexPage = "/index.gmi"

	// Ensure name is relative
	if name == "/" {
		name = "."
	} else {
		name = strings.TrimLeft(name, "/")
	}

	f, err := fsys.Open(name)
	if err != nil {
		w.WriteHeader(toGeminiError(err))
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		w.WriteHeader(toGeminiError(err))
		return
	}

	if stat.IsDir() {
		// Use contents of index file if present
		name = path.Join(name, indexPage)
		index, err := fsys.Open(name)
		if err == nil {
			defer index.Close()
			f = index
		} else {
			// Failed to find index file
			dirList(w, f)
			return
		}
	}

	// Detect mimetype from file extension
	ext := path.Ext(name)
	mimetype := mime.TypeByExtension(ext)
	w.SetMediaType(mimetype)
	io.Copy(w, f)
}

func dirList(w ResponseWriter, f fs.File) {
	var entries []fs.DirEntry
	var err error
	d, ok := f.(fs.ReadDirFile)
	if ok {
		entries, err = d.ReadDir(-1)
	}
	if !ok || err != nil {
		w.WriteHeader(StatusTemporaryFailure, "Error reading directory")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		link := LineLink{
			Name: name,
			URL:  "./" + url.PathEscape(name),
		}
		fmt.Fprintln(w, link.String())
	}
}

func toGeminiError(err error) (status Status, meta string) {
	if errors.Is(err, fs.ErrNotExist) {
		return StatusNotFound, "Not found"
	}
	if errors.Is(err, fs.ErrPermission) {
		return StatusNotFound, "Forbidden"
	}
	return StatusTemporaryFailure, "Internal server error"
}
