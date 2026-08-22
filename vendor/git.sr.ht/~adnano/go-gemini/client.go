package gemini

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/url"
	"unicode/utf8"

	"golang.org/x/net/idna"
)

// A Client is a Gemini client. Its zero value is a usable client.
type Client struct {
	// TrustCertificate is called to determine whether the client should
	// trust the certificate provided by the server.
	// If TrustCertificate is nil or returns nil, the client will accept
	// any certificate. Otherwise, the certificate will not be trusted
	// and the request will be aborted.
	//
	// See the tofu submodule for an implementation of trust on first use.
	TrustCertificate func(hostname string, cert *x509.Certificate) error

	// DialContext specifies the dial function for creating TCP connections.
	// If DialContext is nil, the client dials using package net.
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
}

// Get sends a Gemini request for the given URL.
// The context controls the entire lifetime of a request and its response:
// obtaining a connection, sending the request, and reading the response
// header and body.
//
// An error is returned if there was a Gemini protocol error.
// A non-2x status code doesn't cause an error.
//
// If the returned error is nil, the user is expected to close the Response.
//
// For more control over requests, use NewRequest and Client.Do.
func (c *Client) Get(ctx context.Context, url string) (*Response, error) {
	req, err := NewRequest(url)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

// Do sends a Gemini request and returns a Gemini response.
// The context controls the entire lifetime of a request and its response:
// obtaining a connection, sending the request, and reading the response
// header and body.
//
// An error is returned if there was a Gemini protocol error.
// A non-2x status code doesn't cause an error.
//
// If the returned error is nil, the user is expected to close the Response.
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	if ctx == nil {
		panic("nil context")
	}

	// Punycode request URL host
	host, port := splitHostPort(req.URL.Host)
	punycode, err := punycodeHostname(host)
	if err != nil {
		return nil, err
	}
	if host != punycode {
		host = punycode

		// Copy the URL and update the host
		u := new(url.URL)
		*u = *req.URL
		u.Host = net.JoinHostPort(host, port)

		// Use the new URL in the request so that the server gets
		// the punycoded hostname
		r := new(Request)
		*r = *req
		r.URL = u
		req = r
	}

	// Use request host if provided
	if req.Host != "" {
		host, port = splitHostPort(req.Host)
		host, err = punycodeHostname(host)
		if err != nil {
			return nil, err
		}
	}

	addr := net.JoinHostPort(host, port)

	// Connect to the host
	conn, err := c.dialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	// Setup TLS
	conn = tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		GetClientCertificate: func(_ *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			if req.Certificate != nil {
				return req.Certificate, nil
			}
			return &tls.Certificate{}, nil
		},
		VerifyConnection: func(cs tls.ConnectionState) error {
			return c.verifyConnection(cs, host)
		},
		ServerName: host,
	})

	type result struct {
		resp *Response
		err  error
	}

	res := make(chan result, 1)
	go func() {
		resp, err := c.do(ctx, conn, req)
		res <- result{resp, err}
	}()

	select {
	case <-ctx.Done():
		conn.Close()
		return nil, ctx.Err()
	case r := <-res:
		if r.err != nil {
			conn.Close()
		}
		return r.resp, r.err
	}
}

func (c *Client) do(ctx context.Context, conn net.Conn, req *Request) (*Response, error) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	w := &contextWriter{
		ctx:    ctx,
		done:   done,
		cancel: cancel,
		wc:     conn,
	}
	rc := &contextReader{
		ctx:    ctx,
		done:   done,
		cancel: cancel,
		rc:     conn,
	}

	// Write the request
	if _, err := req.WriteTo(w); err != nil {
		return nil, err
	}

	// Read the response
	resp, err := ReadResponse(rc)
	if err != nil {
		return nil, err
	}
	resp.conn = conn

	return resp, nil
}

func (c *Client) dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if c.DialContext != nil {
		return c.DialContext(ctx, network, addr)
	}
	return (&net.Dialer{}).DialContext(ctx, network, addr)
}

func (c *Client) verifyConnection(cs tls.ConnectionState, hostname string) error {
	// See if the client trusts the certificate
	if c.TrustCertificate != nil {
		cert := cs.PeerCertificates[0]
		return c.TrustCertificate(hostname, cert)
	}
	return nil
}

func splitHostPort(hostport string) (host, port string) {
	var err error
	host, port, err = net.SplitHostPort(hostport)
	if err != nil {
		// Likely no port
		host = hostport
		port = "1965"
	}
	return
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

// punycodeHostname returns the punycoded version of hostname.
func punycodeHostname(hostname string) (string, error) {
	if net.ParseIP(hostname) != nil {
		return hostname, nil
	}
	if isASCII(hostname) {
		return hostname, nil
	}
	return idna.Lookup.ToASCII(hostname)
}
