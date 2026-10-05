package gemini

import (
	"context"
	"log"
)

// LoggingMiddleware returns a handler that wraps h and logs Gemini requests
// and their responses to the log package's standard logger.
// Requests are logged with the format "gemini: {host} {URL} {status code} {bytes written}".
func LoggingMiddleware(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, w ResponseWriter, r *Request) {
		lw := &logResponseWriter{rw: w}
		h.ServeGemini(ctx, lw, r)
		host := r.ServerName()
		log.Printf("gemini: %s %q %d %d", host, r.URL, lw.Status, lw.Wrote)
	})
}

type logResponseWriter struct {
	Status      Status
	Wrote       int
	rw          ResponseWriter
	mediatype   string
	wroteHeader bool
}

func (w *logResponseWriter) SetMediaType(mediatype string) {
	w.mediatype = mediatype
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		meta := w.mediatype
		if meta == "" {
			// Use default media type
			meta = defaultMediaType
		}
		w.WriteHeader(StatusSuccess, meta)
	}
	n, err := w.rw.Write(b)
	w.Wrote += n
	return n, err
}

func (w *logResponseWriter) WriteHeader(status Status, meta string) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.Status = status
	w.Wrote += len(meta) + 5
	w.rw.WriteHeader(status, meta)
}

func (w *logResponseWriter) Flush() error {
	return nil
}
