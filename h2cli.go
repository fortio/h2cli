// Copyright 2023 Fortio Authors
// License: Apache 2.0
//
// Feel free to inspire (copy) from this code but linking back to
// https://github.com/fortio would be appreciated.

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"fortio.org/cli"
	"fortio.org/log"
	"golang.org/x/net/http2"
)

var (
	h2      = flag.Bool("h2", true, "use HTTP2")
	urlFlag = flag.String("url", "http://debug.fortio.org", "URL to fetch")
	method  = flag.String("method", "GET", "HTTP method to use")
	caCert  = flag.String("cacert", "",
		"`Path` to a custom CA certificate file instead standard internet/system CAs")
	stream = flag.Bool("stream", false, "stream stdin to server and back (h2 mode only)")
)

func main() {
	cli.Main() // Will have either called cli.ExitFunction or everything is valid
	client := &http.Client{}
	// Create TLS configuration with the certificate of the server
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}

	if *caCert != "" {
		ca, err := os.ReadFile(*caCert)
		if err != nil {
			log.Fatalf("Reading server certificate: %s", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = caCertPool
	}
	lu := strings.ToLower(*urlFlag)
	if !strings.HasPrefix(lu, "https://") && !strings.HasPrefix(lu, "http://") {
		// be nice and add http:// if missing
		log.LogVf("Adding http:// to url %q", *urlFlag)
		*urlFlag = "http://" + *urlFlag
	}
	u, err := url.Parse(*urlFlag)
	if err != nil {
		log.Fatalf("Failed to parse url %q: %v", *urlFlag, err)
	}
	h2c := ""
	if u.Scheme == "https" || !*h2 {
		// This will do h2 over tls if the server supports it but not h2c
		// but with TLS all is good in the standard package, with ForceAttemptHTTP2
		client.Transport = &http.Transport{
			TLSClientConfig:   tlsConfig, // only used if url is https
			ForceAttemptHTTP2: *h2,       // *h2 for h2 over tls
		}
	} else {
		// h2c
		h2c = "h2c "
		client.Transport = &http2.Transport{
			AllowHTTP: true, // to get h2c
			// Trick to get h2c without TLS:
			// thanks to https://github.com/thrawn01/h2c-golang-example/blob/master/README.md
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		}
	}
	log.Infof("%s%s on %s", h2c, *method, *urlFlag)
	// Perform the request
	var bodyReader io.Reader
	if *stream {
		bodyReader = os.Stdin
	}
	req, err := http.NewRequestWithContext(context.Background(), *method, *urlFlag, bodyReader)
	if err != nil {
		log.Fatalf("Request method %q url %q error: %v", *method, *urlFlag, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed %q %q - error: %v", *method, *urlFlag, err)
	}
	if *stream {
		n, err := io.Copy(os.Stdout, resp.Body)
		log.Infof("Response code %d, proto %s, size %d", resp.StatusCode, resp.Proto, n)
		if err != nil {
			log.Fatalf("Error copying response body: %v", err)
		}
		return
	}
	// else, traditional read all reply mode
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatalf("Failed reading response body: %v", err)
	}
	log.Infof("Response code %d, proto %s, size %d", resp.StatusCode, resp.Proto, len(body))
	os.Stdout.Write(body)
}
