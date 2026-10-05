// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-GO file.

package gemini

import (
	"context"
	"net"
	"net/url"
	"path"
	"sort"
	"strings"
	"sync"
)

// Mux is a Gemini request multiplexer.
// It matches the URL of each incoming request against a list of registered
// patterns and calls the handler for the pattern that
// most closely matches the URL.
//
// Patterns name fixed, rooted paths, like "/favicon.ico",
// or rooted subtrees, like "/images/" (note the trailing slash).
// Longer patterns take precedence over shorter ones, so that
// if there are handlers registered for both "/images/"
// and "/images/thumbnails/", the latter handler will be
// called for paths beginning "/images/thumbnails/" and the
// former will receive requests for any other paths in the
// "/images/" subtree.
//
// Note that since a pattern ending in a slash names a rooted subtree,
// the pattern "/" matches all paths not matched by other registered
// patterns, not just the URL with Path == "/".
//
// Patterns may optionally begin with a host name, restricting matches to
// URLs on that host only. Host-specific patterns take precedence over
// general patterns, so that a handler might register for the two patterns
// "/search" and "search.example.com/" without also taking over requests
// for "gemini://example.com/".
//
// Wildcard patterns can be used to match multiple hostnames. For example,
// the pattern "*.example.com" will match requests for "blog.example.com"
// and "gemini.example.com", but not "example.org".
//
// If a subtree has been registered and a request is received naming the
// subtree root without its trailing slash, Mux redirects that
// request to the subtree root (adding the trailing slash). This behavior can
// be overridden with a separate registration for the path without
// the trailing slash. For example, registering "/images/" causes Mux
// to redirect a request for "/images" to "/images/", unless "/images" has
// been registered separately.
//
// Mux also takes care of sanitizing the URL request path and
// redirecting any request containing . or .. elements or repeated slashes
// to an equivalent, cleaner URL.
type Mux struct {
	mu sync.RWMutex
	m  map[hostpath]Handler
	es []muxEntry // slice of entries sorted from longest to shortest
}

type hostpath struct {
	host string
	path string
}

type muxEntry struct {
	handler Handler
	host    string
	path    string
}

// cleanPath returns the canonical path for p, eliminating . and .. elements.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		// Fast path for common case of p being the string we want:
		if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
			np = p
		} else {
			np += "/"
		}
	}
	return np
}

// Find a handler on a handler map given a path string.
// Most-specific (longest) pattern wins.
func (mux *Mux) match(host, path string) Handler {
	// Check for exact match first.
	if h, ok := mux.m[hostpath{host, path}]; ok {
		return h
	}

	// Check for longest valid match.  mux.es contains all patterns
	// that end in / sorted from longest to shortest.
	for _, e := range mux.es {
		if len(e.host) == len(host) && e.host == host &&
			strings.HasPrefix(path, e.path) {
			return e.handler
		}
	}
	return nil
}

// redirectToPathSlash determines if the given path needs appending "/" to it.
// This occurs when a handler for path + "/" was already registered, but
// not for path itself. If the path needs appending to, it creates a new
// URL, setting the path to u.Path + "/" and returning true to indicate so.
func (mux *Mux) redirectToPathSlash(host, path string, u *url.URL) (*url.URL, bool) {
	mux.mu.RLock()
	shouldRedirect := mux.shouldRedirectRLocked(host, path)
	mux.mu.RUnlock()
	if !shouldRedirect {
		return u, false
	}
	return u.ResolveReference(&url.URL{Path: path + "/"}), true
}

// shouldRedirectRLocked reports whether the given path and host should be redirected to
// path+"/". This should happen if a handler is registered for path+"/" but
// not path -- see comments at Mux.
func (mux *Mux) shouldRedirectRLocked(host, path string) bool {
	if _, exist := mux.m[hostpath{host, path}]; exist {
		return false
	}

	n := len(path)
	if n == 0 {
		return false
	}
	if _, exist := mux.m[hostpath{host, path + "/"}]; exist {
		return path[n-1] != '/'
	}

	return false
}

func getWildcard(hostname string) (string, bool) {
	if net.ParseIP(hostname) == nil {
		split := strings.SplitN(hostname, ".", 2)
		if len(split) == 2 {
			return "*." + split[1], true
		}
	}
	return "", false
}

// Handler returns the handler to use for the given request, consulting
// r.URL.Scheme, r.URL.Host, and r.URL.Path. It always returns a non-nil handler. If
// the path is not in its canonical form, the handler will be an
// internally-generated handler that redirects to the canonical path. If the
// host contains a port, it is ignored when matching handlers.
func (mux *Mux) Handler(r *Request) Handler {
	// Disallow non-Gemini schemes
	if r.URL.Scheme != "gemini" {
		return NotFoundHandler()
	}

	host := r.URL.Hostname()
	path := cleanPath(r.URL.Path)

	// If the given path is /tree and its handler is not registered,
	// redirect for /tree/.
	if u, ok := mux.redirectToPathSlash(host, path, r.URL); ok {
		return StatusHandler(StatusPermanentRedirect, u.String())
	}

	if path != r.URL.Path {
		u := *r.URL
		u.Path = path
		return StatusHandler(StatusPermanentRedirect, u.String())
	}

	mux.mu.RLock()
	defer mux.mu.RUnlock()

	h := mux.match(host, path)

	if h == nil {
		// Try wildcard
		if wildcard, ok := getWildcard(host); ok {
			if u, ok := mux.redirectToPathSlash(wildcard, path, r.URL); ok {
				return StatusHandler(StatusPermanentRedirect, u.String())
			}
			h = mux.match(wildcard, path)
		}
	}

	if h == nil {
		// Try empty host
		if u, ok := mux.redirectToPathSlash("", path, r.URL); ok {
			return StatusHandler(StatusPermanentRedirect, u.String())
		}
		h = mux.match("", path)
	}

	if h == nil {
		h = NotFoundHandler()
	}

	return h
}

// ServeGemini dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *Mux) ServeGemini(ctx context.Context, w ResponseWriter, r *Request) {
	h := mux.Handler(r)
	h.ServeGemini(ctx, w, r)
}

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *Mux) Handle(pattern string, handler Handler) {
	if pattern == "" {
		panic("gemini: invalid pattern")
	}
	if handler == nil {
		panic("gemini: nil handler")
	}

	mux.mu.Lock()
	defer mux.mu.Unlock()

	var host, path string
	// extract hostname and path
	cut := strings.Index(pattern, "/")
	if cut == -1 {
		host = pattern
		path = "/"
	} else {
		host = pattern[:cut]
		path = pattern[cut:]
	}

	// strip port from hostname
	if hostname, _, err := net.SplitHostPort(host); err == nil {
		host = hostname
	}

	if _, exist := mux.m[hostpath{host, path}]; exist {
		panic("gemini: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[hostpath]Handler)
	}
	mux.m[hostpath{host, path}] = handler
	e := muxEntry{handler, host, path}
	if path[len(path)-1] == '/' {
		mux.es = appendSorted(mux.es, e)
	}
}

func appendSorted(es []muxEntry, e muxEntry) []muxEntry {
	n := len(es)
	i := sort.Search(n, func(i int) bool {
		return len(es[i].path) < len(e.path)
	})
	if i == n {
		return append(es, e)
	}
	// we now know that i points at where we want to insert
	es = append(es, muxEntry{}) // try to grow the slice in place, any entry works.
	copy(es[i+1:], es[i:])      // move shorter entries down
	es[i] = e
	return es
}

// HandleFunc registers the handler function for the given pattern.
func (mux *Mux) HandleFunc(pattern string, handler HandlerFunc) {
	mux.Handle(pattern, handler)
}
