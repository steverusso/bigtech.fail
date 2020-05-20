// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

type config struct {
	dir  string
	prod bool
	port int
}

func (c *config) addr() string {
	return fmt.Sprintf(":%d", c.port)
}

func main() {
	// Read cli options into config.
	var c config
	flag.StringVar(&c.dir, "dir", "", "The directory of static resources to serve.")
	flag.BoolVar(&c.prod, "prod", false, "Run server in production mode (use TLS, etc.).")
	flag.IntVar(&c.port, "port", 8080, "The port to use for the non-production server.")
	flag.Parse()

	fs := http.FileServer(http.Dir(c.dir)) // File server for the given directory.

	if c.prod {
		// Production HTTPS server on port 443 (using autocert convenience function).
		go func() {
			if err := http.Serve(autocert.NewListener("bigtech.fail", "www.bigtech.fail"), fs); err != nil {
				log.Fatal(err)
			}
		}()

		// Redirect HTTP requests on port 80 to HTTPS.
		redirToHTTPS := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://bigtech.fail:443"+r.RequestURI, http.StatusMovedPermanently)
		})
		if err := http.ListenAndServe(":80", redirToHTTPS); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(http.ListenAndServe(c.addr(), fs))
	}
}
