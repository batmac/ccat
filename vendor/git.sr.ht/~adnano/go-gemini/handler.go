package gemini

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"strings"
	"time"
)

// A Handler responds to a Gemini request.
//
// ServeGemini should write the response header and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it is not
// valid to use the ResponseWriter after or concurrently with the completion
// of the ServeGemini call.
//
// The provided context is canceled when the client's connection is closed
// or the ServeGemini method returns.
//
// Handlers should not modify the provided Request.
type Handler interface {
	ServeGemini(context.Context, ResponseWriter, *Request)
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions
// as Gemini handlers. If f is a function with the appropriate signature,
// HandlerFunc(f) is a Handler that calls f.
type HandlerFunc func(context.Context, ResponseWriter, *Request)

// ServeGemini calls f(ctx, w, r).
func (f HandlerFunc) ServeGemini(ctx context.Context, w ResponseWriter, r *Request) {
	f(ctx, w, r)
}

// StatusHandler returns a request handler that responds to each request
// with the provided status code and meta.
func StatusHandler(status Status, meta string) Handler {
	return HandlerFunc(func(ctx context.Context, w ResponseWriter, r *Request) {
		w.WriteHeader(status, meta)
	})
}

// NotFoundHandler returns a simple request handler that replies to each
// request with a “51 Not found” reply.
func NotFoundHandler() Handler {
	return StatusHandler(StatusNotFound, "Not found")
}

// StripPrefix returns a handler that serves Gemini requests by removing the
// given prefix from the request URL's Path (and RawPath if set) and invoking
// the handler h. StripPrefix handles a request for a path that doesn't begin
// with prefix by replying with a Gemini 51 not found error. The prefix must
// match exactly: if the prefix in the request contains escaped characters the
// reply is also a Gemini 51 not found error.
func StripPrefix(prefix string, h Handler) Handler {
	if prefix == "" {
		return h
	}
	return HandlerFunc(func(ctx context.Context, w ResponseWriter, r *Request) {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)
		if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
			r2 := new(Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			r2.URL.RawPath = rp
			h.ServeGemini(ctx, w, r2)
		} else {
			w.WriteHeader(StatusNotFound, "Not found")
		}
	})
}

// TimeoutHandler returns a Handler that runs h with the given time limit.
//
// The new Handler calls h.ServeGemini to handle each request, but
// if a call runs for longer than its time limit, the handler responds with a
// 40 Temporary Failure status code and the given message in its response meta.
// After such a timeout, writes by h to its ResponseWriter will return
// context.DeadlineExceeded.
func TimeoutHandler(h Handler, dt time.Duration, message string) Handler {
	return &timeoutHandler{
		h:   h,
		dt:  dt,
		msg: message,
	}
}

type timeoutHandler struct {
	h   Handler
	dt  time.Duration
	msg string
}

func (t *timeoutHandler) ServeGemini(ctx context.Context, w ResponseWriter, r *Request) {
	ctx, cancel := context.WithTimeout(ctx, t.dt)
	defer cancel()

	buf := &bytes.Buffer{}
	tw := &timeoutWriter{
		wr: &contextWriter{
			ctx:    ctx,
			cancel: cancel,
			done:   ctx.Done(),
			wc:     nopCloser{buf},
		},
	}

	done := make(chan struct{})
	go func() {
		t.h.ServeGemini(ctx, tw, r)
		close(done)
	}()

	select {
	case <-done:
		w.WriteHeader(tw.status, tw.meta)
		w.Write(buf.Bytes())
	case <-ctx.Done():
		w.WriteHeader(StatusTemporaryFailure, t.msg)
	}
}

type timeoutWriter struct {
	wr          io.Writer
	status      Status
	meta        string
	mediatype   string
	wroteHeader bool
}

func (w *timeoutWriter) SetMediaType(mediatype string) {
	w.mediatype = mediatype
}

func (w *timeoutWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(StatusSuccess, w.mediatype)
	}
	return w.wr.Write(b)
}

func (w *timeoutWriter) WriteHeader(status Status, meta string) {
	if w.wroteHeader {
		return
	}
	w.status = status
	w.meta = meta
	w.wroteHeader = true
}

func (w *timeoutWriter) Flush() error {
	return nil
}
