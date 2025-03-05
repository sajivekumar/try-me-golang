package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
    caCert, err := os.ReadFile("./certs/server.pem")
	if err != nil {
		log.Fatal(err)
	}
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)
    
    // Create a HTTPS client and supply the created CA pool
	client := &http.Client{
    		Transport: &http.Transport{
    			TLSClientConfig: &tls.Config{
    				RootCAs: caCertPool,
    			},
    		},
    	}

    
	// Request https /hello over port 8443 via the GET method
	r, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		log.Fatal(err)
	}

	// Read the response body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body to stdout
	fmt.Printf("%s\n", body)
}