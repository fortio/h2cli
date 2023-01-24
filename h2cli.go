package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"net/http"
	"os"

	"fortio.org/fortio/log"

	"golang.org/x/net/http2"
)

var (
	h2     = flag.Bool("h2", true, "use HTTP2")
	url    = flag.String("url", "https://localhost:8080/debug", "URL to fetch")
	caCert = flag.String("cacert", "",
		"`Path` to a custom CA certificate file instead standard internet/system CAs")
)

func main() {
	flag.Parse()
	client := &http.Client{}
	// Create TLS configuration with the certificate of the server
	tlsConfig := &tls.Config{}

	if *caCert != "" {
		ca, err := os.ReadFile(*caCert)
		if err != nil {
			log.Fatalf("Reading server certificate: %s", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = caCertPool
	}

	// Use the proper transport in the client
	if *h2 {
		client.Transport = &http2.Transport{TLSClientConfig: tlsConfig}
	} else {
		client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	// Perform the request
	resp, err := client.Get(*url)
	if err != nil {
		log.Fatalf("Failed get: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed reading response body: %v", err)
	}
	log.Printf("Response code %d, proto %s", resp.StatusCode, resp.Proto)
	os.Stdout.Write(body)
}
