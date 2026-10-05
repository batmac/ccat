package gemini

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

// A Server defines parameters for running a Gemini server. The zero value for
// Server is a valid configuration.
type Server struct {
	// Addr optionally specifies the TCP address for the server to listen on,
	// in the form "host:port". If empty, ":1965" (port 1965) is used.
	// See net.Dial for details of the address format.
	Addr string

	// The Handler to invoke.
	Handler Handler

	// ReadTimeout is the maximum duration for reading the entire
	// request.
	//
	// A ReadTimeout of zero means no timeout.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response.
	//
	// A WriteTimeout of zero means no timeout.
	WriteTimeout time.Duration

	// GetCertificate returns a TLS certificate based on the given
	// hostname.
	//
	// If GetCertificate is nil or returns nil, then no certificate
	// will be used and the connection will be aborted.
	//
	// See the certificate submodule for a certificate store that creates
	// and rotates certificates as needed.
	GetCertificate func(hostname string) (*tls.Certificate, error)

	// ErrorLog specifies an optional logger for errors accepting connections,
	// unexpected behavior from handlers, and underlying file system errors.
	// If nil, logging is done via the log package's standard logger.
	ErrorLog interface {
		Printf(format string, v ...interface{})
	}

	listeners map[*net.Listener]context.CancelFunc
	conns     map[*net.Conn]context.CancelFunc
	closed    bool // true if Close or Shutdown called
	shutdown  bool // true if Shutdown called
	doneChan  chan struct{}
	mu        sync.Mutex
}

func (srv *Server) isClosed() bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.closed
}

// done returns a channel that's closed when the server is closed and
// all listeners and connections are closed.
func (srv *Server) done() chan struct{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.doneLocked()
}

func (srv *Server) doneLocked() chan struct{} {
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	return srv.doneChan
}

// tryCloseDone closes srv.done() if the server is closed and
// there are no active listeners or connections.
func (srv *Server) tryCloseDone() {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.tryCloseDoneLocked()
}

func (srv *Server) tryCloseDoneLocked() {
	if !srv.closed {
		return
	}
	if len(srv.listeners) == 0 && len(srv.conns) == 0 {
		ch := srv.doneLocked()
		select {
		case <-ch:
		default:
			close(ch)
		}
	}
}

// Close immediately closes all active net.Listeners and connections.
// For a graceful shutdown, use Shutdown.
func (srv *Server) Close() error {
	srv.mu.Lock()
	{
		if srv.closed {
			srv.mu.Unlock()
			return nil
		}
		srv.closed = true

		srv.tryCloseDoneLocked()

		// Close all active connections and listeners.
		for _, cancel := range srv.listeners {
			cancel()
		}
		for _, cancel := range srv.conns {
			cancel()
		}
	}
	srv.mu.Unlock()

	select {
	case <-srv.done():
		return nil
	}
}

// Shutdown gracefully shuts down the server without interrupting any
// active connections. Shutdown works by first closing all open listeners
// and then waiting indefinitely for connections to close.
// If the provided context expires before the shutdown is complete,
// Shutdown returns the context's error.
//
// When Shutdown is called, Serve and ListenAndServe immediately
// return an error. Make sure the program doesn't exit and waits instead for
// Shutdown to return.
//
// Once Shutdown has been called on a server, it may not be reused;
// future calls to methods such as Serve will return an error.
func (srv *Server) Shutdown(ctx context.Context) error {
	srv.mu.Lock()
	{
		if srv.closed {
			srv.mu.Unlock()
			return nil
		}
		srv.closed = true
		srv.shutdown = true

		srv.tryCloseDoneLocked()

		// Close all active listeners.
		for _, cancel := range srv.listeners {
			cancel()
		}
	}
	srv.mu.Unlock()

	// Wait for active connections to finish.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-srv.done():
		return nil
	}
}

// ListenAndServe listens for requests at the server's configured address.
// ListenAndServe listens on the TCP network address srv.Addr and then calls
// Serve to handle requests on incoming connections. If the provided
// context expires, ListenAndServe closes l and returns the context's error.
//
// If srv.Addr is blank, ":1965" is used.
//
// ListenAndServe always returns a non-nil error.
// After Shutdown or Closed, the returned error is context.Canceled.
func (srv *Server) ListenAndServe(ctx context.Context) error {
	if srv.isClosed() {
		return context.Canceled
	}

	addr := srv.Addr
	if addr == "" {
		addr = ":1965"
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	l = tls.NewListener(l, &tls.Config{
		ClientAuth:     tls.RequestClientCert,
		MinVersion:     tls.VersionTLS12,
		GetCertificate: srv.getCertificate,
	})
	return srv.Serve(ctx, l)
}

func (srv *Server) getCertificate(h *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if srv.GetCertificate == nil {
		return nil, errors.New("gemini: GetCertificate is nil")
	}
	return srv.GetCertificate(h.ServerName)
}

func (srv *Server) trackListener(l *net.Listener, cancel context.CancelFunc) bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.closed {
		return false
	}
	if srv.listeners == nil {
		srv.listeners = make(map[*net.Listener]context.CancelFunc)
	}
	srv.listeners[l] = cancel
	return true
}

func (srv *Server) deleteListener(l *net.Listener) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	delete(srv.listeners, l)
}

// Serve accepts incoming connections on the Listener l, creating a new
// service goroutine for each. The service goroutines read the requests and
// then call the appropriate Handler to reply to them. If the provided
// context expires, Serve closes l and returns the context's error.
//
// Serve always closes l and returns a non-nil error.
// After Shutdown or Close, the returned error is context.Canceled.
func (srv *Server) Serve(ctx context.Context, l net.Listener) error {
	defer l.Close()

	lnctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if !srv.trackListener(&l, cancel) {
		return context.Canceled
	}
	defer srv.tryCloseDone()
	defer srv.deleteListener(&l)

	errch := make(chan error, 1)
	go func() {
		errch <- srv.serve(ctx, l)
	}()

	select {
	case <-lnctx.Done():
		return lnctx.Err()
	case err := <-errch:
		return err
	}
}

func (srv *Server) serve(ctx context.Context, l net.Listener) error {
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, err := l.Accept()
		if err != nil {
			// If this is a temporary error, sleep
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.logf("gemini: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0
		go srv.serveConn(ctx, rw, false)
	}
}

func (srv *Server) trackConn(conn *net.Conn, cancel context.CancelFunc, external bool) bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	// Reject the connection under the following conditions:
	// - Shutdown or Close has been called and conn is external (from ServeConn)
	// - Close (not Shutdown) has been called and conn is internal (from Serve)
	if srv.closed && (external || !srv.shutdown) {
		return false
	}
	if srv.conns == nil {
		srv.conns = make(map[*net.Conn]context.CancelFunc)
	}
	srv.conns[conn] = cancel
	return true
}

func (srv *Server) deleteConn(conn *net.Conn) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	delete(srv.conns, conn)
}

// ServeConn serves a Gemini response over the provided connection.
// It closes the connection when the response has been completed.
// If the provided context expires before the response has completed,
// ServeConn closes the connection and returns the context's error.
func (srv *Server) ServeConn(ctx context.Context, conn net.Conn) error {
	return srv.serveConn(ctx, conn, true)
}

func (srv *Server) serveConn(ctx context.Context, conn net.Conn, external bool) error {
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if !srv.trackConn(&conn, cancel, external) {
		return context.Canceled
	}
	defer srv.tryCloseDone()
	defer srv.deleteConn(&conn)

	if d := srv.ReadTimeout; d != 0 {
		conn.SetReadDeadline(time.Now().Add(d))
	}
	if d := srv.WriteTimeout; d != 0 {
		conn.SetWriteDeadline(time.Now().Add(d))
	}

	errch := make(chan error, 1)
	go func() {
		errch <- srv.goServeConn(ctx, conn)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errch:
		return err
	}
}

func (srv *Server) goServeConn(ctx context.Context, conn net.Conn) error {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	cw := &contextWriter{
		ctx:    ctx,
		done:   done,
		cancel: cancel,
		wc:     conn,
	}
	r := &contextReader{
		ctx:    ctx,
		done:   done,
		cancel: cancel,
		rc:     conn,
	}

	w := newResponseWriter(cw)

	req, err := ReadRequest(r)
	if err != nil {
		w.WriteHeader(StatusBadRequest, "Bad request")
		return w.Flush()
	}
	req.conn = conn

	h := srv.Handler
	if h == nil {
		w.WriteHeader(StatusNotFound, "Not found")
		return w.Flush()
	}

	h.ServeGemini(ctx, w, req)
	return w.Flush()
}

func (srv *Server) logf(format string, args ...interface{}) {
	if srv.ErrorLog != nil {
		srv.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
