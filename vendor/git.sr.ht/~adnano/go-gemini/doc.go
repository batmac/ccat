/*
Package gemini provides Gemini client and server implementations.

Client is a Gemini client.

	client := &gemini.Client{}
	ctx := context.Background()
	resp, err := client.Get(ctx, "gemini://example.com")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	// ...

Server is a Gemini server.

	server := &gemini.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

Servers should be configured with certificates:

	certificates := &certificate.Store{}
	certificates.Register("localhost")
	err := certificates.Load("/var/lib/gemini/certs")
	if err != nil {
		// handle error
	}
	server.GetCertificate = certificates.Get

Mux is a Gemini request multiplexer.
Mux can handle requests for multiple hosts and paths.

	mux := &gemini.Mux{}
	mux.HandleFunc("example.com", func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		fmt.Fprint(w, "Welcome to example.com")
	})
	mux.HandleFunc("example.org/about.gmi", func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		fmt.Fprint(w, "About example.org")
	})
	mux.HandleFunc("/images/", func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		w.WriteHeader(gemini.StatusGone, "Gone forever")
	})
	server.Handler = mux

To start the server, call ListenAndServe:

	ctx := context.Background()
	err := server.ListenAndServe(ctx)
	if err != nil {
		// handle error
	}
*/
package gemini
