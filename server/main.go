// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

//go:embed public
var website embed.FS

func main() {
	host := flag.String("host", "localhost", "The web host.")
	port := flag.Int("p", 0, "The web port. If specified, TLS will not be used.")
	flag.Parse()

	website, err := fs.Sub(website, "public")
	if err != nil {
		log.Fatalf("couldn't get 'public' subdir of embedded fs: %v", err)
	}

	fs := http.FileServer(http.FS(website))

	h := http.NewServeMux()
	h.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		switch uri := r.RequestURI; uri {
		case "/events", "/events/":
			fs.ServeHTTP(w, r)
		default:
			http.Redirect(w, r, "/e/"+uri[8:], http.StatusMovedPermanently)
		}
	})
	h.Handle("/", fs)

	if *port == 0 {
		if err := serveTLS(h); err != nil {
			log.Fatalf("couldn't serve tls: %v", err)
		}
	} else {
		addr := fmt.Sprintf("%s:%d", *host, *port)
		if err := http.ListenAndServe(addr, h); err != nil {
			log.Fatalf("couldn't listen and serve: %v", err)
		}
	}
}

// serveTLS will serve the given handler over HTTPS using the `autocert`
// package and redirect all requests on port 80 to HTTPS on port 443.
func serveTLS(h http.Handler) error {
	errCh := make(chan error)

	go func() {
		if err := http.Serve(autocert.NewListener("bigtech.fail"), h); err != nil {
			errCh <- fmt.Errorf("couldn't serve autocert listener: %w", err)
		}
	}()

	go func() {
		redir := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+":443"+r.RequestURI, http.StatusMovedPermanently)
		})

		if err := http.ListenAndServe(":80", redir); err != nil {
			errCh <- fmt.Errorf("couldn't redirect port 80: %w", err)
		}
	}()

	return <-errCh
}
