package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"net/http"
	"os"

	"fortio.org/fortio/log"
)

var (
	h2     = flag.Bool("h2", true, "use HTTP2")
	url    = flag.String("url", "https://localhost:8080/debug", "URL to fetch")
	method = flag.String("method", "GET", "HTTP method to use")
	caCert = flag.String("cacert", "",
		"`Path` to a custom CA certificate file instead standard internet/system CAs")
)

func main() {
	flag.Parse()
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

	client.Transport = &http.Transport{
		TLSClientConfig:   tlsConfig,
		ForceAttemptHTTP2: *h2, // could also use &http2.Transport{TLSClientConfig: tlsConfig} but that's not necessary to get h2
	}

	// Perform the request
	req, err := http.NewRequestWithContext(context.Background(), *method, *url, nil)
	if err != nil {
		log.Fatalf("Request method %q url %q error: %v", *method, *url, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed %q %q - error: %v", *method, *url, err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatalf("Failed reading response body: %v", err)
	}
	log.Printf("Response code %d, proto %s", resp.StatusCode, resp.Proto)
	os.Stdout.Write(body)
}
