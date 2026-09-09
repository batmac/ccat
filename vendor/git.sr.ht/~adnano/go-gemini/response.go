package gemini

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strconv"
)

// The default media type for responses.
const defaultMediaType = "text/gemini"

// Response represents the response from a Gemini request.
//
// The Client returns Responses from servers once the response
// header has been received. The response body is streamed on demand
// as the Body field is read.
type Response struct {
	// Status is the response status code.
	Status Status

	// Meta returns the response meta.
	// For successful responses, the meta should contain the media type of the response.
	// For failure responses, the meta should contain a short description of the failure.
	Meta string

	// Body represents the response body.
	//
	// The response body is streamed on demand as the Body field
	// is read. If the network connection fails or the server
	// terminates the response, Body.Read calls return an error.
	//
	// The Gemini client guarantees that Body is always
	// non-nil, even on responses without a body or responses with
	// a zero-length body. It is the caller's responsibility to
	// close Body.
	Body io.ReadCloser

	conn net.Conn
}

// ReadResponse reads a Gemini response from the provided io.ReadCloser.
func ReadResponse(r io.ReadCloser) (*Response, error) {
	resp := &Response{}

	// Limit response header size
	lr := io.LimitReader(r, 1029)
	// Wrap the reader to remove the limit later on
	wr := &struct{ io.Reader }{lr}
	br := bufio.NewReader(wr)

	// Read response header
	b, err := br.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return nil, ErrInvalidResponse
		}
		return nil, err
	}
	if len(b) < 3 {
		return nil, ErrInvalidResponse
	}

	// Read the status
	status, err := strconv.Atoi(string(b[:2]))
	if err != nil {
		return nil, ErrInvalidResponse
	}
	resp.Status = Status(status)

	// Read one space
	if b[2] != ' ' {
		return nil, ErrInvalidResponse
	}

	// Read the meta
	meta, ok := bytes.CutSuffix(b[3:], crlf)
	if !ok {
		return nil, ErrInvalidResponse
	}
	if len(meta) == 0 {
		return nil, ErrInvalidResponse
	}
	resp.Meta = string(meta)

	if resp.Status.Class() == StatusSuccess {
		// Use unlimited reader
		wr.Reader = r

		type readCloser struct {
			io.Reader
			io.Closer
		}
		resp.Body = readCloser{br, r}
	} else {
		resp.Body = nopReadCloser{}
		r.Close()
	}
	return resp, nil
}

// Conn returns the network connection on which the response was received.
func (r *Response) Conn() net.Conn {
	return r.conn
}

// TLS returns information about the TLS connection on which the
// response was received.
func (r *Response) TLS() *tls.ConnectionState {
	if tlsConn, ok := r.conn.(*tls.Conn); ok {
		state := tlsConn.ConnectionState()
		return &state
	}
	return nil
}

// WriteTo writes r to w in the Gemini response format, including the
// header and body.
//
// This method consults the Status, Meta, and Body fields of the response.
// The Response Body is closed after it is sent.
func (r *Response) WriteTo(w io.Writer) (int64, error) {
	var wrote int64
	n, err := fmt.Fprintf(w, "%02d %s\r\n", r.Status, r.Meta)
	wrote += int64(n)
	if err != nil {
		return wrote, err
	}
	if r.Body != nil {
		defer r.Body.Close()
		n, err := io.Copy(w, r.Body)
		wrote += n
		if err != nil {
			return wrote, err
		}
	}
	return wrote, nil
}

// A ResponseWriter interface is used by a Gemini handler to construct
// a Gemini response.
//
// A ResponseWriter may not be used after the Handler.ServeGemini method
// has returned.
type ResponseWriter interface {
	// SetMediaType sets the media type that will be sent by Write for a
	// successful response. If no media type is set, a default media type of
	// "text/gemini" will be used.
	//
	// Setting the media type after a call to Write or WriteHeader has
	// no effect.
	SetMediaType(mediatype string)

	// Write writes the data to the connection as part of a Gemini response.
	//
	// If WriteHeader has not yet been called, Write calls WriteHeader with
	// StatusSuccess and the media type set in SetMediaType before writing the data.
	// If no media type was set, Write uses a default media type of
	// "text/gemini".
	Write([]byte) (int, error)

	// WriteHeader sends a Gemini response header with the provided
	// status code and meta.
	//
	// If WriteHeader is not called explicitly, the first call to Write
	// will trigger an implicit call to WriteHeader with a successful
	// status code and the media type set in SetMediaType.
	//
	// The provided code must be a valid Gemini status code.
	// The provided meta must not be longer than 1024 bytes.
	// Only one header may be written.
	WriteHeader(status Status, meta string)

	// Flush sends any buffered data to the client.
	Flush() error
}

type responseWriter struct {
	bw          *bufio.Writer
	mediatype   string
	wroteHeader bool
	bodyAllowed bool
}

func newResponseWriter(w io.Writer) *responseWriter {
	return &responseWriter{
		bw: bufio.NewWriter(w),
	}
}

func (w *responseWriter) SetMediaType(mediatype string) {
	w.mediatype = mediatype
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		meta := w.mediatype
		if meta == "" {
			// Use default media type
			meta = defaultMediaType
		}
		w.WriteHeader(StatusSuccess, meta)
	}
	if !w.bodyAllowed {
		return 0, ErrBodyNotAllowed
	}
	return w.bw.Write(b)
}

func (w *responseWriter) WriteHeader(status Status, meta string) {
	if w.wroteHeader {
		return
	}

	if status.Class() == StatusSuccess {
		w.bodyAllowed = true
	}

	w.bw.WriteString(strconv.Itoa(int(status)))
	w.bw.WriteByte(' ')
	w.bw.WriteString(meta)
	w.bw.Write(crlf)
	w.wroteHeader = true
}

func (w *responseWriter) Flush() error {
	if !w.wroteHeader {
		w.WriteHeader(StatusTemporaryFailure, "Temporary failure")
	}
	// Write errors from WriteHeader will be returned here.
	return w.bw.Flush()
}
